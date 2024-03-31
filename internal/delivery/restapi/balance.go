package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ovrrtd/openidea-bank/internal/helper/common"
	"github.com/ovrrtd/openidea-bank/internal/helper/errorer"
	httpHelper "github.com/ovrrtd/openidea-bank/internal/helper/http"
	"github.com/ovrrtd/openidea-bank/internal/model/request"
	"github.com/ovrrtd/openidea-bank/internal/model/response"
)

func (api *Restapi) AddBalance(w http.ResponseWriter, r *http.Request) {
	var payload request.AddBalance
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpHelper.ResponseJSONHTTP(w, http.StatusBadRequest, "", nil, nil, errorer.ErrInputRequest(err))
		return
	}
	ctx := r.Context()
	user, ok := ctx.Value(common.EncodedUserJwtCtxKey).(*response.User)
	api.log.Debug().Msg("masuk")

	if !ok {
		httpHelper.ResponseJSONHTTP(w, http.StatusInternalServerError, "", nil, nil, errorer.ErrInternalServer)
		return
	}

	payload.UserID = user.ID

	code, err := api.service.AddBalance(r.Context(), payload)
	httpHelper.ResponseJSONHTTP(w, code, "", nil, nil, err)
	api.debugError(err)
}

func (api *Restapi) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var payload request.CreateTransaction
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpHelper.ResponseJSONHTTP(w, http.StatusBadRequest, "", nil, nil, errorer.ErrInputRequest(err))
		return
	}
	ctx := r.Context()
	user, ok := ctx.Value(common.EncodedUserJwtCtxKey).(*response.User)
	if !ok {
		httpHelper.ResponseJSONHTTP(w, http.StatusInternalServerError, "", nil, nil, errorer.ErrInternalServer)
		return
	}

	payload.UserID = user.ID

	code, err := api.service.CreateTransaction(r.Context(), payload)
	httpHelper.ResponseJSONHTTP(w, code, "", nil, nil, err)
	api.debugError(err)
}

func (api *Restapi) GetBalances(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(common.EncodedUserJwtCtxKey).(*response.User)
	if !ok {
		httpHelper.ResponseJSONHTTP(w, http.StatusInternalServerError, "", nil, nil, errorer.ErrInternalServer)
		return
	}
	balances, code, err := api.service.GetBalances(r.Context(), user.ID)
	httpHelper.ResponseJSONHTTP(w, code, "", balances, nil, err)
	api.debugError(err)
}

func (api *Restapi) GetBalancesHistory(w http.ResponseWriter, r *http.Request) {
	var payload request.GetBalancesHistory
	query := r.URL.Query()
	if query.Has("limit") {
		limit, err := strconv.Atoi(query.Get("limit"))
		if err != nil {
			httpHelper.ResponseJSONHTTP(w, http.StatusBadRequest, "", nil, nil, errorer.ErrBadRequest)
			return
		}
		payload.Limit = limit
	} else {
		payload.Limit = 10
	}

	if query.Has("offset") {
		offset, err := strconv.Atoi(query.Get("offset"))
		if err != nil {
			httpHelper.ResponseJSONHTTP(w, http.StatusBadRequest, "", nil, nil, errorer.ErrBadRequest)
			return
		}
		payload.Offset = offset
	} else {
		payload.Offset = 0
	}

	user, ok := r.Context().Value(common.EncodedUserJwtCtxKey).(*response.User)
	if !ok {
		httpHelper.ResponseJSONHTTP(w, http.StatusInternalServerError, "", nil, nil, errorer.ErrInternalServer)
		return
	}
	payload.UserID = user.ID
	balances, code, err := api.service.GetBalancesHistory(r.Context(), payload)
	httpHelper.ResponseJSONHTTP(w, code, "", balances, nil, err)
	api.debugError(err)
}
