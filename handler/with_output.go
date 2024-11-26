package handler

import (
	"context"
	"log/slog"
	"net/http"
)

// FuncWithOutput is an HTTP handler that takes no input and returns a generic output.
type FuncWithOutput[Out any] func(context.Context) (Out, error)

// WithOutput takes a FuncWithOutput and uses it to create an HTTP handler.
func WithOutput[Out any](logger *slog.Logger, f FuncWithOutput[Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("HTTP request", "http.method", r.Method, "http.url", r.URL)

		// Call out to target function.
		out, err := f(r.Context())

		// Serve response.
		writeResponse(w, logger, out, err)
	})
}
