package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	confignet "github.com/maurik77/go-confignet"
	"github.com/maurik77/go-confignet/decrypters"
	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/internal"
	"github.com/maurik77/go-confignet/providers"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func buildJSONConf(filePath string) extensions.IConfiguration {
	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.JSONConfigurationProvider{FilePath: filePath})
	return b.Build()
}

// ---------------------------------------------------------------------------
// InvalidBindError.Error()
// ---------------------------------------------------------------------------

func TestInvalidBindError_Nil(t *testing.T) {
	conf := buildJSONConf("app.json")
	err := conf.Bind("config", nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "nil")
}

func TestInvalidBindError_NonPointer(t *testing.T) {
	conf := buildJSONConf("app.json")
	var cfg myConfig
	err := conf.Bind("config", cfg) // not a pointer
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "non-pointer")
}

func TestInvalidBindError_NilPointer(t *testing.T) {
	conf := buildJSONConf("app.json")
	var cfg *myConfig // nil pointer
	err := conf.Bind("config", cfg)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "nil")
}

// ---------------------------------------------------------------------------
// Configuration.GetProviders()
// ---------------------------------------------------------------------------

func TestConfiguration_GetProviders(t *testing.T) {
	p1 := &providers.JSONConfigurationProvider{FilePath: "app.json"}
	p2 := &providers.YamlConfigurationProvider{FilePath: "app.yaml"}

	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(p1)
	b.Add(p2)
	conf := b.Build()

	infos := conf.GetProviders()
	assert.Len(t, infos, 2)
	assert.Equal(t, p1, infos[0].Provider)
	assert.Equal(t, p2, infos[1].Provider)
}

// ---------------------------------------------------------------------------
// ChainedConfigurationProvider.AddWithEncrypter() & GetDataMultiValues()
// ---------------------------------------------------------------------------

func TestChainedConfigurationProvider_AddWithEncrypter(t *testing.T) {
	// Verify AddWithEncrypter registers the provider+decrypter without panicking.
	// We don't call Load here because the Shamir data is not AES-encrypted.
	chained := &confignet.ChainedConfigurationProvider{}
	dec := &decrypters.AesConfigurationDecrypter{Secret: "mysecret"}
	chained.AddWithEncrypter(&providers.JSONConfigurationProvider{FilePath: "app.json"}, dec)

	// GetProviders is not on ChainedConfigurationProvider, but GetData returns empty before Load
	data := chained.GetData()
	assert.Nil(t, data) // not loaded yet
}

func TestChainedConfigurationProvider_GetDataMultiValues(t *testing.T) {
	chained := &confignet.ChainedConfigurationProvider{}
	chained.Add(&providers.JSONConfigurationProvider{FilePath: "shamir/copy-shamir-2.json"})
	chained.Add(&providers.YamlConfigurationProvider{FilePath: "shamir/copy-shamir-3.yaml"})

	chained.Load(nil)

	multiValues := chained.GetDataMultiValues()
	assert.NotNil(t, multiValues)

	// Both files have the same key, so at least one key should have 2 values
	found := false
	for _, vals := range multiValues {
		if len(vals) == 2 {
			found = true
			break
		}
	}
	assert.True(t, found, "expected at least one key with 2 values from two shamir providers")
}

// ---------------------------------------------------------------------------
// BuildOrPanic()
// ---------------------------------------------------------------------------

func TestBuildOrPanic_Success(t *testing.T) {
	b := &confignet.ConfigurationBuilder{}
	b.Add(&providers.JSONConfigurationProvider{FilePath: "app.json"})

	assert.NotPanics(t, func() {
		conf := b.BuildOrPanic()
		assert.NotNil(t, conf)
	})
}

func TestBuildOrPanic_PanicsOnMissingFile(t *testing.T) {
	b := &confignet.ConfigurationBuilder{}
	b.Add(&providers.JSONConfigurationProvider{FilePath: "nonexistent-file-xyz.json"})

	assert.Panics(t, func() {
		b.BuildOrPanic()
	})
}

// ---------------------------------------------------------------------------
// SetLogger()
// ---------------------------------------------------------------------------

type captureLogger struct {
	messages []string
}

func (l *captureLogger) Printf(format string, args ...interface{}) {
	l.messages = append(l.messages, fmt.Sprintf(format, args...))
}

func TestSetLogger_CustomLogger(t *testing.T) {
	cl := &captureLogger{}
	confignet.SetLogger(cl)

	b := &confignet.ConfigurationBuilder{}
	b.Add(&providers.JSONConfigurationProvider{FilePath: "app.json"})
	_ = b.Build()

	// Builder logs at least one message when adding a provider
	assert.Greater(t, len(cl.messages), 0, "expected logger to receive messages")

	// Restore default logger (stdLogger is unexported, so we use a safe fallback)
	confignet.SetLogger(&captureLogger{}) // keep a non-nil logger
}

func TestSetLogger_NilIsIgnored(t *testing.T) {
	assert.NotPanics(t, func() {
		confignet.SetLogger(nil)
	})
}

// ---------------------------------------------------------------------------
// AES encryption round-trip
// ---------------------------------------------------------------------------

func TestAES_EncryptDecryptRoundTrip(t *testing.T) {
	secret := "mysupersecretkey"
	plaintext := []byte("hello, world!")

	encrypted, err := internal.EncryptBytesToBase64(plaintext, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)

	decrypted, err := internal.DecryptBase64ToBytes(encrypted, secret)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestAES_EncryptDecryptRoundTrip_LongText(t *testing.T) {
	secret := "anothersecret"
	plaintext := []byte("a longer text that exceeds the AES block size boundary definitely yes it does")

	encrypted, err := internal.EncryptBytesToBase64(plaintext, secret)
	assert.NoError(t, err)

	decrypted, err := internal.DecryptBase64ToBytes(encrypted, secret)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestAES_DecryptBadBase64(t *testing.T) {
	_, err := internal.DecryptBase64ToBytes("!!!notbase64!!!", "secret")
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// internal.MarshalToFile() / UnmarshalFromFile()
// ---------------------------------------------------------------------------

type marshalPayload struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestMarshalToFile_RoundTrip(t *testing.T) {
	src := marshalPayload{Name: "test", Value: 42}
	path := filepath.Join(t.TempDir(), "output.json")

	err := internal.MarshalToFile(path, src, json.Marshal)
	assert.NoError(t, err)

	var got marshalPayload
	err = internal.UnmarshalFromFile(path, &got, json.Unmarshal)
	assert.NoError(t, err)
	assert.Equal(t, src.Name, got.Name)
	assert.Equal(t, src.Value, got.Value)
}

func TestUnmarshalFromFile_NotFound(t *testing.T) {
	var dummy map[string]interface{}
	err := internal.UnmarshalFromFile("no-such-file.json", &dummy, json.Unmarshal)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// ---------------------------------------------------------------------------
// AesConfigurationDecrypter.Init() — path with ConfigFilePath set
// ---------------------------------------------------------------------------

func TestAesDecrypter_InitWithConfigFile(t *testing.T) {
	// AesConfigurationDecrypter.Init reads a settings file (meta-config) via
	// ConfigureConfigurationProviders, then calls GetValue(SecretConfigPath).
	// We need:
	//   1. A data file holding the actual secret value.
	//   2. A settings file pointing to that data file as a provider.

	tmpDir := t.TempDir()

	// Step 1: data file with the secret
	dataFile := filepath.Join(tmpDir, "data.json")
	err := os.WriteFile(dataFile, []byte(`{"mysecretpath": "mysecret"}`), 0600)
	assert.NoError(t, err)

	// Step 2: settings file referencing the data file
	settingsContent := fmt.Sprintf(`{"providers":[{"name":"json","properties":{"filePath":%q}}]}`, dataFile)
	settingsFile := filepath.Join(tmpDir, "settings.json")
	err = os.WriteFile(settingsFile, []byte(settingsContent), 0600)
	assert.NoError(t, err)

	dec := &decrypters.AesConfigurationDecrypter{
		ConfigFileType:   "json",
		ConfigFilePath:   settingsFile,
		SecretConfigPath: "mysecretpath",
	}

	dec.Init(&confignet.ConfigurationBuilder{})
	assert.Equal(t, "mysecret", dec.Secret)
}

// ---------------------------------------------------------------------------
// Binder error paths
// ---------------------------------------------------------------------------

func TestBinder_InvalidIntValue(t *testing.T) {
	t.Setenv("config__Obj1__PropertyInt", "notanint")

	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.EnvConfigurationProvider{})
	conf := b.Build()

	cfg := myConfig{Obj1: &subObj{}}
	err := conf.Bind("config", &cfg)
	assert.NoError(t, err)
	assert.Equal(t, 0, cfg.Obj1.PropertyInt) // parse failed, keeps zero
}

func TestBinder_InvalidBoolValue(t *testing.T) {
	t.Setenv("config__Obj1__PropertyBool", "notabool")

	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.EnvConfigurationProvider{})
	conf := b.Build()

	cfg := myConfig{Obj1: &subObj{}}
	err := conf.Bind("config", &cfg)
	assert.NoError(t, err)
	assert.False(t, cfg.Obj1.PropertyBool)
}

func TestBinder_InvalidFloatValue(t *testing.T) {
	t.Setenv("config__Obj1__PropertyFloat32", "notafloat")

	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.EnvConfigurationProvider{})
	conf := b.Build()

	cfg := myConfig{Obj1: &subObj{}}
	err := conf.Bind("config", &cfg)
	assert.NoError(t, err)
	var expected float32
	assert.Equal(t, expected, cfg.Obj1.PropertyFloat32)
}

func TestBinder_InvalidTimeValue(t *testing.T) {
	t.Setenv("config__Obj1__Time", "not-a-time")

	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.EnvConfigurationProvider{})
	conf := b.Build()

	cfg := myConfig{Obj1: &subObj{}}
	err := conf.Bind("config", &cfg)
	assert.NoError(t, err)
	assert.True(t, cfg.Obj1.Time.IsZero())
}

func TestBinder_ArrayOutOfRange(t *testing.T) {
	// ArrayInt is *[3]int — index 99 is out of range, should not panic
	t.Setenv("config__Obj1__ArrayInt__99", "42")

	var b extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	b.Add(&providers.JSONConfigurationProvider{FilePath: "app.json"})
	b.Add(&providers.EnvConfigurationProvider{})
	conf := b.Build()

	cfg := myConfig{}
	err := conf.Bind("config", &cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.Obj1)
	assert.NotNil(t, cfg.Obj1.ArrayInt)
}
