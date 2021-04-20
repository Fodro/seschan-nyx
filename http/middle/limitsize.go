package middle

import (
	"net/http"

	"go.fodro/nyx/config"
)

func LimitSize(c *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
			next.ServeHTTP(w, r)
		})
	}
}
