package main

import (
	"context"
	"fmt"
	"os"
	"plugin"

	"github.com/akutz/gpd/dep"
	"github.com/akutz/gpd/lib"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Yes, we have no bananas,")
		fmt.Println("We have no bananas today.")
		os.Exit(0)
	}
	pluginPath := os.Args[1]
	if !fileExists(pluginPath) {
		fmt.Fprintf(os.Stderr, "error: invalid plugin file: %s\n", pluginPath)
		os.Exit(1)
	}

	// open the plug-in file, causing its package init function to run,
	// thereby registering the module
	if _, err := plugin.Open(pluginPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load plugin: %v\n", err)
	}

	// Instantiate a copy of the module registered by the plug-in.
	modGo := lib.NewModule("mod_go")

	// Initialize mod_go
	modGo.Init(context.Background(), dep.Config{})
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return true
	}
	return false
}
