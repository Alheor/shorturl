package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFlags(t *testing.T) {

	os.Args = append(os.Args, `-a=addr_test_value`)
	os.Args = append(os.Args, `-b=base_host_test_value`)
	os.Args = append(os.Args, `-f=file-storage-path_test_value`)
	os.Args = append(os.Args, `-d=database-dsn_test_value`)

	Load(nil)

	assert.Equal(t, `addr_test_value`, options.Addr)
	assert.Equal(t, `base_host_test_value`, options.BaseHost)
	assert.Equal(t, `file-storage-path_test_value`, options.FileStoragePath)
	assert.Equal(t, `database-dsn_test_value`, options.DatabaseDsn)
}

func TestLoadConfigFromEnv(t *testing.T) {

	err := os.Setenv(`SERVER_ADDRESS`, `addr_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`BASE_URL`, `base_host_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`FILE_STORAGE_PATH`, `file-storage-path_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`DATABASE_DSN`, `database-dsn_test_value`)
	assert.NoError(t, err)

	Load(nil)

	assert.Equal(t, `addr_test_value`, options.Addr)
	assert.Equal(t, `base_host_test_value`, options.BaseHost)
	assert.Equal(t, `file-storage-path_test_value`, options.FileStoragePath)
	assert.Equal(t, `database-dsn_test_value`, options.DatabaseDsn)
}

func TestPriorityLoadingConfig(t *testing.T) {

	os.Args = append(os.Args, `-a=addr_test_value_from_flags`)
	os.Args = append(os.Args, `-b=base_host_test_value_from_flags`)
	os.Args = append(os.Args, `-f=file-storage-path_test_value_from_flags`)
	os.Args = append(os.Args, `-d=database-dsn_test_value_from_flags`)

	err := os.Setenv(`SERVER_ADDRESS`, `addr_test_value_from_env`)
	assert.NoError(t, err)

	err = os.Setenv(`BASE_URL`, `base_host_test_value_from_env`)
	assert.NoError(t, err)

	err = os.Setenv(`FILE_STORAGE_PATH`, `file-storage-path_test_value_from_env`)
	assert.NoError(t, err)

	err = os.Setenv(`DATABASE_DSN`, `database-dsn_test_value_from_env`)
	assert.NoError(t, err)

	Load(nil)

	assert.Equal(t, `addr_test_value_from_env`, options.Addr)
	assert.Equal(t, `base_host_test_value_from_env`, options.BaseHost)
	assert.Equal(t, `file-storage-path_test_value_from_env`, options.FileStoragePath)
	assert.Equal(t, `database-dsn_test_value_from_env`, options.DatabaseDsn)
}
