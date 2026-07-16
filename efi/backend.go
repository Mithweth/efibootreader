package efi

import (
	"github.com/google/uuid"
)

type variableBackend interface {
	GetVariable(name string) (*BootVariable, error)
	IsEfi() bool
	GetBootIds() ([]uint16, error)
}

type BootVariable struct {
	Name       string
	GUID       uuid.UUID
	Attributes uint32
	Data       []byte
}

var backend = newBackend()
