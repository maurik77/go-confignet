package providers

import (
	"log"

	"github.com/maurik77/go-confignet/extensions"
)

type stdLogger struct{}

func (stdLogger) Printf(format string, args ...interface{}) { log.Printf(format, args...) }

// logger is the active logger for the providers package. Use confignet.SetLogger to change it.
var logger extensions.Logger = stdLogger{}

// SetLogger sets the logger used by the providers package.
func SetLogger(l extensions.Logger) {
	logger = l
}
