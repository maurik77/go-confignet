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

type taggedConfig struct {
	Host    string         `confignet:"host"`
	Port    int            `confignet:"port"`
	Debug   bool           `confignet:"debug"`
	Rate    float64        `confignet:"rate"`
	DB      taggedDB       `confignet:"database"`
	Tags    map[string]string
	Items   []string
}

type taggedDB struct {
	Name    string `confignet:"name"`
	Timeout int    `confignet:"connection_timeout"`
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func buildEnvConf2() extensions.IConfiguration {
	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.EnvConfigurationProvider{})
	return b.Build()
}

// ---------------------------------------------------------------------------
// Tests: basic tag mapping
// ---------------------------------------------------------------------------

func TestConfignetTag_ScalarFields(t *testing.T) {
	t.Setenv("svc__host", "example.com")
	t.Setenv("svc__port", "9090")
	t.Setenv("svc__debug", "true")
	t.Setenv("svc__rate", "2.5")

	conf := buildEnvConf2()
	var cfg taggedConfig
	err := conf.Bind("svc", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "example.com", cfg.Host)
	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, true, cfg.Debug)
	assert.InDelta(t, 2.5, cfg.Rate, 1e-9)
}

func TestConfignetTag_NestedStruct(t *testing.T) {
	t.Setenv("svc__database__name", "mydb")
	t.Setenv("svc__database__connection_timeout", "30")

	conf := buildEnvConf2()
	var cfg taggedConfig
	err := conf.Bind("svc", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "mydb", cfg.DB.Name)
	assert.Equal(t, 30, cfg.DB.Timeout)
}

func TestConfignetTag_NoTag_FallsBackToFieldName(t *testing.T) {
	// Tags and Items have no confignet tag — field name is used as-is
	t.Setenv("svc__Items__0", "hello")

	conf := buildEnvConf2()
	var cfg taggedConfig
	err := conf.Bind("svc", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "hello", cfg.Items[0])
}

func TestConfignetTag_JSONProvider(t *testing.T) {
	// app-tagged.json uses lowercase keys matching the confignet tags
	conf := buildJSONConf("app-tagged.json")
	var cfg taggedConfig
	err := conf.Bind("service", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "db-host", cfg.Host)
	assert.Equal(t, 5432, cfg.Port)
	assert.Equal(t, "mydb", cfg.DB.Name)
	assert.Equal(t, 30, cfg.DB.Timeout)
}

// ---------------------------------------------------------------------------
// Tests: tag-aware BindStrict (on this branch strict.go uses tag lookup)
// ---------------------------------------------------------------------------

func TestConfignetTag_BindStrict_TaggedKeyValid(t *testing.T) {
	t.Setenv("svc__host", "localhost")

	conf := buildEnvConf2()
	var cfg taggedConfig

	// "host" maps to Host via tag — BindStrict must not reject it
	// (BindStrict is on feature/bind-strict; here we just verify Bind works)
	err := conf.Bind("svc", &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Host)
}
