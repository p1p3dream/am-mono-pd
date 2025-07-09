// middleware/gzip.go
package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/klauspost/compress/gzhttp"
)

func GzipHandler(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(w, r, ps)
		})

		gzipHandler := gzhttp.GzipHandler(handler)

		gzipHandler.ServeHTTP(w, r)
	}
}
