package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriterWrapper{ResponseWriter: w}
		ip := getIP(r)
		next.ServeHTTP(rw, r)

		if rw.status == 0 {
			rw.status = http.StatusOK
		}

		duration := time.Since(start)
		log.Printf("<- %s %s from=%s status=%d size=%d duration=%s",
			r.Method, r.RequestURI, ip, rw.status, rw.size, duration)
	})
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC recovered: %v\n%s", rec, debug.Stack())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
