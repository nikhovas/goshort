package types

type ModuleStateChangeLog struct {
	ModuleName string
	ModuleType string
	State      string
}

func (log *ModuleStateChangeLog) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Common.ModuleStateChange",
		"type":       "log",
		"moduleName": log.ModuleName,
		"moduleType": log.ModuleType,
		"state":      log.State,
	}
}
