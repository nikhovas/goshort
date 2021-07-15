package types

type ModuleInterface interface {
	Init(config map[string]interface{}) error
	Run() error
	Stop() error
	GetName() string
	GetType() string
	GetMaxRetryAttempts() int
	GetRetryAttemptInterval() int
	IsAvailable() bool
	SetAvailable()
	TryReconnect() error
	//TryGetReconnectControl() bool
	//CloseReconnectingMode()
	SetUnavailableAndTryGetReconnectionControl() bool

	GetDeath() bool
	SetDeath()
	UnsetDeathAndTryGetReconnectionControl() bool
}

type InputControllerInterface interface {
	ModuleInterface
}

type UrlControllerInterface interface {
	ModuleInterface
	Get(key string) (Url, error)
	Post(newUrl Url) (Url, error)
	Patch(patchUrl Url) error
	Delete(url_ Url) error
	GenericKeySupport() bool
}

type LoggerInterface interface {
	ModuleInterface
	Send(element Log) error
	ClientConnectionLogs() bool
	SystemLogs() bool
	IsCommonLogger() bool
	IsExtraLogger() bool
}

type MiddlewareInterface interface {
	ModuleInterface
	Exec(url *Url) error
	BreakOnError() bool
}

type Log interface {
	ToMap() map[string]interface{}
}

//type AdvancedError interface {
//	Log
//	Error() string
//}
//
//func GetErrorDescription(err error) interface{} {
//	v, ok := err.(AdvancedError)
//	if ok {
//		return v.ToMap()
//	} else {
//		return err.Error()
//	}
//}
