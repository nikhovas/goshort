package kernel

import (
	"goshort/src/types"
	"time"
)

type ReconnectionTaskStateChange struct {
	TaskId     uint64
	ModuleName string
	EventName  string
}

func (task ReconnectionTaskStateChange) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Kernel.Reconnection.StateChange",
		"event":      task.EventName,
		"type":       "log",
		"moduleName": task.ModuleName,
		"taskId":     task.TaskId,
	}
}

func (task ReconnectionTaskStateChange) IsError() bool {
	return false
}

type ReconnectionLogModuleInfo struct {
	ModuleName string
	ModuleType string
}

type SuccessReconnection struct {
	ReconnectionLogModuleInfo
}

func (sr SuccessReconnection) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Common.Reconnection.Success",
		"type":       "log",
		"moduleName": sr.ModuleName,
		"moduleType": sr.ModuleType,
	}
}

func (sr SuccessReconnection) IsError() bool {
	return false
}

type BadReconnectionAttempt struct {
	ReconnectionLogModuleInfo
	Attempt int
	Limit   int
	Error   error
}

func (bra *BadReconnectionAttempt) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Common.Reconnection.BadAttempt",
		"type":       "error",
		"moduleName": bra.ModuleName,
		"moduleType": bra.ModuleType,
		"attempt":    bra.Attempt,
		"limit":      bra.Limit,
		"error":      bra.Error,
	}
}

func (bra *BadReconnectionAttempt) IsError() bool {
	return true
}

type ReconnectionAttemptsLimit struct {
	ReconnectionLogModuleInfo
	Limit int
}

func (ral *ReconnectionAttemptsLimit) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Common.Reconnection.ReconnectionAttemptsLimit",
		"type":       "log",
		"moduleName": ral.ModuleName,
		"moduleType": ral.ModuleType,
	}
}

func (ral *ReconnectionAttemptsLimit) IsError() bool {
	return false
}

func getModuleInfo(module types.ModuleInterface) ReconnectionLogModuleInfo {
	return ReconnectionLogModuleInfo{
		ModuleName: module.GetName(),
		ModuleType: module.GetType(),
	}
}

type ReconnectionKernel struct {
	Kernel *Kernel
}

func (reconnectionKernel *ReconnectionKernel) work(module types.ModuleInterface) {
	defer reconnectionKernel.Kernel.wg.Done()

	operationId := reconnectionKernel.Kernel.GetNextOperationNumber()

	maxAttempts := module.GetMaxRetryAttempts()
	retryInterval := module.GetRetryAttemptInterval()

	_ = reconnectionKernel.Kernel.Logger.Send(&ReconnectionTaskStateChange{
		TaskId:     operationId,
		ModuleName: module.GetName(),
		EventName:  "start",
	})
	defer reconnectionKernel.Kernel.Logger.Send(&ReconnectionTaskStateChange{
		TaskId:     operationId,
		ModuleName: module.GetName(),
		EventName:  "end",
	})

	for currentAttempt := 0; currentAttempt < maxAttempts; currentAttempt += 1 {
		if reconnectionKernel.Kernel.Stopped {
			// TODO: logging here
			return
		}
		err := module.TryReconnect()
		if err == nil {
			module.SetAvailable()
			_ = reconnectionKernel.Kernel.Logger.Send(&SuccessReconnection{
				getModuleInfo(module),
			})
			return
		} else {
			_ = reconnectionKernel.Kernel.Logger.Send(&BadReconnectionAttempt{
				ReconnectionLogModuleInfo: getModuleInfo(module),
				Error:                     err,
				Attempt:                   currentAttempt + 1,
				Limit:                     maxAttempts,
			})
		}
		time.Sleep(time.Duration(retryInterval) * time.Second)
	}

	module.SetDeath()
	_ = reconnectionKernel.Kernel.Logger.Send(&ReconnectionAttemptsLimit{
		ReconnectionLogModuleInfo: getModuleInfo(module),
		Limit:                     maxAttempts,
	})
}

func (reconnectionKernel *ReconnectionKernel) Start(module types.ModuleInterface) {
	if module.SetUnavailableAndTryGetReconnectionControl() {
		reconnectionKernel.Kernel.wg.Add(1)
		go reconnectionKernel.work(module)
	}
}
