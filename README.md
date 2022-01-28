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

## Configuration Provider

A configuration provider is responsible for reading the configuration from a specific source. It must implement the following interface:

```go
// IConfigurationProvider is configuration provider interface
type IConfigurationProvider interface {
 Load()
 GetData() map[string]string
 GetSeparator() string
}
```

### Load function

"Load" function is invoked by the configuration builder when the "build" function is called. The function loads the configuration and stores the information in a map of type map[string]string. The key of the map containes the conifguration name, the value of the map the value of the configuration.

Configuration keys:

- Are case-sensitive. For example, ConnectionString and connectionstring are treated as different keys.
- If a key and value is set in more than one configuration providers, the value from the last provider added is used.
- Hierarchical configuration must be represented in flat way in the map. The key will be made concatenating all chained properties with a specific separator. The seperator is specific of every configuration provider.

Configuration values:

- Are strings.

Let's take as example the configuration:

```go
type MyConfig struct {
  Obj1 struct {
    PropertyString string
    PropertyInt    int
    PropertyBool   bool
    Time           time.Time
    Obj2           struct {
       PropertyInt int
    }
  }
  PropertyString string
  PropertyInt    int
  PropertyBool   bool
  Time           time.Time
}
```

Let's assume that the configuration provider uses ":" as separator, the map will contain:

| Key                       | Value        |
| ------------------------- | ------------ |
| **PropertyString**        | "text"       |
| **PropertyInt**           | "3"          |
| **PropertyBool**          | "true"       |
| **Time**                  | "2022-01-01" |
| **Obj1:PropertyString**   | "text2"      |
| **Obj1:PropertyInt**      | "55"         |
| **Obj1:PropertyBool**     | "false"      |
| **Obj1:Time**             | "2022-05-01" |
| **Obj1:Obj2:PropertyInt** | "33"         |

### GetData function

"GetData" function must return the map populated by the "Load" function

### GetSeparator function

"GetSeparator" function must return the separator used by the configuration provider.

## Built-in Configuration Providers

### Json

### Yaml

### Environment variables

### Command line arguments

### Azure Key Vault
