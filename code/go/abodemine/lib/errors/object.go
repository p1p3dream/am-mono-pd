package errors

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type Object struct {
	// The request ID that caused this error, if present.
	// Should be only populated on service layers (http handlers, etc).
	RequestId string `json:"requestId,omitempty"`

	// A globally unique identifier.
	Id string `json:"id,omitempty"`

	// The code this error should be mapped to.
	Code int `json:"code,omitempty"`

	// A label that can be used as a sub error type, or to better identify the error.
	Label string `json:"label,omitempty"`

	// The path to the error.
	// Can be used to identify the field in the request that caused the error.
	Path string `json:"path,omitempty"`

	// A short description of the error.
	Title string `json:"title,omitempty"`

	// A detailed description of the error.
	Detail string `json:"detail,omitempty"`

	// The original error message, from stdlib or external libraries/packages.
	Cause string `json:"cause,omitempty"`

	// Additional metadata about the error.
	Meta map[string]any `json:"meta,omitempty"`
}

func (o *Object) Error() string {
	if o == nil {
		return ""
	}

	b := new(strings.Builder)

	b.WriteString("id:" + o.Id)

	if o.Code > 0 {
		b.WriteString(" code:" + strconv.Itoa(o.Code))
	}

	if o.Label != "" {
		b.WriteString(" label:" + o.Label)
	}

	if o.Detail != "" {
		b.WriteString(" detail:" + o.Detail)
	}

	if o.Cause != "" {
		b.WriteString(" cause:" + o.Cause)
	}

	if o.Meta != nil {
		m, err := json.Marshal(o.Meta)
		if err == nil {
			b.WriteString(" meta:" + string(m))
		} else {
			log.Error().
				Err(err).
				Str("go-value", fmt.Sprintf("%+v", o.Meta)).
				Msg("Failed to marshal error meta.")
		}
	}

	return b.String()
}

// HTTPStatusCode returns the HTTP status code for this error.
func (o *Object) HTTPStatusCode() int {
	if o == nil {
		return 0
	}

	switch o.Code {
	case Code_CANCELED:
		return 499
	case Code_UNKNOWN:
		return 500
	case Code_INVALID_ARGUMENT:
		return 400
	case Code_DEADLINE_EXCEEDED:
		return 504
	case Code_NOT_FOUND:
		return 404
	case Code_ALREADY_EXISTS:
		return 409
	case Code_PERMISSION_DENIED:
		return 403
	case Code_RESOURCE_EXHAUSTED:
		return 429
	case Code_FAILED_PRECONDITION:
		return 400
	case Code_ABORTED:
		return 409
	case Code_OUT_OF_RANGE:
		return 400
	case Code_UNIMPLEMENTED:
		return 501
	case Code_INTERNAL:
		return 500
	case Code_UNAVAILABLE:
		return 503
	case Code_DATA_LOSS:
		return 500
	case Code_UNAUTHENTICATED:
		return 401
	default:
		return o.Code
	}
}

func Internal(id string) *Object {
	return &Object{
		Id:     id,
		Code:   Code_INTERNAL,
		Title:  "Internal error",
		Detail: "An internal application error has occurred. Please contact support for assistance.",
	}
}
