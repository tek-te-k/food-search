package main

import (
	"food-search-backend/routers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// TODO DB用のdocker-composeを作成する
	// database.Connect()
	router := httprouter.New()
	routers.SetupRouter(router)

	server := http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}
	server.ListenAndServe()
}
