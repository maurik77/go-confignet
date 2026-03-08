package internal

import (
	"log"

	"github.com/maurik77/go-confignet/extensions"
)

type stdLogger struct{}

func (stdLogger) Printf(format string, args ...interface{}) { log.Printf(format, args...) }

// Logger is the active logger for the internal package. Use confignet.SetLogger to change it.
var Logger extensions.Logger = stdLogger{}

// SetLogger sets the logger used by the internal package.
func SetLogger(l extensions.Logger) {
	Logger = l
}
