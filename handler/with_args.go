package handler

import (
	"context"
	"log/slog"
	"net/http"
)

// FuncWithArgs is an HTTP handler that takes a generic querystring input and returns an error.
type FuncWithArgs[Args any] func(context.Context, Args) error

// WithArgs takes a FuncWithArgs and uses it to create an HTTP handler that reads the request's querystring.
func WithArgs[Args any](logger *slog.Logger, f FuncWithArgs[Args]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			in  Args
			err error
		)

		logger.Debug("HTTP request", "http.method", r.Method, "http.url", r.URL)

		in, err = InputFromRequest[Args](r)
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
