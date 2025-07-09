package gconf

import (
	"time"

	"aidanwoods.dev/go-paseto"

	"abodemine/lib/errors"
)

const (
	Paseto_V4_SYMMETRIC  = 0
	Paseto_V4_ASYMMETRIC = 1
)

type Paseto struct {
	Type   int           `json:"type,omitempty" yaml:"type,omitempty"`
	Expire time.Duration `json:"expire,omitempty" yaml:"expire,omitempty"`
	Seed   string        `json:"seed,omitempty" yaml:"seed,omitempty"`
	Key    string        `json:"key,omitempty" yaml:"key,omitempty"`
}

type PasetoCacheItem struct {
	Parser                paseto.Parser                `json:"parser,omitempty" yaml:"parser,omitempty"`
	Expire                time.Duration                `json:"expire,omitempty" yaml:"expire,omitempty"`
	V4SymmetricKey        paseto.V4SymmetricKey        `json:"v4_symmetric_key,omitempty" yaml:"v4_symmetric_key,omitempty"`
	V4AsymmetricSecretKey paseto.V4AsymmetricSecretKey `json:"v4_asymmetric_secret_key,omitempty" yaml:"v4_asymmetric_secret_key,omitempty"`
	V4AsymmetricPublicKey paseto.V4AsymmetricPublicKey `json:"v4_asymmetric_public_key,omitempty" yaml:"v4_asymmetric_public_key,omitempty"`
}

func LoadPaseto(config *Paseto) (*PasetoCacheItem, error) {
	if config == nil {
		return nil, &errors.Object{
			Id:     "fc3bbb0f-113f-4b49-a739-080901c93513",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing configuration.",
			Path:   "/",
		}
	}

	item := &PasetoCacheItem{
		Parser: paseto.NewParser(),
		Expire: config.Expire,
	}

	switch config.Type {
	case Paseto_V4_SYMMETRIC:
		key, err := paseto.V4SymmetricKeyFromHex(config.Key)
		if err != nil {
			return nil, &errors.Object{
				Id:     "be8e7a25-bbf5-46f3-aa6d-4cf285ec156f",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to create symmetric key.",
				Cause:  err.Error(),
			}
		}

		item.V4SymmetricKey = key
	case Paseto_V4_ASYMMETRIC:
		key, err := paseto.NewV4AsymmetricSecretKeyFromSeed(config.Seed)
		if err != nil {
			return nil, &errors.Object{
				Id:     "6f1ab605-0954-49fa-be71-b6b13f80acba",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to create asymmetric key.",
				Cause:  err.Error(),
			}
		}

		item.V4AsymmetricSecretKey = key
		item.V4AsymmetricPublicKey = key.Public()
	default:
		return nil, &errors.Object{
			Id:     "bbadc4f8-ffd6-47aa-8ea4-f7d78975bca7",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "The provided Paseto type is unknown.",
			Path:   "/type",
		}
	}

	return item, nil
}
