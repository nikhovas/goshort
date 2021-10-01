package types

import (
	"sync"
)

type KernelInterface interface {
	SetModuleRunState(module ModuleInterface)
	SetModuleStopState(module ModuleInterface)

	GetLoggerCreators() map[string]func(kernel KernelInterface) LoggerInterface
	ErrorManage(err error, module ModuleInterface)
}

type ModuleInterface interface {
	Init(config map[string]interface{}) error
	Run(wg *sync.WaitGroup) error
	Stop() error
	GetName() string
	GetType() string
	GetMaxRetryAttempts() int
	GetRetryAttemptInterval() int
	IsAvailable() bool
	SetAvailable()
	TryReconnect() error
	SetUnavailableAndTryGetReconnectionControl() bool
	GetDeath() bool
	SetDeath()
	UnsetDeathAndTryGetReconnectionControl() bool
}

type InputInterface interface {
	ModuleInterface
}

type DatabaseInterface interface {
	ModuleInterface
	Get(key string) (Url, error)
	Post(newUrl Url) (Url, error)
	Patch(patchUrl Url) error
	Delete(key string) error
	GenericKeySupport() bool
}

type LoggerInterface interface {
	ModuleInterface
	Send(element Log) error
	SendError(err error) error
	SendBatch(batch *LoggingQueueNode) error
}

type MiddlewareInterface interface {
	ModuleInterface
	Exec(url *Url) error
	BreakOnError() bool
}

type Log interface {
	ToMap() map[string]interface{}
}

type LoggingQueueNode struct {
	Next *LoggingQueueNode
	Log  Log
}

func (node *LoggingQueueNode) Len() int {
	counter := 0

	metaElem := node
	for metaElem != nil {
		counter += 1
		metaElem = metaElem.Next
	}

	return counter
}
