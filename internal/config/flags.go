package config

import "flag"

var Options struct {
	Addr     string
	BaseHost string
}

func ParseFlags() {
	flag.StringVar(&Options.Addr, `a`, `localhost:8080`, "server host:port")
	flag.StringVar(&Options.BaseHost, `b`, `http://localhost:8080`, "base host of url")
	flag.Parse()
}
