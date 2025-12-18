package mirage

import "net/http"

type ResponseRecorder struct {
	ResponseWriter http.ResponseWriter
	StatusCode     int
	Body           []byte
}

func (rec *ResponseRecorder) Header() http.Header {
	return rec.ResponseWriter.Header()
}

func (rec *ResponseRecorder) Write(b []byte) (int, error) {
	rec.Body = append(rec.Body, b...)
	return rec.ResponseWriter.Write(b)
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}
