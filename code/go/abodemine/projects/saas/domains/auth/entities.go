package auth

import "github.com/google/uuid"

type SaasServerSessionTokenBody struct {
	OrganizationId uuid.UUID
	SessionId      uuid.UUID
}
