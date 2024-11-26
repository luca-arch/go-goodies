package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// FuncWithInputOutput is an HTTP handler that takes a generic input and returns a monad.
type FuncWithInputOutput[In any, Out any] func(context.Context, In) (Out, error)

// WithInputOutput takes a FuncWithInputOutput and uses it to create an HTTP handler that reads the request's body.
func WithInputOutput[In any, Out any](logger *slog.Logger, f FuncWithInputOutput[In, Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			in  In
			err error
		)

		logger.Debug("HTTP request", "http.method", r.Method, "http.url", r.URL)

		// Read request's body.
		err = json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			//nolint:errcheck // We don't care about this error.
			writeErrResponse(w, err, http.StatusBadRequest)

			return
		}

		// Call out to target function.
		out, err := f(r.Context(), in)

		// Serve response.
		writeResponse(w, logger, out, err)
	})
}
