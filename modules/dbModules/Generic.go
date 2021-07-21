package dbModules

import (
	"goshort/kernel"
	"goshort/types"
	"sync"
)

type Generic struct {
	types.ModuleBase
	GetFunc               func(key string) (types.Url, error)
	PostFunc              func(newUrl types.Url) (types.Url, error)
	PatchFunc             func(patchUrl types.Url) error
	DeleteFunc            func(key string) error
	GenericKeySupportFunc func() bool
	Name                  string
	Kernel                kernel.Kernel
}

func (controller *Generic) Init(_ map[string]interface{}) error {
	return nil
}

func (controller *Generic) Run(wg *sync.WaitGroup) error {
	wg.Done()
	return nil
}

func (controller *Generic) Stop() error {
	return nil
}

func (controller *Generic) Get(key string) (types.Url, error) {
	return controller.GetFunc(key)
}

func (controller *Generic) Post(newUrl types.Url) (types.Url, error) {
	return controller.PostFunc(newUrl)
}

func (controller *Generic) Patch(patchUrl types.Url) error {
	return controller.PatchFunc(patchUrl)
}

func (controller *Generic) Delete(key string) error {
	return controller.DeleteFunc(key)
}

func (controller *Generic) GenericKeySupport() bool {
	return controller.GenericKeySupportFunc()
}

func (controller *Generic) GetName() string {
	return controller.Name
}

func (controller *Generic) GetType() string {
	return "Database.Generic"
}
