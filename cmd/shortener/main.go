package main

import (
	"net/http"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/router"
	"github.com/Alheor/shorturl/internal/urlhasher"
)

func main() {
	config.Load()
	urlhasher.Init()

	err := http.ListenAndServe(config.GetOptions().Addr, router.GetRoutes())
	if err != nil {
		panic(err)
	}
}
