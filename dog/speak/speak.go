package main

import (
	"C"

	"fmt"
)

// dog is (wo)?man's best friend.
type dog interface {
	// Name returns the name of the dog.
	Name() string
}

// Command prints a dog's name to stdout.
func Command(d dog) {
	fmt.Println(d.Name())
}
