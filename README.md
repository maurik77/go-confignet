# go-confignet

A Go configuration framework that reads configuration data from multiple sources using pluggable providers. Freely inspired by the ASP.NET Core Configuration framework.

## Installation

```bash
go get github.com/maurik77/go-confignet
```

Requires **Go 1.18** or later.

## Table of Contents

- [When to choose go-confignet](#when-to-choose-go-confignet)
- [Quick Start](#quick-start)
- [How It Works](#how-it-works)
- [Built-in Providers](#built-in-providers)
  - [JSON](#json)
  - [YAML](#yaml)
  - [TOML](#toml)
  - [Environment Variables](#environment-variables)
  - [Command Line Arguments](#command-line-arguments)
  - [Azure Key Vault](#azure-key-vault)
  - [Split Secrets (Shamir)](#split-secrets-shamir)
  - [AES Encryption](#aes-encryption)
- [Meta-Configuration](#meta-configuration)
- [Supported Field Types](#supported-field-types)
- [Provider Override Order](#provider-override-order)
- [Custom Providers](#custom-providers)

---

## When to choose go-confignet

go-confignet is designed around one core idea: **configuration comes from many places at once, and that should be a first-class concern — not an afterthought**.

### Multiple providers, one struct

Most configuration libraries treat layering as a convenience feature. In go-confignet it is the foundation. Every provider in the stack contributes to the final result, each with its own natural format and separator. You do not need to merge maps, resolve conflicts manually, or write glue code:

```go
confBuilder.Add(&providers.JSONConfigurationProvider{})     // base defaults from file
confBuilder.Add(&providers.YamlConfigurationProvider{})     // team-level overrides
confBuilder.Add(&providers.EnvConfigurationProvider{})      // deployment-specific values
confBuilder.Add(&providers.CmdLineConfigurationProvider{})  // operator overrides at runtime
confBuilder.Add(&providers.KeyVaultConfigurationProvider{   // secrets from Azure Key Vault
    BaseURL: "https://myvault.vault.azure.net",
})

conf := confBuilder.Build()
conf.Bind("app", &cfg) // all sources merged, last writer wins
```

Each provider speaks its own notation (`app.Database.Host`, `app__Database__Host`, `app--Database--Host`). The binder translates them all to the same struct field — no configuration on your part.

### Choose go-confignet when you need to:

**Blend many configuration sources without boilerplate.**
The stacked-provider model is the default, not an advanced feature. Adding or removing a source is one line.

**Store secrets split across multiple locations.**
The built-in Shamir Secret Sharing support lets you distribute secret shares across separate files, vaults, or services. No single location holds the full secret, and reconstruction is automatic at load time.

**Encrypt configuration files at rest.**
The AES decrypter pipeline integrates directly with any provider. Encrypt a file with the bundled `aescrypt` CLI, point the provider at it, and values are transparently decrypted during `Build()` — with the key itself optionally loaded from a separate source.

**Declare your provider stack in a config file instead of code.**
The meta-configuration system lets you list providers, their properties, and their decrypters in a `settings.json` or `settings.yaml` file. The provider stack is assembled at runtime without recompiling.

**Come from an ASP.NET Core background.**
The builder pattern, provider interface, and layered override model are directly inspired by `Microsoft.Extensions.Configuration`. The mental model transfers.

---

## Quick Start

```go
package main

import (
    "fmt"
    confignet "github.com/maurik77/go-confignet"
)

type AppConfig struct {
    Database struct {
        Host string
        Port int
    }
}

func main() {
    cfg := AppConfig{}
    confignet.Bind("app", &cfg) // reads app.json, app.yaml, env vars, and command line
    fmt.Println(cfg.Database.Host)
}
```

`confignet.Bind` uses the default provider stack (JSON → YAML → ENV → CmdLine → KeyVault). Each provider overrides values set by the previous one, so environment variables take precedence over files.

---

## How It Works

### Core flow

```
ConfigurationBuilder  →  Build()  →  Configuration  →  Bind(section, &struct)
        │                                  │
        │  .Add(provider)           iterates providers,
        │  .AddWithEncrypter(       calls filterProperties,
        │    provider, decrypter)   then binds the flat map to the struct
        │
    IConfigurationProvider.Load(decrypter)
        → populates map[string]string
```

### Separators

Each provider uses its own natural key separator. `Bind` translates between them transparently.

| Provider           | Separator | Example key                      |
|--------------------|-----------|----------------------------------|
| JSON / YAML        | `.`       | `app.Database.Host`              |
| Environment        | `__`      | `app__Database__Host`            |
| Command line       | `-`       | `app-Database-Host`              |
| Azure Key Vault    | `--`      | `app--Database--Host`            |

### Sections

`Bind` accepts a section prefix so you can bind a subsection directly:

```go
// bind the entire config
conf.Bind("app", &appCfg)

// bind just the database sub-section
dbCfg := DatabaseConfig{}
conf.Bind("app/Database", &dbCfg)
```

Use `/` as the section separator regardless of which providers are active — the framework translates it automatically.

### GetValue

Retrieve a single value without binding to a struct:

```go
host := conf.GetValue("app/Database/Host")
```

### Provider interface

A configuration provider must implement:

```go
type IConfigurationProvider interface {
    Load(decrypter IConfigurationDecrypter)
    GetData() map[string]string
    GetSeparator() string
}
```

`Load` is called once during `Build()`. It populates an internal `map[string]string` where keys use the provider's separator for hierarchy. Errors (e.g. file not found) are logged and the provider returns an empty map — the rest of the stack still works.

### Key rules

- Keys are **case-sensitive**: `Database` and `database` are different keys.
- If the same key exists in multiple providers, the **last provider wins**.
- All values are stored and returned as **strings**.

---

## Built-in Providers

### JSON

Loads configuration from a JSON file. Separator: `.`

```go
type JSONConfigurationProvider struct {
    FilePath string // default: "app.json"
}
```

**Example file:**

```json
{
  "app": {
    "PropertyInt8": 45,
    "Database": {
      "Host": "localhost",
      "Port": 5432
    }
  }
}
```

**Resulting map:**

| Key                    | Value       |
|------------------------|-------------|
| `app.PropertyInt8`     | `"45"`      |
| `app.Database.Host`    | `"localhost"` |
| `app.Database.Port`    | `"5432"`    |

**Usage:**

```go
confBuilder.Add(&providers.JSONConfigurationProvider{FilePath: "config/app.json"})
```

---

### YAML

Loads configuration from a YAML file. Separator: `.`

```go
type YamlConfigurationProvider struct {
    FilePath string // default: "app.yaml"
}
```

**Example file:**

```yaml
app:
  PropertyInt8: 45
  Database:
    Host: localhost
    Port: 5432
```

**Usage:**

```go
confBuilder.Add(&providers.YamlConfigurationProvider{FilePath: "config/app.yaml"})
```

---

### TOML

Loads configuration from a TOML file. Separator: `.`

```go
type TomlConfigurationProvider struct {
    FilePath string // default: "app.toml"
}
```

**Example file:**

```toml
[app]
PropertyInt8 = 45

[app.Database]
Host = "localhost"
Port = 5432

[[app.Items]]
Name  = "first"
Value = 10

[[app.Items]]
Name  = "second"
Value = 20
```

**Usage:**

```go
confBuilder.Add(&providers.TomlConfigurationProvider{FilePath: "config/app.toml"})
```

**Note:** TOML arrays of tables (`[[...]]`) and inline arrays are both fully supported.

---

### Environment Variables

Loads configuration from environment variables. Separator: `__` (double underscore).

```go
type EnvConfigurationProvider struct {
    Prefix       string // optional: only load vars that start with this prefix
    RemovePrefix bool   // if true, strip the prefix from the key
}
```

**Example:**

```bash
export app__Database__Host=localhost
export app__Database__Port=5432
```

**Usage:**

```go
// Load all environment variables
confBuilder.Add(&providers.EnvConfigurationProvider{})

// Load only vars prefixed with "MYAPP__", strip the prefix
confBuilder.Add(&providers.EnvConfigurationProvider{
    Prefix:       "MYAPP__",
    RemovePrefix: true,
})
```

**Array and map indexing via environment variables:**

```bash
export app__Items__0__Name=first
export app__Items__1__Name=second
export app__Tags__production=true
```

---

### Command Line Arguments

Loads configuration from command line arguments in the form `key=value` or `-key=value`. Separator: `-`.

```go
type CmdLineConfigurationProvider struct {
    Prefix       string                // optional: only load args that start with this prefix
    RemovePrefix bool                  // if true, strip the prefix from the key
    KeyMapper    func(arg string) string // optional: custom key transformation
}
```

**Example:**

```bash
./myapp -app-Database-Host=localhost -app-Database-Port=5432
```

**Resulting map:**

| Key                   | Value       |
|-----------------------|-------------|
| `app-Database-Host`   | `"localhost"` |
| `app-Database-Port`   | `"5432"`    |

**Usage:**

```go
confBuilder.Add(&providers.CmdLineConfigurationProvider{})
```

---

### Azure Key Vault

Loads secrets from Azure Key Vault. Separator: `--` (double hyphen, because Key Vault secret names only allow alphanumeric characters and hyphens).

```go
type KeyVaultConfigurationProvider struct {
    BaseURL      string // Key Vault URL, e.g. "https://myvault.vault.azure.net"
    TenantID     string // optional: Azure tenant ID for service principal auth
    ClientID     string // optional: Azure client ID for service principal auth
    ClientSecret string // optional: Azure client secret for service principal auth
    Prefix       string // optional: only load secrets that start with this prefix
    RemovePrefix bool   // if true, strip the prefix from the key
}
```

**Authentication:**
- If `TenantID`, `ClientID`, and `ClientSecret` are all set, service principal authentication is used.
- Otherwise, `DefaultAzureCredential` is used, which tries the following in order:
  - `AZURE_TENANT_ID` / `AZURE_CLIENT_ID` / `AZURE_CLIENT_SECRET`
  - `AZURE_CLIENT_CERTIFICATE_PATH`
  - `AZURE_USERNAME` / `AZURE_PASSWORD`
  - Workload identity, managed identity, Azure CLI, and more.

**Secret naming convention:**

Since Key Vault names use `-` as separator, a secret named `app--Database--Host` maps to the `Database.Host` field in the `app` section:

```
app--Database--Host  →  app / Database / Host
```

**Usage:**

```go
// Service principal
confBuilder.Add(&providers.KeyVaultConfigurationProvider{
    BaseURL:      "https://myvault.vault.azure.net",
    TenantID:     "00000000-...",
    ClientID:     "00000000-...",
    ClientSecret: os.Getenv("KV_SECRET"),
})

// DefaultAzureCredential (recommended for production)
confBuilder.Add(&providers.KeyVaultConfigurationProvider{
    BaseURL: "https://myvault.vault.azure.net",
})
```

---

### Split Secrets (Shamir)

Implements [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing): a secret is split into N shares stored in separate files or vaults. Any K-of-N shares are sufficient to reconstruct it. No single share reveals the secret.

This is implemented using `ChainedConfigurationProvider` (which collects multiple values for the same key) combined with `ShamirConfigurationDecrypter`.

**Generating shares** (using the `go-shamir` library directly or any compatible tool):

```go
import "github.com/lafriks/go-shamir"

shares, _ := shamir.Split([]byte("my secret value"), 3, 2) // 3 shares, 2 required
```

Store each share (base64-encoded) in a separate file or vault under the same key name.

**Usage (2-of-3 shares from two files):**

```go
var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
var chained extensions.IChainedConfigurationProvider = &confignet.ChainedConfigurationProvider{}

chained.Add(&providers.YamlConfigurationProvider{FilePath: "shares/share-1.yaml"})
chained.Add(&providers.JSONConfigurationProvider{FilePath: "shares/share-2.json"})
// share-3 is kept separately; any two of the three are sufficient

confBuilder.AddWithEncrypter(chained, &decrypters.ShamirConfigurationDecrypter{})
conf := confBuilder.Build()
```

**Share file format** (`share-1.yaml`):

```yaml
config:
  Obj1:
    PropertyString: "<base64-encoded share>"
```

---

### AES Encryption

Encrypts configuration values at rest using AES-256. Values are encrypted with the `aescrypt` CLI tool and decrypted at runtime by `AesConfigurationDecrypter`.

#### Step 1 — Encrypt a config file

Build and run the `aescrypt` CLI tool:

```bash
go build ./cli/
./aescrypt -config app.json -configType json -secret mysecretkey -dest app-encrypted.json
```

#### Step 2 — Wire the decrypter

```go
confBuilder.AddWithEncrypter(
    &providers.JSONConfigurationProvider{FilePath: "app-encrypted.json"},
    &decrypters.AesConfigurationDecrypter{Secret: "mysecretkey"},
)
```

**Loading the secret from another config source:**

To avoid hardcoding the secret, `AesConfigurationDecrypter` can fetch it from a separate configuration at init time:

```go
&decrypters.AesConfigurationDecrypter{
    ConfigFileType:   "yaml",
    ConfigFilePath:   "secrets.yaml",
    SecretConfigPath: "vault/aesKey",
}
```

---

## Meta-Configuration

Instead of wiring providers in code, you can declare them in a `settings.json` or `settings.yaml` file and load them dynamically. This is useful for changing provider configuration without recompiling.

**Loading from a file:**

```go
confBuilder.ConfigureConfigurationProviders("json", "settings.json")
// or
confBuilder.ConfigureConfigurationProviders("yaml", "settings.yaml")
```

**Loading from environment variables** (reads `confignet_configfiletype` and `confignet_configfilepath`):

```bash
export confignet_configfiletype=json
export confignet_configfilepath=settings.json
```

```go
confBuilder.ConfigureConfigurationProvidersFromEnv()
```

**`settings.json` format:**

```json
{
  "providers": [
    {
      "name": "json",
      "properties": { "filePath": "app.json" }
    },
    {
      "name": "yaml",
      "properties": { "filePath": "app.yaml" }
    },
    {
      "name": "env",
      "properties": {}
    },
    {
      "name": "cmdline",
      "properties": {}
    },
    {
      "name": "keyvault",
      "properties": {
        "baseURL": "https://myvault.vault.azure.net"
      }
    },
    {
      "name": "json",
      "properties": { "filePath": "app-encrypted.json" },
      "decrypter": {
        "name": "aes",
        "properties": { "secret": "mysecretkey" }
      }
    },
    {
      "name": "chained",
      "decrypter": { "name": "shamir" },
      "providers": [
        { "name": "yaml", "properties": { "filePath": "shares/share-1.yaml" } },
        { "name": "json", "properties": { "filePath": "shares/share-2.json" } }
      ]
    }
  ]
}
```

**Built-in provider names:**

| Name       | Provider                         |
|------------|----------------------------------|
| `json`     | JSONConfigurationProvider        |
| `yaml`     | YamlConfigurationProvider        |
| `toml`     | TomlConfigurationProvider        |
| `env`      | EnvConfigurationProvider         |
| `cmdline`  | CmdLineConfigurationProvider     |
| `keyvault` | KeyVaultConfigurationProvider    |
| `chained`  | ChainedConfigurationProvider     |

**Built-in decrypter names:**

| Name     | Decrypter                      |
|----------|-------------------------------|
| `aes`    | AesConfigurationDecrypter     |
| `shamir` | ShamirConfigurationDecrypter  |

---

## Supported Field Types

The binder maps string values to the following Go types:

| Go type                              | Notes                                          |
|--------------------------------------|------------------------------------------------|
| `string`                             |                                                |
| `int`, `int8`, `int16`, `int32`, `int64` |                                            |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` |                                       |
| `float32`, `float64`                 |                                                |
| `bool`                               | Accepts `true`, `false`, `1`, `0`, etc.        |
| `time.Time`                          | Must be RFC3339Nano format, e.g. `2006-01-02T15:04:05.999999999Z` |
| Nested structs                       |                                                |
| Slices (`[]T`)                       | Indexed with `__0__`, `__1__`, ...             |
| Fixed arrays (`[N]T`)                | Indexed with `__0__`, `__1__`, ...             |
| Maps (`map[K]V`)                     | Key inserted between separators, e.g. `__mykey__` |
| Pointers to any of the above         | Allocated automatically if nil                 |

**Slice/array example:**

```bash
# environment variables
export app__Items__0__Name=first
export app__Items__1__Name=second
```

```go
type Config struct {
    Items []Item
}
type Item struct {
    Name string
}
```

**Map example:**

```bash
export app__Scores__alice=100
export app__Scores__bob=200
```

```go
type Config struct {
    Scores map[string]int
}
```

---

## Provider Override Order

Providers are applied in the order they are added. The **last provider to set a key wins**. This is the standard pattern for environment-specific overrides:

```go
confBuilder.Add(&providers.JSONConfigurationProvider{})   // base defaults
confBuilder.Add(&providers.YamlConfigurationProvider{})   // optional overrides
confBuilder.Add(&providers.EnvConfigurationProvider{})    // deployment overrides
confBuilder.Add(&providers.CmdLineConfigurationProvider{}) // highest priority
conf := confBuilder.Build()
```

A value set via environment variable will override the same value from a JSON file. A command line argument overrides everything.

---

## Custom Providers

Implement `IConfigurationProvider` to add your own source:

```go
package myprovider

import "github.com/maurik77/go-confignet/extensions"

type MyConfigurationProvider struct {
    data map[string]string
}

func (p *MyConfigurationProvider) Load(decrypter extensions.IConfigurationDecrypter) {
    p.data = map[string]string{
        "app.MyKey": "myValue",
    }

    if decrypter != nil {
        for key, value := range p.data {
            if decrypted, err := decrypter.Decrypt(value); err == nil {
                p.data[key] = decrypted
            }
        }
    }
}

func (p *MyConfigurationProvider) GetData() map[string]string { return p.data }
func (p *MyConfigurationProvider) GetSeparator() string       { return "." }
```

**Register for meta-configuration support:**

To make your provider available by name in `settings.json`, implement `IConfigurationSource` and register it at init time:

```go
type MyConfigurationProviderSource struct{}

func (s *MyConfigurationProviderSource) GetUniqueIdentifier() string { return "myprovider" }

func (s *MyConfigurationProviderSource) NewConfigurationProvider(
    settings extensions.ProviderSettings,
) (extensions.IConfigurationProvider, error) {
    return &MyConfigurationProvider{}, nil
}

func init() {
    confignet.RegisterConfigurationSource(&MyConfigurationProviderSource{})
}
```

Once registered, use it in `settings.json`:

```json
{ "name": "myprovider", "properties": {} }
```
