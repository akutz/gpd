package main

import (
	"C"

	"fmt"

	"github.com/akutz/gpd/dog/stay/lib"
)

// Command prints a dog's name to stdout.
func Command(d lib.Dog) {
	fmt.Println(d.Name())
}
