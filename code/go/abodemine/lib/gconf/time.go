package gconf

import (
	"time"

	"abodemine/lib/errors"
)

func LoadDuration(s string) (time.Duration, error) {
	v, err := time.ParseDuration(s)
	if err != nil {
		return 0, &errors.Object{
			Id:    "d21fcc9e-b0dd-4a5c-a896-711599b3b62e",
			Code:  errors.Code_UNKNOWN,
			Cause: err.Error(),
		}
	}

	return v, nil
}
