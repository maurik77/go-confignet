package providers

import (
	"log"
	"os"
	"strings"

	"github.com/Maurik77/go-confignet/extensions"
)

// CmdLineConfigurationProvider loads configuration from commandline arguments
type CmdLineConfigurationProvider struct {
	Prefix       string
	RemovePrefix bool
	KeyMapper    func(arg string) string
	data         map[string]string
}

// Load configuration from commandline arguments
func (provider *CmdLineConfigurationProvider) Load(decrypter extensions.IConfigurationDecrypter) {
	provider.data = make(map[string]string)

	for _, arg := range os.Args[1:] {
		parts := strings.Split(arg, "=")
		key := parts[0]
		var value string

		if len(parts) > 1 {
			value = strings.Join(parts[1:], "")
		} else {
			value = "true"
		}

		if provider.KeyMapper == nil {
			// Remove the prefix "-" if the argument is -arg or --arg
			key = strings.TrimLeft(key, "-")
		} else {
			key = provider.KeyMapper(key)
		}

		if provider.Prefix != "" && !strings.HasPrefix(key, provider.Prefix) {
			continue
		}

		if provider.Prefix != "" && provider.RemovePrefix {
			key = strings.TrimPrefix(key, provider.Prefix)
		}

		if decrypter != nil {
			var err error
			value, err = decrypter.Decrypt(value)

			if err != nil {
				log.Printf("CmdLineConfigurationProvider:Error calling decryption for key %v. %v", key, err)
			}
		}

		provider.data[key] = value
	}
}

// GetData provides the loaded data
func (provider *CmdLineConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *CmdLineConfigurationProvider) GetSeparator() string {
	return "-"
}
