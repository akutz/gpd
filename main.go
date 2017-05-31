package main

import (
	"context"
	"fmt"
	"os"
	"plugin"

	"github.com/akutz/gpd/lib"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Yes, we have no bananas,")
		fmt.Println("We have no bananas today.")
		os.Exit(0)
	}

	// validate that the plug-in file exists
	pluginPath := os.Args[1]
	if !func(path string) bool {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return true
		}
		return false
	}(pluginPath) {
		fmt.Fprintf(os.Stderr, "error: invalid plugin file: %s\n", pluginPath)
		os.Exit(1)
	}

	// open the plug-in file
	p, err := plugin.Open(pluginPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load plugin: %v\n", err)
		os.Exit(1)
	}

	// lookup the plug-in's Types symbol; it's the type map used to
	// register the plug-in's modules
	tmapObj, err := p.Lookup("Types")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to lookup type map: %v\n", err)
		os.Exit(1)
	}

	// assert that the Types symbol is a *map[string]func() interface{}
	tmapPtr, tmapOk := tmapObj.(*map[string]func() interface{})
	if !tmapOk {
		fmt.Fprintf(os.Stderr, "error: invalid type map: %T\n", tmapObj)
		os.Exit(1)
	}

	// assert that the type map pointer is not nil
	if tmapPtr == nil {
		fmt.Fprintf(
			os.Stderr, "error: nil type map: type=%[1]T val=%[1]v\n", tmapPtr)
		os.Exit(1)
	}

	// dereference the type map pointer
	tmap := *tmapPtr

	// register the plug-in's modules
	for k, v := range tmap {
		lib.RegisterModule(k, v)
	}

	// Instantiate a copy of the module registered by the plug-in.
	modGo := lib.NewModule("mod_go")

	// Create a new context
	ctx := context.Background()

	// Create a new v2 config
	config := &v2Config{}

	// Initialize mod_go with a v2 config implementation
	modGo.Init(ctx, config)

	// Set a value in the config
	config.Set(
		ctx,
		"bananas",
		"Bottom-line, sh*t kicking country choir\n"+
			"You'll see your part come by\n")

	// Initialize mod_go again with the updated config
	modGo.Init(ctx, config)
}

type v2Config struct {
	val interface{}
}

// Get returns the value for the specified key
func (c *v2Config) Get(ctx context.Context, key string) interface{} {
	if c.val != nil {
		return c.val
	}
	if key == "bananas" {
		return "Yes there were thirty, thousand, pounds...\n" +
			"Of...bananas.\n"
	}
	return nil
}

// Set sets the value for the specified key
func (c *v2Config) Set(ctx context.Context, key string, val interface{}) {
	c.val = val
}
