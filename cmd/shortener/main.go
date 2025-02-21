package main

import (
	"net/http"

	"github.com/Alheor/shorturl/internal/controller"
	"github.com/Alheor/shorturl/internal/repository"
)

func main() {

	repository.Init()

	err := http.ListenAndServe(controller.Addr, controller.GetRouter())
	if err != nil {
		panic(err)
	}
}
