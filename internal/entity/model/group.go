package entity

import "github.com/google/uuid"

type Group struct {
	ID        uuid.UUID
	Name      string
	CreatedBy uuid.UUID
}
