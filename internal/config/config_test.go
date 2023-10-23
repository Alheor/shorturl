package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfigFromDefaultValuesSuccess(t *testing.T) {

	Load()

	assert.Equal(t, defaultAddr, Options.Addr)
	assert.Equal(t, defaultBaseHost, Options.BaseHost)
	assert.Equal(t, defaultLogLevel, Options.LogLevel)
}

func TestLoadConfigFromFlagsSuccess(t *testing.T) {

	os.Args = append(os.Args, `-a=addr_test_value`)
	os.Args = append(os.Args, `-b=base_host_test_value`)
	os.Args = append(os.Args, `-l=log_level_test_value`)

	Load()

	assert.Equal(t, `addr_test_value`, Options.Addr)
	assert.Equal(t, `base_host_test_value`, Options.BaseHost)
	assert.Equal(t, `log_level_test_value`, Options.LogLevel)
}

func TestLoadConfigFromEnvSuccess(t *testing.T) {

	err := os.Setenv(envAddr, `addr_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(envBaseHost, `base_host_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(envLogLevel, `log_level_test_value`)
	assert.NoError(t, err)

	Load()

	assert.Equal(t, `addr_test_value`, Options.Addr)
	assert.Equal(t, `base_host_test_value`, Options.BaseHost)
	assert.Equal(t, `log_level_test_value`, Options.LogLevel)
}

func TestLoadConfigPrioritySuccess(t *testing.T) {

	os.Args = append(os.Args, `-a=addr_test_value_from_flags`)
	os.Args = append(os.Args, `-b=base_host_test_value_from_flags`)
	os.Args = append(os.Args, `-l=log_level_test_value_from_flags`)

	err := os.Setenv(envAddr, `addr_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(envBaseHost, `base_host_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(envLogLevel, `log_level_test_value`)
	assert.NoError(t, err)

	Load()

	assert.Equal(t, `addr_test_value`, Options.Addr)
	assert.Equal(t, `base_host_test_value`, Options.BaseHost)
	assert.Equal(t, `log_level_test_value`, Options.LogLevel)
}
