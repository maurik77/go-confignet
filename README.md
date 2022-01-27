# go-confignet

The module is freely ispired by asp.net Configuration framework.
The responsibility for reading the configuration rests with one or more configuration providers.
Configuration providers read configuration data from key-value pairs (map[string]string) using a variety of configuration sources:

- [Json Files](#json)
- [Yaml Files](#yaml)
- [Environment Variables](#environment-variables)
- [Command line arguments](#command-line-arguments)
- [Azure Key Vault](#azure-key-vault)
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

Usage example:

```go
var confBuilder confignet.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

// AddDefaultConfigurationProviders will add:
// 1. JsonConfigurationProvider: default file name app.json
// 2. YamlConfigurationProvider: default file name app.yaml
// 3. EnvConfigurationProvider
// 4. CmdLineConfigurationProvider
// 5. KeyvaultConfigurationProvider: connection settings will be retrieved from the environment variables:
//      AZURE_TENANT_ID
//      AZURE_CLIENT_ID
//      AZURE_CLIENT_SECRET
//      AZURE_CLIENT_CERTIFICATE_PATH
//      AZURE_USERNAME
//      AZURE_PASSWORD
confBuilder.AddDefaultConfigurationProviders()

conf := confBuilder.Build()

myCfg := MyConfig{}
conf.Bind("config", &myCfg)
```

## Configuration Providers

### Json

### Yaml

### Environment variables

### Command line arguments

### Azure Key Vault
