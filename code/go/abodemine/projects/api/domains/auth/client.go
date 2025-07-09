package auth

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Meta      map[string]any

	OrganizationId uuid.UUID
	Name           string
	Description    string
	RedirectUri    string
}
