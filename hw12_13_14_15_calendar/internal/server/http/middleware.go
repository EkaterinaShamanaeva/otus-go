package internalhttp

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
	})
}

func StartPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	text := "Hello!"
	fmt.Fprint(rw, text)
}
