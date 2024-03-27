package restapi

import (
	"encoding/json"
	"net/http"

	httpHelper "github.com/ovrrtd/openidea-bank/internal/helper/http"
	"github.com/ovrrtd/openidea-bank/internal/model/request"
)

func (api *Restapi) Register(w http.ResponseWriter, r *http.Request) {
	var request request.Register

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		httpHelper.ResponseJSONHTTP(w, http.StatusBadRequest, "Error parsing request body", nil, nil, err)
		return
	}

	ret, code, err := api.service.Register(r.Context(), request)
	api.debugError(err)
	httpHelper.ResponseJSONHTTP(w, code, "User registered successfully", ret, nil, err)
}

func (api *Restapi) Login(w http.ResponseWriter, r *http.Request) {
	var request request.Login
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		httpHelper.ResponseJSONHTTP(w, http.StatusBadRequest, "Error parsing request body", nil, nil, err)
		return
	}

	ret, code, err := api.service.Login(r.Context(), request)
	api.debugError(err)
	httpHelper.ResponseJSONHTTP(w, code, "User logged successfully", ret, nil, err)

}
