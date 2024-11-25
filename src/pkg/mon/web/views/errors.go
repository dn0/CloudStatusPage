package views

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/go-playground/form/v4"

	"cspage/pkg/data"
	"cspage/pkg/db"
	htp "cspage/pkg/http"
)

const (
	errSourceContext          = "context"
	errSourceHandler          = "handler"
	statusClientClosedRequest = 499
)

func Render(w http.ResponseWriter, r *http.Request, c templ.Component, err error) {
	ctx := r.Context()
	if handleErrors(ctx, w, err) {
		return
	}
	err = c.Render(ctx, w)
	_ = handleErrors(ctx, w, err)
}

func handleErrors(ctx context.Context, w http.ResponseWriter, err error) bool {
	if err != nil {
		saveError(ctx, err, errSourceHandler)
		setErrorStatus(w, err)
		return true
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		saveError(ctx, ctxErr, errSourceContext)
		setErrorStatus(w, err)
		return true
	}
	return false
}

func saveError(ctx context.Context, err error, src string) {
	ctxErrHolder, ok := ctx.Value(htp.ContextKeyErr).(*htp.CtxErrorHolder)
	if !ok {
		return
	}
	ctxErrHolder.Source = src
	ctxErrHolder.Error = err
}

func setErrorStatus(w http.ResponseWriter, err error) {
	w.Header().Del(cacheHeader)
	var errNotFound *db.ObjectNotFoundError
	var errInvalidInput *data.InvalidInputError
	var errsInvalidInput data.InvalidInputErrors
	var errsFormDecode form.DecodeErrors
	switch {
	case errors.As(err, &errNotFound), errors.Is(err, data.ErrInvalidPage):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.As(err, &errInvalidInput), errors.As(err, &errsInvalidInput):
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	case errors.As(err, &errsFormDecode):
		// Masking details about the decoding error, but we give it a different status code
		for field := range errsFormDecode {
			errsInvalidInput = append(errsInvalidInput, &data.InvalidInputError{Field: field})
		}
		http.Error(w, errsInvalidInput.Error(), http.StatusBadRequest)
	case errors.Is(err, os.ErrDeadlineExceeded), errors.Is(err, context.DeadlineExceeded):
		// NOTE: this can happen during rendering and the status header 200 was already sent
		http.Error(w, err.Error(), http.StatusRequestTimeout)
	case errors.Is(err, context.Canceled):
		http.Error(w, err.Error(), statusClientClosedRequest)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
