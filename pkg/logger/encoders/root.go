package encoders

import (
	"go.uber.org/zap/buffer"
)

const (
	_hex = "0123456789abcdef"
)

//lint:ignore GLOBAL this is okay
var (
	bufferpool = buffer.NewPool()
)

// Setup sets encoders
func Setup() {
}
