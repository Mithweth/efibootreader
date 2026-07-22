package devicepath

import (
	"encoding/binary"
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "You'd call any random byte a partition type and expect me to salute it?"
// "Only MBR and GPT wear this badge legitimately, encoded as a single byte."
type PartitionType uint8

const (
	PartitionMBR PartitionType = 1
	PartitionGPT PartitionType = 2
)

// "None, MBR, or GPT — pick a lane before I pick one for you, permanently!"
// "This single byte enumerates exactly those three signature kinds, iota and all."
type SignatureType uint8

const (
	SignatureNone SignatureType = iota
	SignatureMBR
	SignatureGPT
)

// "A hard drive partition with no signature is an unlocked chest waiting to be looted!"
// "Six fields — number, sector range, signature GUID, and both type bytes — lock down every detail firmware needs."
type HardDriveMediaNode struct {
	PartitionNumber      uint32
	PartitionSectorStart uint64
	PartitionSectorSize  uint64
	Signature            identifiers.GUID
	PartitionType        PartitionType
	SignatureType        SignatureType
}

// "Name your signature or I'll name it for you, in hex, without a shred of mercy!"
// "None, MBR, and GPT get proper names; anything else falls back to a plain hex value."
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

// "You'd print a constant as a naked number and dare call it Go code?"
// "I return the fully qualified devicepath.SignatureXxx identifier, or a typed hex literal for strangers."
func (s SignatureType) GoString() string {
	switch s {
	case SignatureNone:
		return "devicepath.SignatureNone"
	case SignatureMBR:
		return "devicepath.SignatureMBR"
	case SignatureGPT:
		return "devicepath.SignatureGPT"
	default:
		return fmt.Sprintf("devicepath.SignatureType(%#x)", uint8(s))
	}
}

// "MBR or GPT — anything else is a partition type only a scoundrel would claim!"
// "I name the two we recognize and print raw hex for whatever imposter shows up."
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

// "Your Go syntax is as fake as a wooden cutlass carved by a landlubber!"
// "I hand back the real devicepath.PartitionMBR or PartitionGPT constant name, hex fallback included."
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

// "You'd summarize a whole partition table with nothing but a shrug!"
// "I format number, type, signature, and both sector bounds in one HD(...) line, hex where it counts."
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

// "Strike a nil hard drive node and you'll only stab at empty air!"
// "I check for nil before laying out every field as valid, reconstructable Go syntax."
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

// "Sectors mean nothing without knowing where they begin and where they meet their end!"
// "I print start and size, compute the end sector myself, and name the signature if the lookup knows it."
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

// "Thirty-eight bytes is the toll for crossing this partition, pay it or be turned away!"
// "I enforce exactly 38 bytes, decode two little-endian uint64 sector fields, and parse the trailing GUID and type bytes."
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
