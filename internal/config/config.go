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

var opts Options

func init() {
	flag.StringVar(&opts.Addr, `a`, `localhost:8080`, "listen host/ip:port")
	flag.StringVar(&opts.BaseHost, `b`, `http://localhost:8080`, "base host")
	flag.StringVar(&opts.FileStoragePath, `f`, `/tmp/short-url.json`, "Path to storage file")
	flag.StringVar(&opts.DatabaseDsn, `d`, ``, "database dsn")
}

func GetOptions() Options {
	return opts
}

func Load() {
	flag.Parse()

	err := env.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	println(`--- Loaded configuration ---`)

	println(`listen: ` + opts.Addr)
	println(`base host: ` + opts.BaseHost)
	println(`file storage path: ` + opts.FileStoragePath)
	println(`database dsn: ` + opts.DatabaseDsn)
}
