// Package config
// Parsing server configuration flags
package config

import (
	"flag"
	"os"
)

// DefaultAddr default server address
const DefaultAddr = `localhost:8080`

// DefaultBaseHost default host url
const DefaultBaseHost = `http://localhost:8080`

// EnvAddr env variable name
const EnvAddr = `SHORT_URL_ADDR`

// EnvBaseHost env variable name
const EnvBaseHost = `SHORT_URL_BASE_HOST`

// Options Server options
var Options struct {

	// Server address (default: localhost:8080)
	Addr string

	// Host for urls (default: http://localhost:8080)
	BaseHost string
}

func init() {
	flag.StringVar(&Options.Addr, `a`, DefaultAddr, "listening host:port")
	flag.StringVar(&Options.BaseHost, `b`, DefaultBaseHost, "base host of url")
}

// Load loading config
func Load() {
	flag.Parse()

	addr, exist := os.LookupEnv(EnvAddr)
	if exist && len(addr) > 0 {
		Options.Addr = addr
	}

	baseHost, exist := os.LookupEnv(EnvBaseHost)
	if exist && len(baseHost) > 0 {
		Options.BaseHost = baseHost
	}
}
