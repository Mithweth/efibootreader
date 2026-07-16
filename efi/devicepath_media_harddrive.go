package efi

import (
	"encoding/binary"
	"fmt"
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
	Signature            GUID
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

func (h *HardDriveMediaNode) GoString() string {
	if h == nil {
		return "(*efi.HardDriveMediaNode)(nil)"
	}

	return fmt.Sprintf(
		"&efi.HardDriveMediaNode{"+
			"PartitionNumber:%#v, "+
			"PartitionSectorStart:%#v, "+
			"PartitionSectorSize:%#v, "+
			"PartitionType:%#v, "+
			"Signature:%#v, "+
			"SignatureType:%#v}",
		h.PartitionNumber,
		h.PartitionSectorStart,
		h.PartitionSectorSize,
		h.PartitionType,
		h.Signature,
		h.SignatureType,
	)
}

func parseHardDriveMediaNode(data []byte) (*HardDriveMediaNode, error) {
	if len(data) != 38 {
		return nil, fmt.Errorf("invalid hard drive node payload size: got %d, want 38", len(data))
	}
	sig, err := ParseGUID(data[20:36])
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

func (v PartitionType) String() string {
	switch v {
	case PartitionMBR:
		return "MBR"
	case PartitionGPT:
		return "GPT"
	default:
		return fmt.Sprintf("UnknownPartitionType(%#x)", uint8(v))
	}
}

func (v PartitionType) GoString() string {
	switch v {
	case PartitionMBR:
		return "efi.PartitionMBR"
	case PartitionGPT:
		return "efi.PartitionGPT"
	default:
		return fmt.Sprintf("efi.Partition(%#x)", uint8(v))
	}
}