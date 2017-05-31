package main

import (
	"fmt"
	"os"
	"plugin"
)

// lucy is a type that implements the Dog interface, returning
// "Lucy" from the Name() string function.
type lucy struct{}

func (l *lucy) Name() string { return "Lucy" }
func (l *lucy) Self() Dog    { return l }

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "invalid args\n")
		os.Exit(1)
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

	// lookup the plug-in's Command symbol; it's the function used to
	// print a dog's name
	cmdFunc, err := p.Lookup("Command")
	if err != nil {
		fmt.Fprintf(
			os.Stderr, "error: failed to lookup Command func: %v\n", err)
		os.Exit(1)
	}

	// assert that the Command symbol is a func(Dog)
	cmd, ok := cmdFunc.(func(Dog))
	if !ok {
		fmt.Fprintf(os.Stderr, "error: invalid Command func: %T\n", cmdFunc)
		os.Exit(1)
	}

	// issue the command to Lucy
	cmd(&lucy{})
}
