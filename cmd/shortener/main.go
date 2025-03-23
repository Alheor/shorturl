package main

import (
	"net/http"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/router"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"go.uber.org/zap"
)

func main() {
	var err error
	config.Load()
	urlhasher.Init()

	err = logger.Init(nil)
	if err != nil {
		panic(err)
	}

	err = repository.Init()
	if err != nil {
		panic(err)
	}

	logger.Info("Starting server", zap.String("addr", config.GetOptions().Addr))
	err = http.ListenAndServe(config.GetOptions().Addr, router.GetRoutes())
	if err != nil {
		panic(err)
	}
}
