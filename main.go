package traefikswr

import (
	"context"
	"net/http"

	"github.com/motoki317/swr-cache"
)

type Config struct {
	TTL   string `json:"ttl"`
	Grace string `json:"grace"`
}

func CreateConfig() *Config {
	return &Config{}
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return swrcache.NewWithProxy(&swrcache.Config{
		TTL:   config.TTL,
		Grace: config.Grace,
	}, next)
}
