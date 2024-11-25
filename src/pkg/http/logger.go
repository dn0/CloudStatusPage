package http

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	logLevelDefault            = slog.LevelInfo
	logLevelClientError        = slog.LevelWarn
	logLevelServerError        = slog.LevelError
	ContextKeyErr       ctxKey = "err"
)

type ctxKey string

type CtxErrorHolder struct {
	Source string
	Error  error
}

type bodyReader struct {
	io.ReadCloser
	bytes int
}

func (r *bodyReader) Read(b []byte) (int, error) {
	n, err := r.ReadCloser.Read(b)
	r.bytes += n
	//nolint:wrapcheck // Nothing to wrap here.
	return n, err
}

//nolint:funlen // There are many things to log.
func newLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == healthzURL {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			bodyRead := &bodyReader{ReadCloser: r.Body}
			r.Body = bodyRead
			writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ctxErrHolder := &CtxErrorHolder{}
			r = r.WithContext(context.WithValue(r.Context(), ContextKeyErr, ctxErrHolder))

			//nolint:contextcheck // No context here.
			defer func() {
				status := writer.Status()
				end := time.Now()

				requestAttributes := []slog.Attr{
					slog.Time("time", start),
					slog.String("proto", r.Proto),
					slog.String("method", r.Method),
					slog.String("host", r.Host),
					slog.String("path", r.URL.Path),
					slog.String("query", r.URL.RawQuery),
					slog.String("addr", r.RemoteAddr),
					slog.String("referer", r.Referer()),
					slog.String("user-agent", r.UserAgent()),
					slog.Int("length", bodyRead.bytes),
				}

				bytes := writer.BytesWritten()
				duration := end.Sub(start)
				responseAttributes := []slog.Attr{
					slog.Time("time", end),
					slog.Duration("took", duration),
					slog.Int("status", status),
					slog.Int("length", bytes),
				}

				attributes := []slog.Attr{
					{
						Key:   "request",
						Value: slog.GroupValue(requestAttributes...),
					},
					{
						Key:   "response",
						Value: slog.GroupValue(responseAttributes...),
					},
				}

				if ctxErrHolder.Error != nil {
					attributes = append(attributes, slog.Attr{
						Key: "err",
						Value: slog.GroupValue(
							slog.String("source", ctxErrHolder.Source),
							slog.String("type", fmt.Sprintf("%T", ctxErrHolder.Error)),
							slog.String("msg", ctxErrHolder.Error.Error()),
						),
					})
				}

				var level slog.Level
				switch {
				case status >= http.StatusInternalServerError:
					level = logLevelServerError
				case status >= http.StatusBadRequest, ctxErrHolder.Error != nil:
					level = logLevelClientError
				default:
					level = logLevelDefault
				}

				msg := r.Method + " " +
					r.URL.Path + " " +
					r.Proto + " " +
					strconv.Itoa(status) + " " +
					strconv.Itoa(bytes) + "B " +
					duration.String()
				logger.LogAttrs(context.Background(), level, msg, attributes...)
			}()

			next.ServeHTTP(writer, r)
		})
	}
}
