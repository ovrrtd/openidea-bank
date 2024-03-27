package restapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "openidea_bank",
		Help:    "Histogram of server request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10), // Adjust bucket sizes as needed
	}, []string{"path", "method", "status"})
)

func (api *Restapi) MakeRoute(mr *mux.Router) {

	mr.Use(api.middleware.LoggingMiddleware)
	// user
	mr.HandleFunc("/v1/user/register", api.Register).Methods(http.MethodPost)
	mr.HandleFunc("/v1/user/login", api.Login).Methods(http.MethodPost)
	// image
	mr.HandleFunc("/v1/image", api.middleware.Authentication(true)(http.HandlerFunc(api.UploadImage))).Methods(http.MethodPost)
}

func NewRoute(app *echo.Echo, method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	app.Add(method, path, wrapHandlerWithMetrics(path, method, handler), middleware...)
}

func wrapHandlerWithMetrics(path string, method string, handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()

		// Execute the actual handler and catch any errors
		err := handler(c)

		// Regardless of whether an error occurred, record the metrics
		duration := time.Since(startTime).Seconds()

		requestHistogram.WithLabelValues(path, method, strconv.Itoa(c.Response().Status)).Observe(duration)
		return err
	}
}
