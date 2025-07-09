package arc

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"abodemine/lib/errors"
)

type httpApiResponse struct {
	Data   any              `json:"data,omitempty"`
	Errors []*errors.Object `json:"errors,omitempty"`
}

func HttpApiDataResponse(dom Domain, w http.ResponseWriter, code int, data any) {
	HttpApiResponse(dom, w, code, data, "", nil)
}

func HttpApiErrorResponse(dom Domain, w http.ResponseWriter, requestId string, err error) {
	HttpApiResponse(dom, w, 0, nil, requestId, err)
}

func HttpApiResponse(dom Domain, w http.ResponseWriter, code int, data any, requestId string, err error) {
	w.Header().Set("Content-Type", "application/json")

	var objects []*errors.Object

	if err != nil {
		objects = errors.Sanitize(
			err,
			// dom.DeploymentEnvironment() == app.DeploymentEnvironment_LOCAL,
			false,
		).Objects

		firstError := objects[0]
		firstError.RequestId = requestId

		code = firstError.HTTPStatusCode()

		if firstError.Code == errors.Code_INTERNAL {
			// Ensure we log the full chain on internal errors.
			log.Error().
				Err(err).
				Str("label", "Internal Error").
				Str("request_id", requestId).
				Str("layer", "api_response").
				Send()
		}
	}

	out := httpApiResponse{
		Data:   data,
		Errors: objects,
	}

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Error().
			Str("id", "74603aee-4187-439b-95f2-cd98952a7964").
			Err(err).
			Msg("Failed to encode http api response.")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
