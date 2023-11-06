package repository

import (
	"github.com/Alheor/shorturl/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const targetURL = `https://practicum.yandex.ru/`
const shortName = `testShortName`

func TestAddURLAndGetURLMapSuccess(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `` //режим без записи в файл
	r := Init()

	err := r.Add(shortName, targetURL)
	require.NoError(t, err)

	url, err := r.Get(shortName)
	require.NoError(t, err)

	assert.Equal(t, targetURL, url)
}

func TestAddURLShortNameExistsMapError(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `` //режим без записи в файл
	r := Init()

	err := r.Add(shortName, targetURL)
	require.NoError(t, err)

	err = r.Add(shortName, targetURL)
	require.Error(t, err)
}

func TestAddURLURLExistsMapError(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `` //режим без записи в файл
	r := Init()

	err := r.Add(shortName, targetURL)
	require.NoError(t, err)

	err = r.Add(`otherShortName`, targetURL)
	require.Error(t, err)
}

func TestGetURLMapError(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `` //режим без записи в файл
	r := Init()

	_, err := r.Get(shortName)
	require.Error(t, err)
}

func TestAddURLAndGetURLFileSuccess(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `/tmp/test.json`
	removeFile(config.Options.FileStoragePath)

	r := Init()
	testData := getTestData()

	for index, val := range testData {
		err := r.Add(index, val)
		require.NoError(t, err)
	}

	for index, val := range testData {
		URL, err := r.Get(index)
		require.NoError(t, err)
		assert.Equal(t, val, URL)
	}

	removeFile(config.Options.FileStoragePath)
}

func TestAddURLShortNameExistsFileError(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `/tmp/test.json`
	removeFile(config.Options.FileStoragePath)

	r := Init()

	err := r.Add(shortName, targetURL)
	require.NoError(t, err)

	err = r.Add(shortName, targetURL)
	require.Error(t, err)

	removeFile(config.Options.FileStoragePath)
}

func TestAddURLURLExistsFileError(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `/tmp/test.json`
	removeFile(config.Options.FileStoragePath)

	r := Init()

	err := r.Add(shortName, targetURL)
	require.NoError(t, err)

	err = r.Add(`otherShortName`, targetURL)
	require.Error(t, err)

	removeFile(config.Options.FileStoragePath)
}

func TestGetURLFileError(t *testing.T) {
	config.Load()
	removeFile(config.Options.FileStoragePath)

	r := Init()

	_, err := r.Get(shortName)
	require.Error(t, err)
}

func TestLodCreatedSuccess(t *testing.T) {

	config.Load()
	config.Options.FileStoragePath = `/tmp/test.json`
	removeFile(config.Options.FileStoragePath)

	r := Init()
	testData := getTestData()

	for index, val := range testData {
		err := r.Add(index, val)
		require.NoError(t, err)
	}

	assert.FileExists(t, config.Options.FileStoragePath)

	removeFile(config.Options.FileStoragePath)
}

func TestLodFromFileSuccess(t *testing.T) {

	config.Load()
	config.Options.FileStoragePath = `/tmp/test.json`
	removeFile(config.Options.FileStoragePath)

	r := Init()
	testData := getTestData()

	for index, val := range testData {
		err := r.Add(index, val)
		require.NoError(t, err)
	}

	r = nil
	r = Init()

	for index, val := range testData {
		URL, err := r.Get(index)
		require.NoError(t, err)
		assert.Equal(t, val, URL)
	}

	removeFile(config.Options.FileStoragePath)
}

func removeFile(path string) {

	_, err := os.Stat(path)
	if err != nil {
		return
	}

	err = os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func getTestData() map[string]string {
	return map[string]string{
		shortName + `1`: targetURL + `1`,
		shortName + `2`: targetURL + `2`,
		shortName + `3`: targetURL + `3`,
		shortName + `4`: targetURL + `4`,
		shortName + `5`: targetURL + `5`,
	}
}
