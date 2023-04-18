# go-confignet

The Configuration Framework is a module that provides a way to read configuration data from a variety of sources using configuration providers. This framework is freely inspired by the asp.net Configuration framework.

## Configuration Providers

Configuration providers read configuration data from key-value pairs (map[string]string) using a variety of configuration sources, including:

- [Json Files](#json)
- [Yaml Files](#yaml)
- [Environment Variables](#environment-variables)
- [Command line arguments](#command-line-arguments)
- [Azure Key Vault](#azure-key-vault)
- [Splitted Secrets](#splitted-secrets)
- Custom providers

## Usage

Using the Configuration Framework is simple and can be broken down into a few simple steps:

1. Create a struct that represents the configuration.
2. Create a configuration builder.
3. Add one or more configuration providers, each with its own configuration.
4. Build the configuration to retrieve a configuration struct.
5. Invoke the Configuration.Bind function to apply the configuration to your custom object.

## Getting Started

Simple configuration struct example

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

type subObjItem struct {
    PropertyString string
    PropertyInt    int
    PropertyBool   bool
}
```

Complex configuration struct example

```go
type myConfig struct {
    Obj1         *subObj
    PropertyInt8 *int8
}

type subObj struct {
    PropertyString string
    PropertyInt    int
    PropertyInt8   int8
    PropertyInt16  int16
    PropertyInt64  int64
    PropertyBool   bool
    Time           time.Time
    ArrayStr       []string
    ArrayInt       *[3]int
    ArrayObj       []subObjItem
    ArrayObjPtr    []*subObjItem
    MapStr         map[string]string
    MapInt         map[int]int
    MapObj         map[int]subObjItem
    MapObjPtr      map[bool]*subObjItem
}

type subObjItem struct {
    PropertyString string
    PropertyInt    int
    PropertyBool   bool
}
```

Basic usage example:

```go
// Default configuration providers:
// Bind applies the configuration to the given object using the default configuration providers(AddDefaultConfigurationProviders)

myCfg := MyConfig{}
confignet.Bind("config", &myCfg)
```

Intermediate usage example:

```go
var confBuilder confignet.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

// AddDefaultConfigurationProviders will add:
// 1. JsonConfigurationProvider: default file name app.json
// 2. YamlConfigurationProvider: default file name app.yaml
// 3. EnvConfigurationProvider
// 4. CmdLineConfigurationProvider
// 5. KeyVaultConfigurationProvider: connection settings will be retrieved from the environment variables:
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

The "Load" function is invoked by the configuration builder when the "build" function is called. The function loads the configuration and stores the information in a map of type map[string]string. The key of the map contains the configuration name, the value of the map the value of the configuration.
Usually the execution of the function should be safe, it means that should never throw an error. Yaml and Json configuration providers, for instance, write in the standard error stream the reason if they cannot find the file or if they are not able to read it correctly.

Configuration keys:

- Are case-sensitive. For example, ConnectionString and connectionstring are treated as different keys.
- If a key and value is set in more than one configuration providers, the value from the last provider added is used.
- Hierarchical configuration must be represented in flat way in the map. The key will be made concatenating all chained properties with a specific separator. The separator is specific of every configuration provider.

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

| OOP dotted notation (myConfig struct) | Map Key                   | Map Value    |
| ------------------------------------- | ------------------------- | ------------ |
| myConfig.PropertyString               | **PropertyString**        | "text"       |
| myConfig.PropertyInt                  | **PropertyInt**           | "3"          |
| myConfig.PropertyBool                 | **PropertyBool**          | "true"       |
| myConfig.Time                         | **Time**                  | "2022-01-01" |
| myConfig.Obj1.PropertyString          | **Obj1:PropertyString**   | "text2"      |
| myConfig.Obj1.PropertyInt             | **Obj1:PropertyInt**      | "55"         |
| myConfig.Obj1.PropertyBool            | **Obj1:PropertyBool**     | "false"      |
| myConfig.Obj1.Time                    | **Obj1:Time**             | "2022-05-01" |
| myConfig.Obj1.Obj2.PropertyInt        | **Obj1:Obj2:PropertyInt** | "33"         |

### GetData function

"GetData" function must return the map populated by the "Load" function

### GetSeparator function

"GetSeparator" function must return the separator used by the configuration provider.

## Built-in Configuration Providers

### Json

JSONConfigurationProvider loads configuration from JSON file. It uses "." (dot) as separator for hierarchical configuration. It exposes just one public property: the file path. If the path is not provided the default value "app.json" is used.

```go
// JSONConfigurationProvider loads configuration from JSON file key-value pairs
type JSONConfigurationProvider struct {
   FilePath string
}
```

Example:

```json
{
  "config": {
    "PropertyInt8": 45,
    "Obj1": {
      "PropertyString": "TestObj1",
      "PropertyInt": 1,
      "PropertyBool": true,
      "Time": "2022-01-01"
    }
  }
}
```

Map:

| Map Key                        | Map Value    |
| ------------------------------ | ------------ |
| **config.PropertyInt8**        | "45"         |
| **config.Obj1.PropertyString** | "TestObj1"   |
| **config.Obj1.PropertyInt**    | "1"          |
| **config.Obj1.PropertyBool**   | "true"       |
| **config.Obj1.Time**           | "2022-01-01" |

### Yaml

YamlConfigurationProvider loads configuration from YAML file. It uses "." (dot) as separator for hierarchical configuration. It exposes just one public property: the file path. If the path is not provided the default value "app.yaml" is used.

```go
// YamlConfigurationProvider loads configuration from YAML file key-value pairs
type YamlConfigurationProvider struct {
  FilePath string
}
```

Example:

```yaml
config:
  PropertyInt8: 45
  Obj1:
    PropertyString: "TestObj1"
    PropertyInt: 1
    PropertyBool: true
```

Map:

| Map Key                        | Map Value  |
| ------------------------------ | ---------- |
| **config.PropertyInt8**        | "45"       |
| **config.Obj1.PropertyString** | "TestObj1" |
| **config.Obj1.PropertyInt**    | "1"        |
| **config.Obj1.PropertyBool**   | "true"     |

### Environment variables

EnvConfigurationProvider loads configuration from environment variables. It uses "\_\_" (double underscore) as separator for hierarchical configuration. It exposes the following properties that change the behavior of the provider.

```go
// EnvConfigurationProvider loads configuration from environment variables
type EnvConfigurationProvider struct {
  Prefix       string
  RemovePrefix bool
}
```

Properties:

- Prefix (optional): if set only the environment variables starting with the Prefix value will be loaded. E.g.:
  - Prefix value: "secrets"
  - Environment Variables:
    - secrets\_\_cred\_\_password=pwd123
    - cred\_\_username=test1
  - Map: only secrets\_\_cred\_\_password is loaded
    - Key: secrets\_\_cred\_\_password
    - Value: pwd123
- RemovePrefix (optional): It is ignored if Prefix is not set. If set to true the environment variable will be added to the map removing from the key the prefix value and the first separator. E.g. using as example the previous settings:
  - RemovePrefix: true
  - Map: only secrets\_\_cred\_\_password is loaded
    - Key: cred\_\_password
    - Value: pwd123

Example:

```bash
# Environment variables
export config__PropertyInt8=45
export config__Obj1__PropertyString=TestObj1
export config__Obj1__PropertyInt=1
export config__Obj1__PropertyBool=true
```

Map:

| Map Key                              | Map Value  |
| ------------------------------------ | ---------- |
| **config\_\_PropertyInt8**           | "45"       |
| **config\_\_Obj1\_\_PropertyString** | "TestObj1" |
| **config\_\_Obj1\_\_PropertyInt**    | "1"        |
| **config\_\_Obj1\_\_PropertyBool**   | "true"     |

### Command line arguments

CmdConfigurationProvider loads configuration from the command line arguments. It uses "-" (hyphen) as separator for hierarchical configuration. It exposes the following properties that change the behavior of the provider.

```go
// EnvConfigurationProvider loads configuration from environment variables
type CmdLineConfigurationProvider struct {
 Prefix       string
 RemovePrefix bool
}
```

Properties:

- Prefix (optional): if set only the environment variables starting with the Prefix value will be loaded. E.g.:
  - Prefix value: "secrets"
  - Command line arguments:
    - secrets-cred-password=pwd123
    - cred-username=test1
  - Map: only secrets-cred-password is loaded
    - Key: secrets-cred-password
    - Value: pwd123
- RemovePrefix (optional): It is ignored if Prefix is not set. If set to true the environment variable will be added to the map removing from the key the prefix value and the first separator. E.g. using as example the previous settings:
  - RemovePrefix: true
  - Map: only secrets-cred-password is loaded
    - Key: cred-password
    - Value: pwd123

Example:

```bash
# Command line arguments
[Exe path] -config-PropertyInt8 45 -config-Obj1-PropertyString TestObj1 -config-Obj1-PropertyInt 1 -config-Obj1-PropertyBool true
```

Map:

| Map Key                        | Map Value  |
| ------------------------------ | ---------- |
| **config-PropertyInt8**        | "45"       |
| **config-Obj1-PropertyString** | "TestObj1" |
| **config-Obj1-PropertyInt**    | "1"        |
| **config-Obj1-PropertyBool**   | "true"     |

### Azure Key Vault

KeyVaultConfigurationProvider loads configuration from an Azure KeyVault service. Azure KeyVault is a cloud service for securely storing and accessing secrets; for further information follow the official documentation [official link](https://learn.microsoft.com/en-Us/azure/key-vault/general/basic-concepts).
It uses "\-\-_" (double hyphen) as separator for hierarchical configuration. It exposes the following properties that change the behavior of the provider.

```go
// KeyVaultConfigurationProvider loads configuration from Azure Key Vault
type KeyVaultConfigurationProvider struct {
 Prefix       string
 RemovePrefix bool
 TenantID     string
 ClientID     string
 ClientSecret string
 BaseURL      string
}
```

Properties:

- Prefix (optional): if set only the environment variables starting with the Prefix value will be loaded. E.g.:
  - Prefix value: "secrets"
  - Environment Variables:
    - secrets--cred--password=pwd123
    - cred--username=test1
  - Map: only secrets--cred--password is loaded
    - Key: secrets--cred--password
    - Value: pwd123
- RemovePrefix (optional): It is ignored if Prefix is not set. If set to true the environment variable will be added to the map removing from the key the prefix value and the first separator. E.g. using as example the previous settings:
  - RemovePrefix: true
  - Map: only secrets--cred--password is loaded
    - Key: cred--password
    - Value: pwd123
- TenantID:
- ClientID:
- ClientSecret
