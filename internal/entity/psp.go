package entity

import "github.com/google/uuid"

type PSP struct {
	ID     uuid.UUID
	Name   string
	Prefix string
}
