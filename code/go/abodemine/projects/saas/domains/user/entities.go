package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Meta      map[string]any

	OrganizationId uuid.UUID
	Username       string
	Email          string
	RoleId         uuid.UUID
	RoleName       string

	ExternalId string
}
