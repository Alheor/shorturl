package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Options struct {
	Addr            string `env:"SERVER_ADDRESS"`
	BaseHost        string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
}

var options Options

func init() {
	flag.StringVar(&options.Addr, `a`, `localhost:8080`, "listen host/ip:port")
	flag.StringVar(&options.BaseHost, `b`, `http://localhost:8080`, "base host")
	flag.StringVar(&options.FileStoragePath, `f`, `/tmp/short-url.json`, "Path to storage file")
	flag.StringVar(&options.DatabaseDsn, `d`, ``, "database dsn")
}

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

	return options
}
