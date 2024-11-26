package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// FuncWithArgsInput is an HTTP handler that takes a generic querystring input and returns a monad.
type FuncWithArgsInput[Args any, In any] func(context.Context, Args, In) error

// WithArgsInput takes a FuncWithArgsInput and uses it to create an HTTP handler that reads the request's querystring and body.
func WithArgsInput[Args any, In any](logger *slog.Logger, f FuncWithArgsInput[Args, In]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			args Args
			in   In
			err  error
		)

		logger.Debug("HTTP request", "http.method", r.Method, "http.url", r.URL)

		args, err = InputFromRequest[Args](r)
		if err != nil {
			//nolint:errcheck // We don't care about this error.
			writeErrResponse(w, err, http.StatusBadRequest)

			return
		}

		err = json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			//nolint:errcheck // We don't care about this error.
			writeErrResponse(w, err, http.StatusBadRequest)

			return
		}

		// Call out to target function.
		err = f(r.Context(), args, in)
		if err != nil {
			//nolint:errcheck // We don't care about this error.
			writeErrResponse(w, err, http.StatusInternalServerError)
		}

		// Serve response.
		writeResponse(w, logger, success, err)
	})
}
