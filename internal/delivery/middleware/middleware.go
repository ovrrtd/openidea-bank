package middleware

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ovrrtd/openidea-bank/internal/helper/common"
	"github.com/ovrrtd/openidea-bank/internal/helper/errorer"
	httpHelper "github.com/ovrrtd/openidea-bank/internal/helper/http"
	"github.com/ovrrtd/openidea-bank/internal/helper/jwt"
	"github.com/ovrrtd/openidea-bank/internal/service"

	"github.com/rs/zerolog"
)

type middleware struct {
	logger  zerolog.Logger
	service service.Service
}

type Middleware interface {
	Authentication(isThrowError bool) func(next http.HandlerFunc) http.HandlerFunc
	LoggingMiddleware(h http.Handler) http.Handler
	RemoveTrailingSlash(h http.Handler) http.Handler
}

func New(logger zerolog.Logger, service service.Service) Middleware {
	return &middleware{
		logger:  logger,
		service: service,
	}
}

func (m *middleware) Authentication(isThrowError bool) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.logger.Info().Msg("Authentication")
			token := httpHelper.GetJWTFromRequest(r)
			ctx := r.Context()

			if token == "" && isThrowError {
				httpHelper.ResponseJSONHTTP(w, http.StatusUnauthorized, "", nil, nil, errorer.ErrUnauthorized)
				return
			}
			if token != "" {
				claims := &common.UserClaims{}
				err := jwt.VerifyJwt(token, claims, os.Getenv("JWT_SECRET"))
				if err != nil {
					if err == errorer.ErrUnauthorized {
						httpHelper.ResponseJSONHTTP(w, http.StatusUnauthorized, "", nil, nil, errorer.ErrUnauthorized)
						return
					}
					httpHelper.ResponseJSONHTTP(w, http.StatusForbidden, "", nil, nil, errorer.ErrForbidden)
					return
				}

				usr, code, err := m.service.GetUserByID(ctx, claims.Id)
				if err != nil {
					httpHelper.ResponseJSONHTTP(w, code, "", nil, nil, err)
					return
				}
				ctx = context.WithValue(ctx, common.EncodedUserJwtCtxKey, usr)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func (m *middleware) LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				m.logger.Error().Msgf("recovering from err %v\n %s", err, buf)
				httpHelper.ResponseJSONHTTP(w, http.StatusInternalServerError, "", nil, nil, errorer.ErrInternalServer)
			}
		}()

		start := time.Now()
		wrapped := m.wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		m.logger.Info().
			Int("status", wrapped.status).
			Str("method", r.Method).
			Str("path", r.URL.EscapedPath()).
			Int64("duration", int64(time.Since(start))).
			Msg("request")
	}

	return http.HandlerFunc(fn)
}

func (m *middleware) RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (m *middleware) wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}
