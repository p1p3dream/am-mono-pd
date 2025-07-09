package middleware

import "net/http"

type ResposeWriter struct {
	StatusCode int
	http.ResponseWriter
}

func (rw *ResposeWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
