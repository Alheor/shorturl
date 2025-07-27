package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var filePath = `/tmp/config.json`

func TestLoadConfigFromFlagsWithoutFile(t *testing.T) {

	options = Options{}

	os.Args = append(os.Args, `-a=addr_test_value`)
	os.Args = append(os.Args, `-b=base_host_test_value`)
	os.Args = append(os.Args, `-f=file-storage-path_test_value`)
	os.Args = append(os.Args, `-d=database-dsn_test_value`)
	os.Args = append(os.Args, `-k=signature-key_test_value`)
	os.Args = append(os.Args, `-tlscert=cert_test_value`)
	os.Args = append(os.Args, `-tlskey=cert-key_test_value`)
	os.Args = append(os.Args, `-s=true`)

	Load()

	assert.Equal(t, `addr_test_value`, options.Addr)
	assert.Equal(t, `base_host_test_value`, options.BaseHost)
	assert.Equal(t, `file-storage-path_test_value`, options.FileStoragePath)
	assert.Equal(t, `database-dsn_test_value`, options.DatabaseDsn)
	assert.Equal(t, `signature-key_test_value`, options.SignatureKey)
	assert.Equal(t, `cert_test_value`, options.TLSCert)
	assert.Equal(t, `cert-key_test_value`, options.TLSKey)
	assert.True(t, options.EnableHTTPS)
	assert.Equal(t, ``, options.FileConfig)
}

func TestLoadConfigFromFlagsWithEmptyFile(t *testing.T) {

	options = Options{}

	f, err := os.OpenFile(filePath, os.O_CREATE, 0755)
	assert.NoError(t, err)
	f.Close()

	os.Args = append(os.Args, `-a=addr_test_value`)
	os.Args = append(os.Args, `-b=base_host_test_value`)
	os.Args = append(os.Args, `-f=file-storage-path_test_value`)
	os.Args = append(os.Args, `-d=database-dsn_test_value`)
	os.Args = append(os.Args, `-k=signature-key_test_value`)
	os.Args = append(os.Args, `-tlscert=cert_test_value`)
	os.Args = append(os.Args, `-tlskey=cert-key_test_value`)
	os.Args = append(os.Args, `-s=true`)
	os.Args = append(os.Args, `-c=`+filePath)

	Load()

	assert.Equal(t, `addr_test_value`, options.Addr)
	assert.Equal(t, `base_host_test_value`, options.BaseHost)
	assert.Equal(t, `file-storage-path_test_value`, options.FileStoragePath)
	assert.Equal(t, `database-dsn_test_value`, options.DatabaseDsn)
	assert.Equal(t, `signature-key_test_value`, options.SignatureKey)
	assert.Equal(t, `cert_test_value`, options.TLSCert)
	assert.Equal(t, `cert-key_test_value`, options.TLSKey)
	assert.True(t, options.EnableHTTPS)
	assert.Equal(t, filePath, options.FileConfig)

	err = os.Remove(filePath)
	require.NoError(t, err)
}

func TestLoadConfigFromEnvWithoutFile(t *testing.T) {

	options = Options{}

	err := os.Setenv(`SERVER_ADDRESS`, `addr_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`BASE_URL`, `base_host_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`FILE_STORAGE_PATH`, `file-storage-path_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`DATABASE_DSN`, `database-dsn_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`SIGNATURE_KEY`, `signature-key_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`TLS_CERT`, `cert_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`TLS_KEY`, `cert-key_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`ENABLE_HTTPS`, `true`)
	assert.NoError(t, err)

	os.Args = append(os.Args, `-c=`)

	Load()

	assert.Equal(t, `addr_test_value`, options.Addr)
	assert.Equal(t, `base_host_test_value`, options.BaseHost)
	assert.Equal(t, `file-storage-path_test_value`, options.FileStoragePath)
	assert.Equal(t, `database-dsn_test_value`, options.DatabaseDsn)
	assert.Equal(t, `signature-key_test_value`, options.SignatureKey)
	assert.Equal(t, `cert_test_value`, options.TLSCert)
	assert.Equal(t, `cert-key_test_value`, options.TLSKey)
	assert.True(t, options.EnableHTTPS)
	assert.Equal(t, ``, options.FileConfig)
}

func TestLoadConfigFromEnvWithEmptyFile(t *testing.T) {

	options = Options{}

	f, err := os.OpenFile(filePath, os.O_CREATE, 0755)
	assert.NoError(t, err)
	f.Close()

	err = os.Setenv(`SERVER_ADDRESS`, `addr_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`BASE_URL`, `base_host_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`FILE_STORAGE_PATH`, `file-storage-path_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`DATABASE_DSN`, `database-dsn_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`SIGNATURE_KEY`, `signature-key_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`TLS_CERT`, `cert_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`TLS_KEY`, `cert-key_test_value`)
	assert.NoError(t, err)

	err = os.Setenv(`ENABLE_HTTPS`, `true`)
	assert.NoError(t, err)

	err = os.Setenv(`CONFIG`, filePath)
	assert.NoError(t, err)

	Load()

	assert.Equal(t, `addr_test_value`, options.Addr)
	assert.Equal(t, `base_host_test_value`, options.BaseHost)
	assert.Equal(t, `file-storage-path_test_value`, options.FileStoragePath)
	assert.Equal(t, `database-dsn_test_value`, options.DatabaseDsn)
	assert.Equal(t, `signature-key_test_value`, options.SignatureKey)
	assert.Equal(t, `cert_test_value`, options.TLSCert)
	assert.Equal(t, `cert-key_test_value`, options.TLSKey)
	assert.True(t, options.EnableHTTPS)
	assert.Equal(t, filePath, options.FileConfig)

	err = os.Remove(filePath)
	require.NoError(t, err)
}

func TestLoadFromFileWithoutFile(t *testing.T) {

	options = Options{}

	options.Addr = `Addr value not changed`
	options.BaseHost = `BaseHost value not changed`
	options.FileStoragePath = `FileStoragePath value not changed`
	options.DatabaseDsn = `DatabaseDsn value not changed`
	options.SignatureKey = `SignatureKey value not changed`
	options.TLSCert = `TLSCert value not changed`
	options.TLSKey = `TLSKey value not changed`
	options.EnableHTTPS = true
	options.FileConfig = ``

	err := loadFromFile(&options)
	assert.NoError(t, err)

	assert.Equal(t, `Addr value not changed`, options.Addr)
	assert.Equal(t, `BaseHost value not changed`, options.BaseHost)
	assert.Equal(t, `FileStoragePath value not changed`, options.FileStoragePath)
	assert.Equal(t, `DatabaseDsn value not changed`, options.DatabaseDsn)
	assert.Equal(t, `SignatureKey value not changed`, options.SignatureKey)
	assert.Equal(t, `TLSCert value not changed`, options.TLSCert)
	assert.Equal(t, `TLSKey value not changed`, options.TLSKey)
	assert.Equal(t, true, options.EnableHTTPS)
}

func TestLoadFromFileWithEmptyFile(t *testing.T) {

	options = Options{}

	f, err := os.OpenFile(filePath, os.O_CREATE, 0755)
	assert.NoError(t, err)
	f.Close()

	options = Options{}

	options.Addr = `Addr value not changed`
	options.BaseHost = `BaseHost value not changed`
	options.FileStoragePath = `FileStoragePath value not changed`
	options.DatabaseDsn = `DatabaseDsn value not changed`
	options.SignatureKey = `SignatureKey value not changed`
	options.TLSCert = `TLSCert value not changed`
	options.TLSKey = `TLSKey value not changed`
	options.EnableHTTPS = true
	options.FileConfig = filePath

	err = loadFromFile(&options)
	assert.NoError(t, err)

	assert.Equal(t, `Addr value not changed`, options.Addr)
	assert.Equal(t, `BaseHost value not changed`, options.BaseHost)
	assert.Equal(t, `FileStoragePath value not changed`, options.FileStoragePath)
	assert.Equal(t, `DatabaseDsn value not changed`, options.DatabaseDsn)
	assert.Equal(t, `SignatureKey value not changed`, options.SignatureKey)
	assert.Equal(t, `TLSCert value not changed`, options.TLSCert)
	assert.Equal(t, `TLSKey value not changed`, options.TLSKey)
	assert.Equal(t, true, options.EnableHTTPS)

	err = os.Remove(filePath)
	require.NoError(t, err)
}

func TestLoadFromFileFileWithConfig(t *testing.T) {

	options = Options{}

	configStr := `{
    "server_address": "Addr value is changed",
    "base_url": "BaseHost value is changed",
    "file_storage_path": "FileStoragePath value is changed",
    "database_dsn": "DatabaseDsn value is changed",
    "signature_key": "SignatureKey value is changed",
    "tls_cert": "TLSCert value is changed",
    "tls_key": "TLSKey value is changed",
    "enable_https": true
} `

	err := os.WriteFile(filePath, []byte(configStr), 0755)
	assert.NoError(t, err)

	options.FileConfig = filePath

	err = loadFromFile(&options)
	assert.NoError(t, err)

	assert.Equal(t, `Addr value is changed`, options.Addr)
	assert.Equal(t, `BaseHost value is changed`, options.BaseHost)
	assert.Equal(t, `FileStoragePath value is changed`, options.FileStoragePath)
	assert.Equal(t, `DatabaseDsn value is changed`, options.DatabaseDsn)
	assert.Equal(t, `SignatureKey value is changed`, options.SignatureKey)
	assert.Equal(t, `TLSCert value is changed`, options.TLSCert)
	assert.Equal(t, `TLSKey value is changed`, options.TLSKey)
	assert.True(t, options.EnableHTTPS)

	err = os.Remove(filePath)
	require.NoError(t, err)
}

func TestPriorityLoadingConfig(t *testing.T) {

	configStr := `{
    "server_address": "Addr first value",
    "base_url": "BaseHost first value",
    "file_storage_path": "FileStoragePath first value",
    "database_dsn": "DatabaseDsn first value",
    "signature_key": "SignatureKey first value",
    "tls_cert": "TLSCert first value",
    "tls_key": "TLSKey first value",
    "enable_https": true
} `

	err := os.WriteFile(filePath, []byte(configStr), 0755)
	assert.NoError(t, err)

	options.FileConfig = filePath

	os.Args = append(os.Args, `-a=addr_test_value_from_flags`)
	os.Args = append(os.Args, `-b=base_host_test_value_from_flags`)
	os.Args = append(os.Args, `-f=file-storage-path_test_value_from_flags`)
	os.Args = append(os.Args, `-d=database-dsn_test_value_from_flags`)
	os.Args = append(os.Args, `-k=signature_key_test_value_from_flags`)
	os.Args = append(os.Args, `-tlscert=tls_cert_test_value_from_flags`)
	os.Args = append(os.Args, `-tlskey=tls_key_test_value_from_flags`)

	err = os.Setenv(`SERVER_ADDRESS`, `addr_test_value_from_env`)
	assert.NoError(t, err)

	err = os.Setenv(`BASE_URL`, `base_host_test_value_from_env`)
	assert.NoError(t, err)

	err = os.Setenv(`FILE_STORAGE_PATH`, `file-storage-path_test_value_from_env`)
	assert.NoError(t, err)

	err = os.Setenv(`DATABASE_DSN`, `database-dsn_test_value_from_env`)
	assert.NoError(t, err)

	err = os.Setenv(`TLS_CERT`, `tls_cert_test_value_from_env`)
	assert.NoError(t, err)

	err = os.Setenv(`TLS_KEY`, `tls_key_test_value_from_env`)
	assert.NoError(t, err)

	err = os.Setenv(`SIGNATURE_KEY`, `signature_key_test_value_from_env`)
	assert.NoError(t, err)

	Load()

	assert.Equal(t, `addr_test_value_from_env`, options.Addr)
	assert.Equal(t, `base_host_test_value_from_env`, options.BaseHost)
	assert.Equal(t, `file-storage-path_test_value_from_env`, options.FileStoragePath)
	assert.Equal(t, `database-dsn_test_value_from_env`, options.DatabaseDsn)
	assert.Equal(t, `signature_key_test_value_from_env`, options.SignatureKey)
	assert.Equal(t, `tls_cert_test_value_from_env`, options.TLSCert)
	assert.Equal(t, `tls_key_test_value_from_env`, options.TLSKey)

	err = os.Remove(filePath)
	require.NoError(t, err)
}
