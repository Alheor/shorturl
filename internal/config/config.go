// Package config
// Parsing server configuration flags
package config

import (
	"flag"
	"os"
)

const defaultAddr = `localhost:8080`
const defaultBaseHost = `http://localhost:8080`
const defaultLogLevel = `info`
const defaultLFileStoragePath = `/tmp/short-url-db.json`
const envAddr = `SHORT_URL_ADDR`
const envBaseHost = `SHORT_URL_BASE_HOST`
const envFileStoragePath = `FILE_STORAGE_PATH`
const envLogLevel = `LOG_LEVEl`
const envDatabaseDsn = `DATABASE_DSN`

// Options Server options
var Options struct {
	Addr            string
	BaseHost        string
	LogLevel        string
	FileStoragePath string
	DatabaseDsn     string
}

func init() {
	flag.StringVar(&Options.Addr, `a`, defaultAddr, "listening host:port")
	flag.StringVar(&Options.BaseHost, `b`, defaultBaseHost, "base host of url")
	flag.StringVar(&Options.LogLevel, `l`, defaultLogLevel, "log handler level")
	flag.StringVar(&Options.FileStoragePath, `f`, defaultLFileStoragePath, "Path to storage file")
	flag.StringVar(&Options.DatabaseDsn, `d`, ``, "Path to storage file")
}

// Load loading config
func Load() {
	flag.Parse()

	addr, exist := os.LookupEnv(envAddr)
	if exist && len(addr) > 0 {
		Options.Addr = addr
	}

	baseHost, exist := os.LookupEnv(envBaseHost)
	if exist && len(baseHost) > 0 {
		Options.BaseHost = baseHost
	}

	logLevel, exist := os.LookupEnv(envLogLevel)
	if exist && len(logLevel) > 0 {
		Options.LogLevel = logLevel
	}

	fileStoragePath, exist := os.LookupEnv(envFileStoragePath)
	if exist {
		Options.FileStoragePath = fileStoragePath
	}

	databaseDsn, exist := os.LookupEnv(envDatabaseDsn)
	if exist {
		Options.DatabaseDsn = databaseDsn
	}
}
