package lib

import (
	"context"

	"github.com/akutz/gpd/dep"
)

// Module is the interface implementated by types that
// register themselves as modular plug-ins.
type Module interface {

	// Init initializes the module.
	Init(ctx context.Context, config dep.Config)
}

var mods = map[string]func() Module{}

// RegisterModule registers a new module with its name and function
// that returns a new, uninitialized instance of the module type.
func RegisterModule(name string, ctor func() Module) {
	mods[name] = ctor
}

// NewModule instantiates a new instance of the module type with the
// specified name.
func NewModule(name string) Module {
	return mods[name]()
}
