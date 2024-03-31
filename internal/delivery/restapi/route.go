package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (api *Restapi) MakeRoute(mr *mux.Router) {

	mr.Use(api.middleware.LoggingMiddleware)

	// prometheus
	mr.Handle("/metrics", promhttp.Handler())
	mr.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		rJson, _ := json.Marshal(map[string]interface{}{"service": "ok"})
		w.Header().Add("content-length", strconv.Itoa(len(rJson)))
		w.Write([]byte(rJson))
	})
	// user
	api.middleware.NewRoute(mr, http.MethodPost, "/v1/user/register", api.Register)
	api.middleware.NewRoute(mr, http.MethodPost, "/v1/user/login", api.Login)
	// image
	api.middleware.NewRoute(mr, http.MethodPost, "/v1/image", api.middleware.Authentication(true)(http.HandlerFunc(api.UploadImage)))
	// balance
	api.middleware.NewRoute(mr, http.MethodGet, "/v1/balance", api.middleware.Authentication(true)(http.HandlerFunc(api.GetBalances)))
	api.middleware.NewRoute(mr, http.MethodPost, "/v1/balance", api.middleware.Authentication(true)(http.HandlerFunc(api.AddBalance)))
	api.middleware.NewRoute(mr, http.MethodGet, "/v1/balance/history", api.middleware.Authentication(true)(http.HandlerFunc(api.GetBalancesHistory)))
	// transaction
	api.middleware.NewRoute(mr, http.MethodPost, "/v1/transaction", api.middleware.Authentication(true)(http.HandlerFunc(api.CreateTransaction)))
}
