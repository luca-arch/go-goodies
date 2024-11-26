package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// FuncWith is an HTTP handler that takes a generic input with query arguments and request body, and returns a generic output.
type FuncWith[In any, Args any, Out any] func(context.Context, In, Args) (Out, error)

// HandleWithMultipleInput takes a FuncWith and uses it to create an HTTP handler that reads the request's body and the query arguments.
func With[In any, Args any, Out any](logger *slog.Logger, f FuncWith[In, Args, Out]) http.Handler {
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
		out, err := f(r.Context(), in, args)

		// Serve response.
		writeResponse(w, logger, out, err)
	})
}
