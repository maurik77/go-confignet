package confignet

import (
	"log"

	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/internal"
	"github.com/maurik77/go-confignet/providers"
)

type stdLogger struct{}

func (stdLogger) Printf(format string, args ...interface{}) { log.Printf(format, args...) }

// logger is the active logger for the confignet package.
var logger extensions.Logger = stdLogger{}

// SetLogger sets the logger used by all confignet components.
// Call this before building any configuration.
//
// Example with logrus:
//
//	confignet.SetLogger(logrus.StandardLogger())
//
// Example with a custom logger:
//
//	confignet.SetLogger(myLogger)
func SetLogger(l extensions.Logger) {
	if l != nil {
		logger = l
		internal.SetLogger(l)
		providers.SetLogger(l)
	}
}
