package tests

import (
	"testing"

	confignet "github.com/maurik77/go-confignet"
	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/providers"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Test structs
// ---------------------------------------------------------------------------

type defaultConfig struct {
	Port    int     `default:"8080"`
	Host    string  `default:"localhost"`
	Debug   bool    `default:"true"`
	Rate    float64 `default:"1.5"`
	Nested  defaultNested
	PtrSub  *defaultNested
}

type defaultNested struct {
	Timeout int    `default:"30"`
	Name    string `default:"worker"`
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func buildEmptyConf() extensions.IConfiguration {
	// EnvConfigurationProvider with a prefix that will never match anything,
	// so the provider supplies no keys — defaults must fill the gaps.
	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.EnvConfigurationProvider{Prefix: "zzznomatch__"})
	return b.Build()
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestDefaultTag_AppliedWhenMissing(t *testing.T) {
	conf := buildEmptyConf()
	var cfg defaultConfig
	err := conf.Bind("defaults", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, true, cfg.Debug)
	assert.InDelta(t, 1.5, cfg.Rate, 1e-9)
}

func TestDefaultTag_NestedStruct(t *testing.T) {
	conf := buildEmptyConf()
	var cfg defaultConfig
	err := conf.Bind("defaults", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, 30, cfg.Nested.Timeout)
	assert.Equal(t, "worker", cfg.Nested.Name)
}

func TestDefaultTag_OverriddenByProvider(t *testing.T) {
	// The JSON file has defaults__Port = 9090 and defaults__Host = "example.com"
	t.Setenv("defaults__Port", "9090")
	t.Setenv("defaults__Host", "example.com")

	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.EnvConfigurationProvider{})
	conf := b.Build()

	var cfg defaultConfig
	err := conf.Bind("defaults", &cfg)
	assert.NoError(t, err)

	// Provider values win over tag defaults
	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, "example.com", cfg.Host)
	// Fields not set by provider still get their defaults
	assert.Equal(t, true, cfg.Debug)
	assert.InDelta(t, 1.5, cfg.Rate, 1e-9)
}

func TestDefaultTag_ProviderExplicitZeroOverridesDefault(t *testing.T) {
	// Provider explicitly sets Port to 0 — this should win over the default.
	t.Setenv("defaults__Port", "0")

	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.EnvConfigurationProvider{})
	conf := b.Build()

	var cfg defaultConfig
	err := conf.Bind("defaults", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, 0, cfg.Port)
}

func TestDefaultTag_NoDefaultTag_KeepsZero(t *testing.T) {
	// A struct without default tags should still have zero values when no provider matches.
	conf := buildEmptyConf()
	var cfg struct {
		Port int
		Host string
	}
	err := conf.Bind("notag", &cfg)
	assert.NoError(t, err)
	assert.Equal(t, 0, cfg.Port)
	assert.Equal(t, "", cfg.Host)
}

func TestDefaultTag_PtrSubStruct_NilStaysNil(t *testing.T) {
	// PtrSub is a *defaultNested — when no provider sets any of its fields,
	// it should remain nil (we don't allocate just to apply defaults).
	conf := buildEmptyConf()
	var cfg defaultConfig
	err := conf.Bind("defaults", &cfg)
	assert.NoError(t, err)
	assert.Nil(t, cfg.PtrSub)
}
