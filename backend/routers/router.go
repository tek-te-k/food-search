package routers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter(router *httprouter.Router) {
	router.GET("/", Index)
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("index")
}
