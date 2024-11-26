package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// FuncWithInput is an HTTP handler that takes a generic input and returns an error.
type FuncWithInput[In any] func(context.Context, In) error

// WithInput takes a FuncWithInput and uses it to create an HTTP handler that reads the request's body.
func WithInput[In any](logger *slog.Logger, f FuncWithInput[In]) http.Handler {
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
		err = f(r.Context(), in)

		// Serve response.
		writeResponse(w, logger, success, err)
	})
}
