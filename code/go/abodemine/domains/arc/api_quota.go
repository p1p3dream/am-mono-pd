package arc

import (
	"time"

	"github.com/google/uuid"
)

const (
	QuotaExhaustedDaily   = "daily"
	QuotaExhaustedMonthly = "monthly"
)

type ApiQuota struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Meta      map[string]any

	OrganizationId uuid.UUID
	DailyQuota     int64
	MonthlyQuota   int64

	AddressLayoutEnabled      bool
	AssessorLayoutEnabled     bool
	CompsLayoutEnabled        bool
	ListingLayoutEnabled      bool
	RecorderLayoutEnabled     bool
	RentEstimateLayoutEnabled bool
	SaleEstimateLayoutEnabled bool
}

type ApiQuotaTransaction struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Meta      map[string]any

	OrganizationId uuid.UUID
	ApiKeyId       *uuid.UUID
	TrxTimestamp   time.Time
	Description    *string

	BaseReqAmount            int32
	AddressLayoutAmount      int32
	AssessorLayoutAmount     int32
	CompsLayoutAmount        int32
	ListingLayoutAmount      int32
	RecorderLayoutAmount     int32
	RentEstimateLayoutAmount int32
	SaleEstimateLayoutAmount int32
}
