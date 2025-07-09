package entities

import "github.com/google/uuid"

type TokenExchangeBody struct {
	OrganizationId uuid.UUID `json:"organization_id,omitempty"`
	Expire         int       `json:"expire,omitempty"`
	ExternalId     string    `json:"external_id,omitempty"`
	RedirectUri    string    `json:"redirect_uri,omitempty"`
}
