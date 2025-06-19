// Package config - сервис конфигурации.
//
// # Описание
//
// Загружает конфигурацию из параметров командной строки, либо их переменных окружения. Переменные окружения имеют приоритет загрузки.
//
// # Описание конфигурационных параметров
//
// Addr - localhost:8080.
//
// BaseHost - my-host.com.
//
// DatabaseDsn - user=app password=pass host=localhost port=5432 dbname=app, либо postgresql://app:pass@chc_postgres:5432/app.
// При конфигурации, сервис в первую очередь смотреть на параметр DatabaseDsn, и только после на FileStoragePath.
//
// FileStoragePath - отвечает за возможность сохранения данных сервиса в файл, либо в память.
// При указании пути к файлу, сервис попытается использовать указанны файл, либо создать его, если его нет.
// Для работы сервиса в режиме хранения данных в памяти, нужно установить этот параметр как пустую строку.
package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

// DefaultLSignatureKey signature key for user authentication.
const DefaultLSignatureKey = `40d40c8d1b5fff17e7edcabc6b2fa4ab`

// Options - конфигурационные параметры.
type Options struct {
	// Addr - адрес, который будет слушать сервис.
	Addr string `env:"SERVER_ADDRESS"`
	// BaseHost - хост сервиса.
	BaseHost string `env:"BASE_URL"`
	// FileStoragePath - путь к файлу для хранения данных (если сервис должен хранить данные в файле или в памяти).
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	// DatabaseDsn - Dsn базы данных (если сервис должен хранить данные в БД).
	DatabaseDsn string `env:"DATABASE_DSN"`
	// SignatureKey  - ключ подписи cookie
	SignatureKey string `env:"SIGNATURE_KEY"`
}

var options Options

func init() {
	flag.StringVar(&options.Addr, `a`, `localhost:8080`, "listen host/ip:port")
	flag.StringVar(&options.BaseHost, `b`, `http://localhost:8080`, "base host")
	flag.StringVar(&options.FileStoragePath, `f`, `/tmp/short-url.json`, "path to storage file")
	flag.StringVar(&options.DatabaseDsn, `d`, ``, "database dsn")
	flag.StringVar(&options.SignatureKey, `k`, DefaultLSignatureKey, "signature key")
}

// Load - загрузка конфигурации.
func Load() Options {

	flag.Parse()

	err := env.Parse(&options)
	if err != nil {
		log.Fatal(err)
	}

	println(`--- Loaded configuration ---`)

	println(`listen: ` + options.Addr)
	println(`base host: ` + options.BaseHost)
	println(`file storage path: ` + options.FileStoragePath)
	println(`database dsn: ` + options.DatabaseDsn)

	if options.SignatureKey == DefaultLSignatureKey {
		println(`signature key status: used default key`)
	} else {
		println(`signature key status: key specified by parameter`)
	}

	return options
}
