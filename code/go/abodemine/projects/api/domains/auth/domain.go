package auth

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"abodemine/domains/arc"
	"abodemine/domains/token"
	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/ptr"
	"abodemine/lib/val"
)

const (
	// 5 minutes.
	DefaultApiSessionTtl int32 = 300
)

type Domain interface {
	Authenticate(ctx context.Context, in *AuthenticateInput) (*AuthenticateOutput, error)

	InsertApiQuotaTransaction(r *arc.Request, in *InsertApiQuotaTransactionInput) (*InsertApiQuotaTransactionOutput, error)

	TokenExchange(r *arc.Request, in *TokenExchangeInput) (*TokenExchangeOutput, error)
}

type domain struct {
	repository Repository

	ArcDomain   arc.Domain
	TokenDomain token.Domain
}

type NewDomainInput struct {
	Repository Repository

	ArcDomain   arc.Domain
	TokenDomain token.Domain
}

func NewDomain(in *NewDomainInput) *domain {
	rep := in.Repository

	if rep == nil {
		rep = &repository{}
	}

	return &domain{
		repository:  rep,
		ArcDomain:   in.ArcDomain,
		TokenDomain: in.TokenDomain,
	}
}

type AuthenticateInput struct {
	AuthorizationHeader []string
}

type AuthenticateOutput struct {
	Request *arc.Request
}

func (dom *domain) Authenticate(ctx context.Context, in *AuthenticateInput) (*AuthenticateOutput, error) {
	if in == nil {
		return nil, &errors.Object{
			Id:     "80d44621-aa24-4ecd-959f-20944377a716",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing input.",
		}
	}

	if len(in.AuthorizationHeader) == 0 {
		return nil, &errors.Object{
			Id:     "76b06b33-0945-47b7-9bb2-5040863f6ae9",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Missing auth header.",
		}
	}

	keyHash := strings.TrimPrefix(in.AuthorizationHeader[0], "Bearer ")

	var keyType ApiKeyType

	if strings.HasPrefix(keyHash, "AM.p.") {
		keyType = ApiKeyTypeLegacy
	} else {
		return nil, &errors.Object{
			Id:     "a3428d48-31f0-427c-bf18-a0d59cfc563b",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Invalid api key.",
		}
	}

	createRequest := func(session arc.ServerSession) (*arc.Request, error) {
		request, err := dom.ArcDomain.CreateRequest(&arc.CreateRequestInput{
			Context: ctx,
			Flags:   session.Flags(),
		})
		if err != nil {
			return nil, errors.Forward(err, "a8d2e762-7ec4-4eaf-b84d-ec4f2485619d")
		}

		request.SetSession(session)

		return request, nil
	}

	systemRequest, err := dom.createSystemRequest(ctx)
	if err != nil {
		return nil, errors.Forward(err, "1d59d1c1-73b2-4d20-b54e-02506800af62")
	}

	selectApiSessionOut, err := dom.SelectApiSession(systemRequest, &SelectApiSessionInput{
		KeyHash: keyHash,
		KeyType: keyType,
	})
	if err == nil {
		if selectApiSessionOut.Invalid {
			return nil, &errors.Object{
				Id:     "97a02772-faea-4dbd-bf12-da7428cb2125",
				Code:   errors.Code_UNAUTHENTICATED,
				Detail: "Unauthorized.",
			}
		}

		switch selectApiSessionOut.QuotaExhausted {
		case arc.QuotaExhaustedDaily:
			return nil, &errors.Object{
				Id:     "4f43dcf0-1022-455b-894c-04f1a8cc1ce2",
				Code:   errors.Code_RESOURCE_EXHAUSTED,
				Label:  "DAILY_QUOTA_EXHAUSTED",
				Detail: "Not enough daily quota to complete the request.",
			}
		case arc.QuotaExhaustedMonthly:
			return nil, &errors.Object{
				Id:     "b9422b42-1d54-4869-8483-f6029b7d63a4",
				Code:   errors.Code_RESOURCE_EXHAUSTED,
				Label:  "MONTHLY_QUOTA_EXHAUSTED",
				Detail: "Not enough monthly quota to complete the request.",
			}
		}

		request, err := createRequest(selectApiSessionOut.Session)
		if err != nil {
			return nil, errors.Forward(err, "552d3969-41af-48c7-8a81-dbe9ce7cf1f4")
		}

		return &AuthenticateOutput{
			Request: request,
		}, nil
	}

	lastErr := errors.Last(err)

	// A problem other than "KEY_NOT_FOUND" occurred.
	if !strings.HasPrefix(lastErr.Id, "c0a66e0c") {
		return nil, errors.Forward(err, "3423d918-14d7-4959-8621-ba2820c0ab39")
	}

	// Check if key exists in the database.
	// We'll cache it either way to prevent hitting
	// the database again for a few minutes.

	selectApiKeyRecordOut, err := dom.repository.SelectApiKeyRecord(systemRequest, &SelectApiKeyRecordInput{
		KeyType:   keyType,
		KeyHash:   keyHash,
		KeyStatus: ApiKeyStatusActive,
	})
	if err != nil {
		return nil, errors.Forward(err, "978d503c-1473-4d78-88da-208d133048a8")
	}

	var updateApiKeyInput *UpdateApiSessionInput

	record := selectApiKeyRecordOut.Record

	if record == nil {
		updateApiKeyInput = &UpdateApiSessionInput{
			Invalid: true,
			KeyType: keyType,
			KeyHash: keyHash,
			TTL:     DefaultApiSessionTtl,
		}

		_, err = dom.UpdateApiSession(systemRequest, updateApiKeyInput)
		if err != nil {
			return nil, errors.Forward(err, "33d112fb-2bab-453b-ada4-646f82193536")
		}

		return nil, &errors.Object{
			Id:     "a0158f12-5773-4127-b8be-5cae91aec360",
			Code:   errors.Code_UNAUTHENTICATED,
			Detail: "Unauthorized.",
		}
	}

	selectApiQuotaAvailabilityOut, err := dom.repository.SelectApiQuotaAvailability(systemRequest, &SelectApiQuotaAvailabilityInput{
		OrganizationId: record.OrganizationId,
	})
	if err != nil {
		return nil, errors.Forward(err, "bd07a2d0-b5a8-4335-9b61-c33dfd80ef32")
	}

	hasDailyQuota := selectApiQuotaAvailabilityOut.HasDailyQuota
	hasMonthlyQuota := selectApiQuotaAvailabilityOut.HasMonthlyQuota

	if !hasDailyQuota || !hasMonthlyQuota {
		updateApiKeyInput = &UpdateApiSessionInput{
			QuotaExhausted: val.Ternary(
				!hasDailyQuota,
				arc.QuotaExhaustedDaily,
				arc.QuotaExhaustedMonthly,
			),
			KeyType: keyType,
			KeyHash: keyHash,
			TTL:     DefaultApiSessionTtl,
		}

		_, err = dom.UpdateApiSession(systemRequest, updateApiKeyInput)
		if err != nil {
			return nil, errors.Forward(err, "53757426-dc49-495d-ac14-76728707a1f4")
		}

		switch {
		case !hasDailyQuota:
			return nil, &errors.Object{
				Id:     "63e88aec-dafd-4592-9f48-718babf7c082",
				Code:   errors.Code_RESOURCE_EXHAUSTED,
				Label:  "DAILY_QUOTA_EXHAUSTED",
				Detail: "Not enough daily quota to complete the request.",
			}
		case !hasMonthlyQuota:
			return nil, &errors.Object{
				Id:     "856d32fb-7afe-476d-aae4-de8e185a4e04",
				Code:   errors.Code_RESOURCE_EXHAUSTED,
				Label:  "MONTHLY_QUOTA_EXHAUSTED",
				Detail: "Not enough monthly quota to complete the request.",
			}
		}
	}

	var ttl int32

	if record.ExpiresAt.IsZero() {
		ttl = DefaultApiSessionTtl
	} else {
		ttl = int32(min(
			int64(DefaultApiSessionTtl),
			// int64 to prevent overflow.
			int64(time.Until(record.ExpiresAt).Seconds()),
		))
	}

	session, err := dom.ArcDomain.CreateServerSession(
		systemRequest,
		&arc.CreateServerSessionInput{
			OrganizationId: record.OrganizationId,
			UserId:         record.UserId,
			KeyHash:        keyHash,
			KeyId:          record.Id,
			KeyType:        int16(keyType),
			RoleName:       record.RoleName,
			SessionType:    arc.SessionTypeApiServer,
			Timezone:       "UTC",
			TTL:            1,
			Flags:          selectApiQuotaAvailabilityOut.EnabledLayouts,
			DoNotSave:      true,
		},
	)
	if err != nil {
		return nil, errors.Forward(err, "9b15bf77-1bbd-4330-99f3-950032efcb65")
	}

	_, err = dom.UpdateApiSession(systemRequest, &UpdateApiSessionInput{
		KeyType: keyType,
		KeyHash: keyHash,
		Session: session,
		TTL:     ttl,
	})
	if err != nil {
		return nil, errors.Forward(err, "a10ddbca-fd8c-45a4-82bb-c50b77c008f8")
	}

	request, err := createRequest(session)
	if err != nil {
		return nil, errors.Forward(err, "50a891b0-57da-4983-a9e3-0c96f25ba20c")
	}

	out := &AuthenticateOutput{
		Request: request,
	}

	return out, nil
}

type InsertApiQuotaTransactionInput struct {
	Entity *arc.ApiQuotaTransaction
}

type InsertApiQuotaTransactionOutput struct{}

func (dom *domain) InsertApiQuotaTransaction(r *arc.Request, in *InsertApiQuotaTransactionInput) (*InsertApiQuotaTransactionOutput, error) {
	if in.Entity.Id == uuid.Nil {
		id, err := val.NewUUID7()
		if err != nil {
			return nil, errors.Forward(err, "4dd17b98-9357-495b-b551-ed14ccda21c9")
		}

		in.Entity.Id = id
	}

	in.Entity.OrganizationId = r.Session().OrganizationId()

	if in.Entity.ApiKeyId == nil {
		in.Entity.ApiKeyId = ptr.UUID(r.Session().KeyId())
	}

	insertOut, err := dom.repository.InsertApiQuotaTransactionRecord(r, &InsertApiQuotaTransactionRecordInput{
		Record: in.Entity,
	})
	if err != nil {
		return nil, errors.Forward(err, "0f22efd1-6ac1-4129-b3b1-a6a9deef07f4")
	}

	if !insertOut.HasDailyQuota {
		if insertOut.DailyQuota-insertOut.DailyUsage+insertOut.TrxLayoutSum > 0 {
			return nil, &errors.Object{
				Id:     "bafd8e03-9ed8-4b5d-9836-bd9b3aa1c7ac",
				Code:   errors.Code_RESOURCE_EXHAUSTED,
				Label:  "INSUFFICIENT_DAILY_QUOTA",
				Detail: "Not enough daily quota to complete the request.",
			}
		}

		// Only update the cached session if the daily quota is exhausted.

		systemRequest, err := dom.createSystemRequest(r.Context())
		if err != nil {
			return nil, errors.Forward(err, "48a0bc41-3d42-4065-bfc4-f57fd2414a86")
		}

		session := r.Session()
		updateApiKeyInput := &UpdateApiSessionInput{
			QuotaExhausted: arc.QuotaExhaustedDaily,
			KeyType:        ApiKeyType(session.KeyType()),
			KeyHash:        session.KeyHash(),
			TTL:            DefaultApiSessionTtl,
		}

		_, err = dom.UpdateApiSession(systemRequest, updateApiKeyInput)
		if err != nil {
			return nil, errors.Forward(err, "6d0b265c-92e0-4f9e-a850-a1cdf21c439c")
		}

		return nil, &errors.Object{
			Id:     "a020f270-6367-41e2-b5b6-849c7e7ab332",
			Code:   errors.Code_RESOURCE_EXHAUSTED,
			Label:  "DAILY_QUOTA_EXHAUSTED",
			Detail: "Not enough daily quota to complete the request.",
		}
	}

	if !insertOut.HasMonthlyQuota {
		if insertOut.MonthlyQuota-insertOut.MonthlyUsage+insertOut.TrxLayoutSum > 0 {
			return nil, &errors.Object{
				Id:     "f79b4e65-4713-432c-b0f1-b3872ea2e19b",
				Code:   errors.Code_RESOURCE_EXHAUSTED,
				Label:  "INSUFFICIENT_MONTHLY_QUOTA",
				Detail: "Not enough monthly quota to complete the request.",
			}
		}

		// Only update the cached session if the monthly quota is exhausted.

		systemRequest, err := dom.createSystemRequest(r.Context())
		if err != nil {
			return nil, errors.Forward(err, "af5a2133-9130-45c9-8538-f7100026779e")
		}

		session := r.Session()
		updateApiKeyInput := &UpdateApiSessionInput{
			QuotaExhausted: arc.QuotaExhaustedMonthly,
			KeyType:        ApiKeyType(session.KeyType()),
			KeyHash:        session.KeyHash(),
			TTL:            DefaultApiSessionTtl,
		}

		_, err = dom.UpdateApiSession(systemRequest, updateApiKeyInput)
		if err != nil {
			return nil, errors.Forward(err, "b10405c2-52fa-4cc3-9b16-cc62b7d39452")
		}

		return nil, &errors.Object{
			Id:     "43e398ef-e725-43f0-ae26-6ca121a30954",
			Code:   errors.Code_RESOURCE_EXHAUSTED,
			Label:  "MONTHLY_QUOTA_EXHAUSTED",
			Detail: "Not enough monthly quota to complete the request.",
		}
	}

	out := &InsertApiQuotaTransactionOutput{}

	return out, nil
}

func (dom *domain) createSystemRequest(ctx context.Context) (*arc.Request, error) {
	req, err := dom.ArcDomain.CreateRequest(&arc.CreateRequestInput{
		Context: ctx,
	})
	if err != nil {
		return nil, errors.Forward(err, "8e4703f3-ba61-4e7b-8c0b-091d2832a582")
	}

	systemSession, err := dom.ArcDomain.CreateServerSession(
		req,
		&arc.CreateServerSessionInput{
			OrganizationId: consts.AbodeMineOrganizationId(),
			UserId:         consts.AbodeMineBotUserId(),
			RoleName:       consts.RoleSystemAuthCheckUser,
			SessionType:    arc.SessionTypeSystem,
			TTL:            1,
			DoNotSave:      true,
		})
	if err != nil {
		return nil, errors.Forward(err, "5df875fe-2c1a-4b1e-a105-b577233f5924")
	}

	req.SetSession(systemSession)

	return req, nil
}
