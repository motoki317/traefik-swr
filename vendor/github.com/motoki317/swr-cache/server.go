package swrcache

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/motoki317/sc"
)

type Config struct {
	TTL   string
	Grace string
}

type cacheKey struct {
	method string
	url    string
}

type server struct {
	proxy http.Handler
	cache *sc.Cache[cacheKey, *response]
}

func New(proxyTarget string, config *Config) (http.Handler, error) {
	targetURL, err := url.Parse(proxyTarget)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy target: %v", err)
	}
	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.User = targetURL.User
		req.URL.Host = targetURL.Host

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(fmt.Sprintf("error on making request: %v", err)))
			return
		}
		for k, v := range res.Header {
			rw.Header()[k] = v
		}
		rw.WriteHeader(res.StatusCode)
		_, err = io.Copy(rw, res.Body)
		if err != nil {
			_, _ = rw.Write([]byte(fmt.Sprintf("error on writing response body: %v", err)))
			return
		}
	})
	return NewWithProxy(config, proxy)
}

func NewWithProxy(config *Config, proxy http.Handler) (http.Handler, error) {
	ttl, err := time.ParseDuration(config.TTL)
	if err != nil {
		return nil, fmt.Errorf("invalid ttl: %w", err)
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("ttl needs to be a positive value")
	}
	grace, err := time.ParseDuration(config.Grace)
	if err != nil {
		return nil, fmt.Errorf("invalid grace: %w", err)
	}

	p := &server{
		proxy: proxy,
	}
	p.cache, err = sc.New(p.replace, ttl, ttl+grace, sc.WithCleanupInterval(ttl), sc.EnableStrictCoalescing())
	if err != nil {
		return nil, fmt.Errorf("failed to generate cache instance: %w", err)
	}
	return p, nil
}

var cacheableMethods = []string{
	"HEAD",
	"GET",
}

func (s *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !contains(cacheableMethods, req.Method) {
		s.proxy.ServeHTTP(rw, req)
		return
	}

	res, err := s.cache.Get(req.Context(), cacheKey{
		method: req.Method,
		url:    req.URL.String(),
	})
	if err != nil {
		log.Printf("error on making request to upstream: %v\n", err)
		return
	}

	header := rw.Header()
	for k, v := range res.headers {
		header[k] = v
	}
	rw.WriteHeader(res.status)
	_, err = rw.Write(res.body.Bytes())
	if err != nil {
		log.Printf("error on writing response: %v\n", err)
		return
	}
}

func (s *server) replace(ctx context.Context, key cacheKey) (*response, error) {
	res := newResponse()
	req, err := http.NewRequestWithContext(ctx, key.method, key.url, nil)
	if err != nil {
		return nil, err
	}
	s.proxy.ServeHTTP(res, req)
	return res, nil
}
