package gconf

import (
	"github.com/casbin/casbin/v2"

	"abodemine/lib/errors"
)

type Casbin struct {
	Model  string `json:"model,omitempty" yaml:"model,omitempty"`
	Policy string `json:"policy,omitempty" yaml:"policy,omitempty"`
}

func LoadCasbin(config *Casbin) (*casbin.Enforcer, error) {
	if config == nil {
		return nil, &errors.Object{
			Id:     "f1302542-5386-43af-bb01-4b1e95225d88",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing config.",
		}
	}

	enforcer, err := casbin.NewEnforcer(config.Model, config.Policy)
	if err != nil {
		return nil, &errors.Object{
			Id:     "f46e84e4-40b4-4a59-8350-46d074a214a4",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create enforcer.",
			Cause:  err.Error(),
		}
	}

	return enforcer, nil
}
