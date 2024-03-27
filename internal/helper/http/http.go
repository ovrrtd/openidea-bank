package util

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ovrrtd/openidea-bank/internal/helper/common"

	"github.com/pkg/errors"
)

func ResponseJSONHTTP(w http.ResponseWriter, code int, msg string, data interface{}, meta *common.Meta, err error) {
	res := map[string]interface{}{
		"data":    data,
		"message": strings.ToLower(http.StatusText(code)),
	}
	if err != nil {
		res["message"] = errors.Cause(err).Error()
	} else {
		if msg != "" {
			res["message"] = msg
		}
	}

	if meta != nil {
		res["meta"] = meta
	}
	rJson, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(rJson)))
	w.WriteHeader(code)
	w.Write(rJson)
}

func GetJWTFromRequest(r *http.Request) string {
	// From query.
	query := r.URL.Query().Get("jwt")
	if query != "" {
		return query
	}

	// From header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}

	// From cookie.
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}

	return cookie.Value
}
