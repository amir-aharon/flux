package mirage

import (
	"net/http"

	"github.com/amir-aharon/flux/stash"
	"github.com/amir-aharon/flux/zen"
)

type CDN struct {
	LB    *zen.LoadBalancer
	Cache *stash.Cache[string, CachedResponse]
}

func NewCDN(backends []string, cacheCapacity int) (*CDN, error) {
	lb, err := zen.NewLoadBalancer(backends)
	if err != nil {
		return nil, err
	}
	return &CDN{
		LB:    lb,
		Cache: stash.NewCache[string, CachedResponse](cacheCapacity),
	}, nil
}

func (cdn *CDN) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		cdn.LB.ServeHTTP(rw, req)

		if req.Method == http.MethodPut || req.Method == http.MethodDelete {
			cdn.Cache.Remove(http.MethodGet + ":" + req.RequestURI)
		}
		return
	}
	key := req.Method + ":" + req.RequestURI

	if resp, found := cdn.Cache.Get(key); found {
		for k, v := range resp.Headers {
			for _, val := range v {
				rw.Header().Add(k, val)
			}
		}
		rw.WriteHeader(resp.StatusCode)
		rw.Write(resp.Body)
		return
	}

	recorder := &ResponseRecorder{
		ResponseWriter: rw,
		StatusCode:     http.StatusOK,
		Body:           make([]byte, 0),
	}
	cdn.LB.ServeHTTP(recorder, req)

	if recorder.StatusCode/100 == 2 {
		newEntry := CachedResponse{
			StatusCode: recorder.StatusCode,
			Headers:    recorder.Header().Clone(),
			Body:       recorder.Body,
		}
		cdn.Cache.Put(key, newEntry)
	}
}
