package endpoint

import (
	"net/http"
	"time"
	"log"
)

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		diff := time.Since(start)
		log.Printf(
			"%s\t%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			diff,
			r.RemoteAddr,
		)
		err := conf.MetricProvider.PutResponseTime("API", float64(diff) / float64(time.Millisecond))
		if err != nil {
			log.Printf("Could not publish metric: %q", err)
		}
	})
}
