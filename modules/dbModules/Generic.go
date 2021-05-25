package dbModules

import (
	"goshort/kernel"
)

type Generic struct {
	get               func(key string) (kernel.Url, error)
	post              func(newUrl kernel.Url) (kernel.Url, error)
	patch             func(patchUrl kernel.Url) error
	delete            func(url_ kernel.Url) error
	genericKeySupport func() bool
}

func Create(get func(key string) (kernel.Url, error),
	post func(newUrl kernel.Url) (kernel.Url, error),
	patch func(patchUrl kernel.Url) error,
	delete func(url_ kernel.Url) error,
	genericKeySupport func() bool) Generic {
	return Generic{get: get, post: post, patch: patch, delete: delete, genericKeySupport: genericKeySupport}
}

func (controller *Generic) Run() error {
	return nil
}

func (controller *Generic) Get(key string) (kernel.Url, error) {
	return controller.get(key)
}

func (controller *Generic) Post(newUrl kernel.Url) (kernel.Url, error) {
	return controller.post(newUrl)
}

func (controller *Generic) Patch(patchUrl kernel.Url) error {
	return controller.patch(patchUrl)
}

func (controller *Generic) Delete(url_ kernel.Url) error {
	return controller.delete(url_)
}

func (controller *Generic) GenericKeySupport() bool {
	return controller.genericKeySupport()
}
