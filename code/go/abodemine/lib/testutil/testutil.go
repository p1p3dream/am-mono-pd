package testutil

import (
	"encoding/base64"

	"github.com/google/uuid"

	"abodemine/lib/val"
)

func NewBase64UUID[K comparable](c *val.Cache[K, any], key K) string {
	u := uuid.New()
	return base64.StdEncoding.EncodeToString(c.SetGet(key, u[:]).([]byte))
}
