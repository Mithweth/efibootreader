package efi

import (
	"github.com/google/uuid"
)

type variableBackend interface {
	GetVariable(name string) (*Variable, error)
	IsEfi() bool
}

type Variable struct {
	Name       string
	GUID       uuid.UUID
	Attributes uint32
	Data       []byte
}

var backend = newBackend()
