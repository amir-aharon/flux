package mirage

import "net/http"

type CachedResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}
