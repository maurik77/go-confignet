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

type strictConfig struct {
	Host    string
	Port    int
	Nested  strictNested
	Tags    map[string]string
	Items   []string
}

type strictNested struct {
	Timeout int
	Name    string
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func buildEnvConf() extensions.IConfiguration {
	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.EnvConfigurationProvider{})
	return b.Build()
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestBindStrict_NoUnknownKeys(t *testing.T) {
	t.Setenv("app__Host", "localhost")
	t.Setenv("app__Port", "8080")

	conf := buildEnvConf()
	var cfg strictConfig
	err := conf.BindStrict("app", &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)
}

func TestBindStrict_UnknownKey(t *testing.T) {
	t.Setenv("app__Hsot", "localhost") // typo: Hsot instead of Host

	conf := buildEnvConf()
	var cfg strictConfig
	err := conf.BindStrict("app", &cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Hsot")
}

func TestBindStrict_MultipleUnknownKeys(t *testing.T) {
	t.Setenv("app__Hsot", "localhost") // typo
	t.Setenv("app__Prot", "8080")     // typo

	conf := buildEnvConf()
	var cfg strictConfig
	err := conf.BindStrict("app", &cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Hsot")
	assert.Contains(t, err.Error(), "Prot")
}

func TestBindStrict_NestedValidKey(t *testing.T) {
	t.Setenv("app__Nested__Timeout", "30")

	conf := buildEnvConf()
	var cfg strictConfig
	err := conf.BindStrict("app", &cfg)
	assert.NoError(t, err)
	assert.Equal(t, 30, cfg.Nested.Timeout)
}

func TestBindStrict_NestedUnknownKey(t *testing.T) {
	t.Setenv("app__Nested__Tiemout", "30") // typo

	conf := buildEnvConf()
	var cfg strictConfig
	err := conf.BindStrict("app", &cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Tiemout")
}

func TestBindStrict_SliceIndexValid(t *testing.T) {
	t.Setenv("app__Items__0", "hello")

	conf := buildEnvConf()
	var cfg strictConfig
	err := conf.BindStrict("app", &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "hello", cfg.Items[0])
}

func TestBindStrict_MapKeyValid(t *testing.T) {
	t.Setenv("app__Tags__env", "prod")

	conf := buildEnvConf()
	var cfg strictConfig
	err := conf.BindStrict("app", &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "prod", cfg.Tags["env"])
}

func TestBindStrict_NilTarget(t *testing.T) {
	conf := buildEnvConf()
	err := conf.BindStrict("app", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil")
}

func TestBindStrict_NonPointerTarget(t *testing.T) {
	conf := buildEnvConf()
	var cfg strictConfig
	err := conf.BindStrict("app", cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non-pointer")
}
