package devicepath

import (
	"encoding/binary"
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
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
	Signature            identifiers.GUID
	PartitionType        PartitionType
	SignatureType        SignatureType
}

func (s SignatureType) String() string {
	switch s {
	case SignatureNone:
		return "None"
	case SignatureMBR:
		return "MBR"
	case SignatureGPT:
		return "GPT"
	default:
		return fmt.Sprintf("%#x", uint8(s))
	}
}

func (s SignatureType) GoString() string {
	switch s {
	case SignatureNone:
		return "devicepath.SignatureNone"
	case SignatureMBR:
		return "devicepath.SignatureMBR"
	case SignatureGPT:
		return "devicepath.SignatureGPT"
	default:
		return fmt.Sprintf("devicepath.Signature(%#x)", uint8(s))
	}
}

func (v PartitionType) String() string {
	switch v {
	case PartitionMBR:
		return "MBR"
	case PartitionGPT:
		return "GPT"
	default:
		return fmt.Sprintf("%#x", uint8(v))
	}
}

func (v PartitionType) GoString() string {
	switch v {
	case PartitionMBR:
		return "devicepath.PartitionMBR"
	case PartitionGPT:
		return "devicepath.PartitionGPT"
	default:
		return fmt.Sprintf("devicepath.Partition(%#x)", uint8(v))
	}
}

func (h *HardDriveMediaNode) String() string {
	return fmt.Sprintf(
		"HD(%d,%s,%s,%x,%x)",
		h.PartitionNumber,
		h.PartitionType,
		h.Signature,
		h.PartitionSectorStart,
		h.PartitionSectorSize,
	)
}

func (h *HardDriveMediaNode) GoString() string {
	if h == nil {
		return "(*devicepath.HardDriveMediaNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.HardDriveMediaNode{"+
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

func (h *HardDriveMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sHard Drive Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Partition Number\t\t : %d\n", indent, h.PartitionNumber)
	_, _ = fmt.Fprintf(w, "%s  Partition Start (Sectors)\t : %d\n", indent, h.PartitionSectorStart)
	_, _ = fmt.Fprintf(w, "%s  Partition Size (Sectors)\t : %d\n", indent, h.PartitionSectorSize)
	_, _ = fmt.Fprintf(w, "%s  Partition End (Sectors)\t : %d\n", indent, h.PartitionSectorStart+h.PartitionSectorSize)
	_, _ = fmt.Fprintf(w, "%s  Partition Type\t\t : %s\n", indent, h.PartitionType)
	_, _ = fmt.Fprintf(w, "%s  Signature\t\t\t : %s", indent, h.Signature)
	if description, ok := identifiers.LookupGUID(h.Signature); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n%s  Signature Type\t\t : %s\n", indent, h.SignatureType)
}

func parseHardDriveMediaNode(data []byte) (*HardDriveMediaNode, error) {
	if len(data) != 38 {
		return nil, fmt.Errorf("invalid hard drive node payload size: got %d, want 38", len(data))
	}
	sig, err := identifiers.ParseGUID(data[20:36])
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
