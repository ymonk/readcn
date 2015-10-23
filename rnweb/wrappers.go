package main

import (
	"net/http"
	"runtime/debug"
	"time"
)

// LoggingWrapperHandler is a middleware that provides logging service
func LoggingWrapperHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		tracer.Tracef("[%s] %q %v\n", r.Method, r.URL.String(), end.Sub(start))
	}
	return http.HandlerFunc(fn)
}

// RecoverWrapperHandler is a middleware that recovers from panic
func RecoverWrapperHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				tracer.Tracef("panic: %+v\n", err)
				tracer.Trace(string(debug.Stack()))
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
