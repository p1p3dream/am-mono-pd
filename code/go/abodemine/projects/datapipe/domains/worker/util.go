package worker

import (
	"github.com/google/uuid"

	"abodemine/lib/errors"
	"abodemine/projects/datapipe/domains/partners/abodemine"
	"abodemine/projects/datapipe/domains/partners/attom_data"
	"abodemine/projects/datapipe/domains/partners/first_american"
)

func PartnerNameById(id uuid.UUID) (string, error) {
	switch id {
	case abodemine.PartnerId:
		return "abodemine", nil
	case attom_data.PartnerId:
		return "attom-data", nil
	case first_american.PartnerId:
		return "first-american", nil
	}

	return "", &errors.Object{
		Id:     "89d413cc-7320-446b-84d0-60fba8b2eda0",
		Code:   errors.Code_INTERNAL,
		Detail: "Invalid partner id.",
	}
}
