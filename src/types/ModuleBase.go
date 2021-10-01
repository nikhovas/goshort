package types

import (
	"goshort/kernel/utils"
	"goshort/types/errors"
	"sync/atomic"
)

type ModuleBase struct {
	MaxRetryAttempts     int
	RetryAttemptInterval int
	IsAvailableVal       int32
	death                int32
}

func (mb *ModuleBase) Init(config map[string]interface{}) error {
	mb.MaxRetryAttempts = utils.UnwrapFieldOrDefault(config, "max_retry_attempts", 10).(int)
	mb.RetryAttemptInterval = utils.UnwrapFieldOrDefault(config, "retry_attempt_interval", 1).(int)
	mb.IsAvailableVal = 0
	return nil
}

func (mb *ModuleBase) GetMaxRetryAttempts() int {
	return mb.MaxRetryAttempts
}

func (mb *ModuleBase) GetRetryAttemptInterval() int {
	return mb.RetryAttemptInterval
}

func (mb *ModuleBase) IsAvailable() bool {
	return mb.IsAvailableVal == 1
}

func (mb *ModuleBase) SetAvailable() {
	mb.IsAvailableVal = 1
}

func (mb *ModuleBase) SetUnavailableAndTryGetReconnectionControl() bool {
	prev := atomic.SwapInt32(&mb.IsAvailableVal, 0)
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
