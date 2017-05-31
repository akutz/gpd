package main

import "C"

import (
	"context"
	"fmt"
	"os"
)

// Types is the symbol the host process uses to
// retrieve the plug-in's type map
var Types = map[string]func() interface{}{
	"mod_go": func() interface{} { return &module{} },
}

type module struct{}

func (m *module) Init(ctx context.Context, configObj interface{}) error {

	config, configOk := configObj.(Config)
	if !configOk {
		return errInvalidConfig
	}

	fmt.Fprintf(os.Stdout, "%T\n", config)
	fmt.Fprintln(os.Stdout, config.Get(ctx, "bananas"))
	return nil
}
