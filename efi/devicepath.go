package efi

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"
	"github.com/google/uuid"
)


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

func ParseDevicePath(data []byte) ([]DevicePathNode, error) {
	var nodes []DevicePathNode

	offset := 0
	for {
		if offset+4 > len(data) {
			return nil, fmt.Errorf("truncated device path header at offset %d", offset)
		}

		length := int(binary.LittleEndian.Uint16(data[offset+2 : offset+4]))

		if length < 4 {
			return nil, fmt.Errorf("invalid device path node length %d at offset %d", length, offset)
		}

		if offset+length > len(data) {
			return nil, fmt.Errorf("device path node exceeds buffer: offset=%d length=%d total=%d", offset, length, len(data))
		}

		nodes = append(nodes, DevicePathNode{
			Type:    data[offset],
			SubType: data[offset+1],
			Data:    data[offset+4 : offset+length],
		})

		offset += length
		if offset >= len(data) {
			break
		}
	}

	return nodes, nil
}

func ParseHardDriveNode(data []byte) (*HardDriveNode, error) {
	if len(data) != 38 {
		return nil, fmt.Errorf("invalid hard drive node payload size: got %d, want 38", len(data))
	}
	sig, err := ParseEFIGUID(data[20:36])
	if err != nil {
		return nil, err
	}

	return &HardDriveNode{
		PartitionNumber: binary.LittleEndian.Uint32(data[0:4]),
		PartitionStart:  binary.LittleEndian.Uint64(data[4:12]),
		PartitionSize:   binary.LittleEndian.Uint64(data[12:20]),
		Signature:       sig,
		MBRType:         data[36],
		SignatureType:   data[37],
	}, nil
}

func ParseFilePathNode(data []byte) (string, error) {
	if len(data)%2 != 0 {
		return "", fmt.Errorf("invalid UTF-16 file path size: %d", len(data))
	}

	var codeUnits []uint16

	for i := 0; i < len(data); i += 2 {
		value := binary.LittleEndian.Uint16(data[i : i+2])
		if value == 0 {
			break
		}
		codeUnits = append(codeUnits, value)
	}

	return string(utf16.Decode(codeUnits)), nil
}
