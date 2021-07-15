package types

import (
	"goshort/kernel/utils"
	"goshort/types/errors"
	"sync/atomic"
)

type ModuleBase struct {
	MaxRetryAttempts     int
	RetryAttemptInterval int
	isAvailable          int32
	death                int32
}

func (mb *ModuleBase) Init(config map[string]interface{}) error {
	mb.MaxRetryAttempts = utils.UnwrapFieldOrDefault(config, "max_retry_attempts", 10).(int)
	mb.RetryAttemptInterval = utils.UnwrapFieldOrDefault(config, "retry_attempt_interval", 1).(int)
	mb.isAvailable = 1
	return nil
}

func (mb *ModuleBase) GetMaxRetryAttempts() int {
	return mb.MaxRetryAttempts
}

func (mb *ModuleBase) GetRetryAttemptInterval() int {
	return mb.RetryAttemptInterval
}

func (mb *ModuleBase) IsAvailable() bool {
	return mb.isAvailable == 1
}

func (mb *ModuleBase) SetAvailable() {
	mb.isAvailable = 1
}

func (mb *ModuleBase) SetUnavailableAndTryGetReconnectionControl() bool {
	prev := atomic.SwapInt32(&mb.isAvailable, 0)
	return prev == 1
}

func (mb *ModuleBase) GetDeath() bool {
	return mb.death == 1
}

func (mb *ModuleBase) SetDeath() {
	mb.death = 1
}

func (mb *ModuleBase) UnsetDeathAndTryGetReconnectionControl() bool {
	prev := atomic.SwapInt32(&mb.death, 0)
	return prev == 1
}

func (mb *ModuleBase) TryReconnect() error {
	return errors.NotImplementedError
}

//func (mb * ModuleBase) TryGetReconnectControl() bool {
//	prev := atomic.SwapInt32(&mb.isReconnectingNow, 1)
//	return prev == 0
//}
//
//func (mb * ModuleBase) CloseReconnectingMode() {
//	mb.isReconnectingNow = 0
//}

type LoggerBase struct {
	ModuleBase
	SystemLogs_           bool
	ClientConnectionLogs_ bool
	CommonLogger          bool
	ExtraLogger           bool
}

func (lb *LoggerBase) Init(config map[string]interface{}) error {
	if err := lb.ModuleBase.Init(config); err != nil {
		return err
	}
	lb.ClientConnectionLogs_ = utils.UnwrapFieldOrDefault(config, "client_connection_logs", true).(bool)
	lb.SystemLogs_ = utils.UnwrapFieldOrDefault(config, "system_logs", true).(bool)
	lb.CommonLogger = utils.UnwrapFieldOrDefault(config, "common_logger", true).(bool)
	lb.ExtraLogger = utils.UnwrapFieldOrDefault(config, "extra_logger", false).(bool)
	return nil
}

func (lb *LoggerBase) ClientConnectionLogs() bool {
	return lb.ClientConnectionLogs_
}

func (lb *LoggerBase) SystemLogs() bool {
	return lb.SystemLogs_
}

func (lb *LoggerBase) IsCommonLogger() bool {
	return lb.CommonLogger
}

func (lb *LoggerBase) IsExtraLogger() bool {
	return lb.ExtraLogger
}
