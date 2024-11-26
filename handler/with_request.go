package handler

import (
	"log/slog"
	"net/http"
)

// FuncWithRequest is an HTTP handler that takes a generic input + an HTTP request, and returns a generic output.
type FuncWithRequest[Out any] func(*http.Request) (Out, error)

// WithRequest takes a FuncWithRequest and uses it to create an HTTP handler.
func WithRequest[Out any](logger *slog.Logger, f FuncWithRequest[Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("HTTP request", "http.method", r.Method, "http.url", r.URL)

		// Call out to target function.
		out, err := f(r)

		// Serve response.
		writeResponse(w, logger, out, err)
	})
}
