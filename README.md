# go-confignet

The module is freely ispired by asp.net Configuration framework.
The responsibility for reading the configuration rests with one or more configuration providers.
Configuration providers read configuration data from key-value pairs (map[string]string) using a variety of configuration sources:

- Json Files
- Yaml Files
- Environment Variables
- Azure Key Vault
- Custom providers

The usage of the module consists in few simlpe steps:

- Create the struct which represents the configuration
- Create a configuration builder
- Add one or more configuration providers, each one with his configuration
- Build the configuration to retrieve a configuration struct
- Invoke the Configuration.Bind function to apply the configuration to your custom object

Configuration example

```go
type MyConfig struct {
    Obj1 SubObj
}

type SubObj struct {
    PropertyString string
    PropertyInt    int
    PropertyInt8   int8
    PropertyInt16  int16
    PropertyInt64  int64
    PropertyBool   bool
    Time           time.Time
}
```

Usage example

```go
var confBuilder confignet.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

// AddDefaultConfigurationProviders will add:
// JsonConfigurationProvider: default file name app.json
// YamlConfigurationProvider: default file name app.yaml
// EnvConfigurationProvider
// CmdLineConfigurationProvider
// KeyvaultConfigurationProvider: connection settings will be retrieved from the environment variables
confBuilder.AddDefaultConfigurationProviders()

conf := confBuilder.Build()

myCfg := MyConfig{}
conf.Bind("config", &myCfg)
```
