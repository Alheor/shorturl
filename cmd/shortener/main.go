package main

import (
	"net/http"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/router"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"go.uber.org/zap"
)

func main() {
	config.Load()
	urlhasher.Init()
	logger.Init(nil)

	logger.Info("Starting server", zap.String("addr", config.GetOptions().Addr))
	err := http.ListenAndServe(config.GetOptions().Addr, router.GetRoutes())
	if err != nil {
		panic(err)
	}
}
