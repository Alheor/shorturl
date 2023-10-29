// Package config
// Parsing server configuration flags
package config

import (
	"flag"
	"os"
)

// DefaultAddr default server address
const defaultAddr = `localhost:8080`

// DefaultBaseHost default host url
const defaultBaseHost = `http://localhost:8080`

// DefaultLogLevel default loghandler level
const defaultLogLevel = `info`

// EnvAddr env variable name
const envAddr = `SHORT_URL_ADDR`

// EnvBaseHost env variable name
const envBaseHost = `SHORT_URL_BASE_HOST`

// EnvDefaultLogLevel env variable loghandler level
const envLogLevel = `LOG_LEVEl`

// Options Server options
var Options struct {

	// Server address (default: localhost:8080)
	Addr string

	// Host for urls (default: http://localhost:8080)
	BaseHost string

	LogLevel string
}

func init() {
	flag.StringVar(&Options.Addr, `a`, defaultAddr, "listening host:port")
	flag.StringVar(&Options.BaseHost, `b`, defaultBaseHost, "base host of url")
	flag.StringVar(&Options.LogLevel, `l`, defaultLogLevel, "loghandler level")
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
}
