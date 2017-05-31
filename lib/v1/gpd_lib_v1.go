package v1

import (
	"context"
	"errors"
)

// errInvalidConfig is the error returned when a module is initialized
// with an invalid configuration argument
var errInvalidConfig = errors.New("invalid config")

// Config is a configuration provider
type Config interface {
	// Get returns the value for the specified key
	Get(ctx context.Context, key string) interface{}
}
