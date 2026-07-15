package efi

import "github.com/google/uuid"

type DevicePathNode struct {
	Type    uint8
	SubType uint8
	Data    []byte
}

type HardDriveNode struct {
	PartitionNumber uint32
	PartitionStart  uint64
	PartitionSize   uint64
	Signature       uuid.UUID
	MBRType         uint8
	SignatureType   uint8
}

type BootEntry struct {
	Attributes     uint32
	Description    string
	FilePathLength uint16
	DevicePath     []DevicePathNode
	OptionalData   []byte
}

type Variable struct {
	Name       string
	GUID       uuid.UUID
	Attributes uint32
	Data       []byte
}
