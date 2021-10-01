package inputModules

import (
	"goshort/kernel"
	"goshort/types"
	"sync"
)

type Generic struct {
	types.ModuleBase
	moduleName string
	Kernel     *kernel.Kernel
}

func (generic *Generic) Get(key string) (types.Url, error) {
	return generic.Kernel.Database.Get(key)
}

func (generic *Generic) Post(url types.Url) (types.Url, error) {
	return generic.Kernel.Database.Post(url)
}

func (generic *Generic) Patch(url types.Url) error {
	return generic.Kernel.Database.Patch(url)
}

func (generic *Generic) Delete(key string) error {
	return generic.Kernel.Database.Delete(key)
}

func (generic *Generic) Run(wg *sync.WaitGroup) error {
	wg.Done()
	return nil
}

func (generic *Generic) Stop() error {
	return nil
}

func (generic *Generic) GetName() string {
	return generic.moduleName
}

func (generic *Generic) GetType() string {
	return "Input.Generic"
}
