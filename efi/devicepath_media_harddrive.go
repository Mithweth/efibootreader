package efi

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
)

type PartitionType uint8

const (
	PartitionMBR PartitionType = 1
	PartitionGPT PartitionType = 2
)

type SignatureType uint8

const (
	SignatureNone SignatureType = iota
	SignatureMBR
	SignatureGPT
)

type HardDriveMediaNode struct {
	PartitionNumber      uint32
	PartitionSectorStart uint64
	PartitionSectorSize  uint64
	Signature            uuid.UUID
	PartitionType        PartitionType
	SignatureType        SignatureType
}

func (h *HardDriveMediaNode) String() string {
	var partitionType string
	switch h.PartitionType {
	case PartitionMBR:
		partitionType = "MBR"
	case PartitionGPT:
		partitionType = "GPT"
	default:
	}
	return fmt.Sprintf(
		"HD(%d,%s,%s,%x,%x)",
		h.PartitionNumber,
		partitionType,
		h.Signature,
		h.PartitionSectorStart,
		h.PartitionSectorSize,
	)
}

func parseHardDriveMediaNode(data []byte) (*HardDriveMediaNode, error) {
	if len(data) != 38 {
		return nil, fmt.Errorf("invalid hard drive node payload size: got %d, want 38", len(data))
	}
	sig, err := ParseEFIGUID(data[20:36])
	if err != nil {
		return nil, err
	}

	return &HardDriveMediaNode{
		PartitionNumber:      binary.LittleEndian.Uint32(data[0:4]),
		PartitionSectorStart: binary.LittleEndian.Uint64(data[4:12]),
		PartitionSectorSize:  binary.LittleEndian.Uint64(data[12:20]),
		Signature:            sig,
		PartitionType:        PartitionType(data[36]),
		SignatureType:        SignatureType(data[37]),
	}, nil
}
