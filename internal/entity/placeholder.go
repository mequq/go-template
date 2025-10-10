package entity

import "github.com/google/uuid"

type Placeholder struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
