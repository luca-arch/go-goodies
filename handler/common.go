// package handler provides a set of helpers to create HTTP handlers.
// Inspired by https://www.willem.dev/articles/generic-http-handlers/
package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrInvalidArg   = errors.New("invalid query argument")
	ErrInvalidInput = errors.New("invalid input")
	success         = &SuccessResponse{Ok: true} //nolint:gochecknoglobals
)

type ErrResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Ok bool `json:"ok"`
}

func writeErrResponse(w http.ResponseWriter, err error, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(&ErrResponse{Error: err.Error()}) //nolint:wrapcheck
}

// writeResponse is an helper that writes JSON-encoded data into the ResponseWriter.
func writeResponse[T any](w http.ResponseWriter, logger *slog.Logger, out T, err error) {
	w.Header().Set("Content-Type", "application/json")

	var wErr error

	switch {
	case err == nil:
		w.WriteHeader(http.StatusOK)
		wErr = json.NewEncoder(w).Encode(out)
	case errors.Is(err, ErrInvalidInput), errors.Is(err, ErrInvalidArg):
		wErr = writeErrResponse(w, err, http.StatusBadRequest)
	default:
		wErr = writeErrResponse(w, err, http.StatusInternalServerError)
	}

	if wErr != nil {
		logger.Warn("failed to serve HTTP response", "error", wErr)
	}
}
