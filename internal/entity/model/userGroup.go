package entity

import (
	"time"

	"github.com/google/uuid"
)

type GroupRole string

const (
	AdminGroupRole GroupRole = "admin"
	UserGroupRole  GroupRole = "user"
)

type UserGroup struct {
	UserID   uuid.UUID
	GroupID  uuid.UUID
	Role     GroupRole
	JoinedAt time.Time
}
