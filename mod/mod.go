package main

import "C"

import (
	"context"
	"fmt"
	"os"

	"github.com/akutz/gpd/dep"
	"github.com/akutz/gpd/lib"
)

type module struct{}

func init() {
	lib.RegisterModule("mod_go", func() lib.Module {
		return &module{}
	})
}

func (m *module) Init(ctx context.Context, config dep.Config) {
	fmt.Fprintln(os.Stdout, "Yes there were thirty, thousand, pounds...")
	fmt.Fprintln(os.Stdout, "Of...bananas.")
}
