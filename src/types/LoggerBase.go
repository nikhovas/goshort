package types

import "goshort/src/kernel/utils"

type LoggerBase struct {
	ModuleBase
	SystemLogs_           bool
	ClientConnectionLogs_ bool
}

func (lb *LoggerBase) Init(config map[string]interface{}) error {
	if err := lb.ModuleBase.Init(config); err != nil {
		return err
	}
	lb.ClientConnectionLogs_ = utils.UnwrapFieldOrDefault(config, "client_connection_logs", true).(bool)
	lb.SystemLogs_ = utils.UnwrapFieldOrDefault(config, "system_logs", true).(bool)
	return nil
}

func (lb *LoggerBase) ClientConnectionLogs() bool {
	return lb.ClientConnectionLogs_
}

func (lb *LoggerBase) SystemLogs() bool {
	return lb.SystemLogs_
}
