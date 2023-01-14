package main

import (
	"github.com/lackone/go-web/framework"
	"net/http"
)

func main() {
	core := framework.NewCore()

	registerRouter(core)

	s := &http.Server{
		Addr:    ":8080",
		Handler: core,
	}

	s.ListenAndServe()
}
