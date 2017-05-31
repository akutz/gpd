// +build !self

package main

// Dog is (wo)?man's best friend.
type Dog interface {
	// Name returns the name of the dog.
	Name() string
}
