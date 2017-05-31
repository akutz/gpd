package main

import (
	"C"

	"fmt"
)

// Command prints a dog's name to stdout.
func Command(d Dog) {
	fmt.Println(d.Name())
}
