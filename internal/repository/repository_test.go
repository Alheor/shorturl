package repository

import (
	"context"
	"errors"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

const targetURL = `https://practicum.yandex.ru/`
const shortName = `tstName`

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

func TestStorageIsReadyMapSuccess(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `` //режим без записи в файл
	r := Init()

	assert.True(t, r.StorageIsReady())
}

func TestAddBatchMapSuccess(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `` //режим без записи в файл

	r := Init()

	testData := getTestData()

	list := make([]BatchEl, 0, len(testData))
	for i, v := range testData {
		list = append(list, BatchEl{OriginalURL: v, ShortURL: i})
	}

	err := r.AddBatch(list)
	require.NoError(t, err)

	for index, val := range testData {
		URL, err := r.Get(index)
		require.NoError(t, err)
		assert.Equal(t, val, URL)
	}
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

func TestCreatedFileSuccess(t *testing.T) {

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

func TestLoadFromFileSuccess(t *testing.T) {

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

func TestStorageIsReadyFileSuccess(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `/tmp/test.json`
	removeFile(config.Options.FileStoragePath)
	r := Init()

	assert.True(t, r.StorageIsReady())
}

func TestAddBatchFileSuccess(t *testing.T) {
	config.Load()
	config.Options.FileStoragePath = `/tmp/test.json`
	removeFile(config.Options.FileStoragePath)
	r := Init()

	testData := getTestData()

	list := make([]BatchEl, 0, len(testData))
	for i, v := range testData {
		list = append(list, BatchEl{OriginalURL: v, ShortURL: i})
	}

	err := r.AddBatch(list)
	require.NoError(t, err)

	r = nil
	r = Init()

	for index, val := range testData {
		URL, err := r.Get(index)
		require.NoError(t, err)
		assert.Equal(t, val, URL)
	}

	removeFile(config.Options.FileStoragePath)
}

func TestCreateDBSchemaSuccess(t *testing.T) {
	config.Load()

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	if config.Options.DatabaseDsn == `` {
		t.Skip(`Run with database only`)
		return
	}

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, config.Options.DatabaseDsn)
	require.NoError(t, err)

	_, err = conn.Exec(ctx, `DROP TABLE IF EXISTS `+strings.ToLower(tableName))
	require.NoError(t, err)

	createDBSchema(ctx, conn)

	var tableExists bool
	row := conn.QueryRow(ctx, `SELECT true FROM pg_tables WHERE tablename = $1`, strings.ToLower(tableName))
	err = row.Scan(&tableExists)
	require.NoError(t, err)

	assert.Equal(t, true, tableExists)
}

func TestAddURLAndGetURLDBSuccess(t *testing.T) {
	config.Load()

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	if config.Options.DatabaseDsn == `` {
		t.Skip(`Run with database only`)
		return
	}

	r := Init()

	err := prepareDB()
	require.NoError(t, err)

	err = r.Add(shortName, targetURL)
	require.NoError(t, err)

	url, err := r.Get(shortName)
	require.NoError(t, err)

	assert.Equal(t, targetURL, url)
}

func TestGetURLNotExistDBSuccess(t *testing.T) {
	config.Load()

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	if config.Options.DatabaseDsn == `` {
		t.Skip(`Run with database only`)
		return
	}

	r := Init()

	err := prepareDB()
	require.NoError(t, err)

	url, err := r.Get(shortName)
	require.Error(t, err)

	assert.Equal(t, ``, url)
}

func TestAddURLAndGetURLDBUniqueError(t *testing.T) {
	config.Load()

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	if config.Options.DatabaseDsn == `` {
		t.Skip(`Run with database only`)
		return
	}

	r := Init()

	err := prepareDB()
	require.NoError(t, err)

	err = r.Add(shortName, targetURL)
	require.NoError(t, err)

	err = r.Add(shortName+`1`, targetURL)
	assert.Error(t, err)

	var uErr *UniqueError

	assert.True(t, errors.As(err, &uErr))
	assert.Equal(t, shortName, uErr.ShortKey)

	err = r.Add(shortName, targetURL+`1`)
	assert.Error(t, err)

	assert.True(t, errors.As(err, &uErr))
	assert.Equal(t, shortName, uErr.ShortKey)
}

func TestRemoveURLDBSuccess(t *testing.T) {
	config.Load()

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	if config.Options.DatabaseDsn == `` {
		t.Skip(`Run with database only`)
		return
	}

	r := Init()

	err := prepareDB()
	require.NoError(t, err)

	err = r.Add(shortName, targetURL)
	require.NoError(t, err)

	r.Remove(shortName)

	url, err := r.Get(shortName)
	require.Error(t, err)

	assert.Equal(t, ``, url)
}

func TestStorageIsReadyDBSuccess(t *testing.T) {
	config.Load()

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	if config.Options.DatabaseDsn == `` {
		t.Skip(`Run with database only`)
		return
	}

	r := Init()

	err := prepareDB()
	require.NoError(t, err)

	assert.True(t, r.StorageIsReady())
}

func TestAddURLAndGetURLBatchDBSuccess(t *testing.T) {
	config.Load()

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	if config.Options.DatabaseDsn == `` {
		t.Skip(`Run with database only`)
		return
	}

	r := Init()

	err := prepareDB()
	require.NoError(t, err)

	testData := getTestData()

	list := make([]BatchEl, 0, len(testData))
	for i, v := range testData {
		list = append(list, BatchEl{OriginalURL: v, ShortURL: i})
	}

	err = r.AddBatch(list)
	require.NoError(t, err)

	for index, val := range testData {
		URL, err := r.Get(index)
		require.NoError(t, err)
		assert.Equal(t, val, URL)
	}
}

func TestAddURLAndGetURLBatchDBUniqueError(t *testing.T) {
	config.Load()

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	if config.Options.DatabaseDsn == `` {
		t.Skip(`Run with database only`)
		return
	}

	r := Init()

	err := prepareDB()
	require.NoError(t, err)

	testData := getTestData()
	testData[shortName+`6`] = targetURL + `5`

	list := make([]BatchEl, 0, len(testData))
	for i, v := range testData {
		list = append(list, BatchEl{OriginalURL: v, ShortURL: i})
	}

	err = r.AddBatch(list)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), ErrValueAlreadyExist)
}

func prepareDB() error {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.Options.DatabaseDsn)

	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, `TRUNCATE `+strings.ToLower(tableName))
	if err != nil {
		return err
	}

	return nil
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
