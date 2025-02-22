package main

import (
	"net/http"

	"github.com/Alheor/shorturl/internal/controller"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/urlhasher"
)

func main() {
	repository.Init()
	urlhasher.Init()

	err := http.ListenAndServe(controller.Addr, controller.GetRouter())
	if err != nil {
		panic(err)
	}
}
