package tests

import (
	"testing"

	confignet "github.com/maurik77/go-confignet"
	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/providers"
	"github.com/stretchr/testify/assert"
)

// arrayConfig is a self-contained struct for array syntax tests.
type arrayConfig struct {
	Strings  []string
	Ints     []int
	Items    []arrayItem
	ItemsPtr []*arrayItem
	Fixed    [3]int
}

type arrayItem struct {
	Name  string
	Value int
}

// TestArraySyntax_Slice_Primitive verifies: Strings__0 = "hello"
func TestArraySyntax_Slice_Primitive(t *testing.T) {
	t.Setenv("arr__Strings__0", "hello")
	t.Setenv("arr__Strings__1", "world")

	conf := buildEnvConf(t)
	cfg := arrayConfig{}
	assert.Nil(t, conf.Bind("arr", &cfg))

	assert.Equal(t, "hello", cfg.Strings[0])
	assert.Equal(t, "world", cfg.Strings[1])
}

// TestArraySyntax_Slice_Object verifies: Items__0__Name = "first"
func TestArraySyntax_Slice_Object(t *testing.T) {
	t.Setenv("arr__Items__0__Name", "first")
	t.Setenv("arr__Items__0__Value", "10")
	t.Setenv("arr__Items__1__Name", "second")
	t.Setenv("arr__Items__1__Value", "20")

	conf := buildEnvConf(t)
	cfg := arrayConfig{}
	assert.Nil(t, conf.Bind("arr", &cfg))

	assert.Equal(t, "first", cfg.Items[0].Name)
	assert.Equal(t, 10, cfg.Items[0].Value)
	assert.Equal(t, "second", cfg.Items[1].Name)
	assert.Equal(t, 20, cfg.Items[1].Value)
}

// TestArraySyntax_Slice_ObjectPtr verifies: ItemsPtr__0__Name = "first" (slice of pointers)
func TestArraySyntax_Slice_ObjectPtr(t *testing.T) {
	t.Setenv("arr__ItemsPtr__0__Name", "ptr-first")
	t.Setenv("arr__ItemsPtr__0__Value", "99")

	conf := buildEnvConf(t)
	cfg := arrayConfig{}
	assert.Nil(t, conf.Bind("arr", &cfg))

	assert.Equal(t, "ptr-first", cfg.ItemsPtr[0].Name)
	assert.Equal(t, 99, cfg.ItemsPtr[0].Value)
}

// TestArraySyntax_FixedArray_Primitive verifies: Fixed__1 = "42" (fixed-size array)
func TestArraySyntax_FixedArray_Primitive(t *testing.T) {
	t.Setenv("arr__Fixed__0", "10")
	t.Setenv("arr__Fixed__1", "42")
	t.Setenv("arr__Fixed__2", "30")

	conf := buildEnvConf(t)
	cfg := arrayConfig{}
	assert.Nil(t, conf.Bind("arr", &cfg))

	assert.Equal(t, 10, cfg.Fixed[0])
	assert.Equal(t, 42, cfg.Fixed[1])
	assert.Equal(t, 30, cfg.Fixed[2])
}

// TestArraySyntax_JSON verifies the same patterns via JSON (separator: ".")
func TestArraySyntax_JSON(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	confBuilder.Add(&providers.JSONConfigurationProvider{FilePath: "app-array-syntax.json"})
	conf := confBuilder.Build()

	cfg := arrayConfig{}
	assert.Nil(t, conf.Bind("arr", &cfg))

	assert.Equal(t, "hello", cfg.Strings[0])
	assert.Equal(t, "world", cfg.Strings[1])
	assert.Equal(t, "first", cfg.Items[0].Name)
	assert.Equal(t, 10, cfg.Items[0].Value)
	assert.Equal(t, "second", cfg.Items[1].Name)
	assert.Equal(t, 20, cfg.Items[1].Value)
	assert.Equal(t, 42, cfg.Fixed[1])
}

func buildEnvConf(t *testing.T) extensions.IConfiguration {
	t.Helper()
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	confBuilder.Add(&providers.EnvConfigurationProvider{})
	return confBuilder.Build()
}
