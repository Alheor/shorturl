package config

import (
	"flag"
)

var Options struct {
	Addr     string
	BaseHost string
}

func init() {
	flag.StringVar(&Options.Addr, `a`, `localhost:8080`, "listen host/ip:port")
	flag.StringVar(&Options.BaseHost, `b`, `http://localhost:8080`, "base host")
}

func Load() {
	flag.Parse()

	println(`--- Loaded configuration ---`)
	println(`listen: ` + Options.Addr)
	println(`base host: ` + Options.BaseHost)
}
