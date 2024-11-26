package handler

import (
	"context"
	"log/slog"
	"net/http"
)

// FuncWithArgsOutput is an HTTP handler that takes a generic querystring input and returns a monad.
type FuncWithArgsOutput[Args any, Out any] func(context.Context, Args) (Out, error)

// WithArgsOutput takes a FuncWithArgsOutput and uses it to create an HTTP handler that reads the request's querystring and serves a result or an error.
func WithArgsOutput[Args any, Out any](logger *slog.Logger, f FuncWithArgsOutput[Args, Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			args Args
			err  error
		)

		logger.Debug("HTTP request", "http.method", r.Method, "http.url", r.URL)

		args, err = InputFromRequest[Args](r)
		if err != nil {
			//nolint:errcheck // We don't care about this error.
			writeErrResponse(w, err, http.StatusBadRequest)

			return
		}

		// Call out to target function.
		out, err := f(r.Context(), args)

		// Serve response.
		writeResponse(w, logger, out, err)
	})
}
