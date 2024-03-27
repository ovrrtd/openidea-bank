package restapi

import (
	"net/http"

	httpHelper "github.com/ovrrtd/openidea-bank/internal/helper/http"
)

func (api *Restapi) UploadImage(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		httpHelper.ResponseJSONHTTP(w, http.StatusBadRequest, "", nil, nil, err)
		return
	}
	defer file.Close()

	imgUrl, code, err := api.service.UploadImage(r.Context(), fileHeader)
	api.debugError(err)
	if err != nil {
		httpHelper.ResponseJSONHTTP(w, code, "", nil, nil, err)
		return
	}
	httpHelper.ResponseJSONHTTP(w, code, "File uploaded sucessfully", map[string]string{"imageUrl": imgUrl}, nil, err)
}
