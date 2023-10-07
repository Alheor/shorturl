package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfigFromDefaultValuesSuccess(t *testing.T) {

	Load()

	assert.Equal(t, DefaultAddr, Options.Addr)
	assert.Equal(t, DefaultBaseHost, Options.BaseHost)
}

func TestLoadConfigFromFlagsSuccess(t *testing.T) {

	os.Args = append(os.Args, `-a=addr_test_value`)
	os.Args = append(os.Args, `-b=base_host_test_value`)

	Load()

	assert.Equal(t, `addr_test_value`, Options.Addr)
	assert.Equal(t, `base_host_test_value`, Options.BaseHost)
}

func TestLoadConfigFromEnvSuccess(t *testing.T) {

	err := os.Setenv(EnvAddr, `addr_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(EnvBaseHost, `base_host_test_value`)
	assert.NoError(t, err)

	Load()

	assert.Equal(t, `addr_test_value`, Options.Addr)
	assert.Equal(t, `base_host_test_value`, Options.BaseHost)
}

func TestLoadConfigPrioritySuccess(t *testing.T) {

	os.Args = append(os.Args, `-a=addr_test_value_from_flags`)
	os.Args = append(os.Args, `-b=base_host_test_value_from_flags`)

	err := os.Setenv(EnvAddr, `addr_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(EnvBaseHost, `base_host_test_value`)
	assert.NoError(t, err)

	Load()

	assert.Equal(t, `addr_test_value`, Options.Addr)
	assert.Equal(t, `base_host_test_value`, Options.BaseHost)
}
