package dbModules

import (
	"goshort/types"
)

type Generic struct {
	GetFunc               func(key string) (types.Url, error)
	PostFunc              func(newUrl types.Url) (types.Url, error)
	PatchFunc             func(patchUrl types.Url) error
	DeleteFunc            func(url_ types.Url) error
	GenericKeySupportFunc func() bool
	Name                  string
}

func (controller *Generic) Init(config map[string]interface{}) error {
	return nil
}

func (controller *Generic) Run() error {
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

func (controller *Generic) Delete(url_ types.Url) error {
	return controller.DeleteFunc(url_)
}

func (controller *Generic) GenericKeySupport() bool {
	return controller.GenericKeySupportFunc()
}

func (controller *Generic) GetName() string {
	return controller.Name
}

func (controller *Generic) GetType() string {
	return "Generic"
}
