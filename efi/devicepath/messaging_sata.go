package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SataMessagingNode struct {
	HBAPortNumber            uint16
	PortMultiplierPortNumber uint16
	LogicalUnitNumber        uint16
}

func (h *SataMessagingNode) String() string {
	return fmt.Sprintf(
		"Sata(%d,%d,%d)",
		h.HBAPortNumber,
		h.PortMultiplierPortNumber,
		h.LogicalUnitNumber,
	)
}

func (h *SataMessagingNode) GoString() string {
	if h == nil {
		return "(*efi.SataMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&efi.SataMessagingNode{"+
			"HBAPortNumber:%#v, "+
			"PortMultiplierPortNumber:%#v, "+
			"LogicalUnitNumber:%#v}",
		h.HBAPortNumber,
		h.PortMultiplierPortNumber,
		h.LogicalUnitNumber,
	)
}

func (h *SataMessagingNode) dump(w io.Writer, indent string) {
	fmt.Fprintf(w, "%sSATA Messaging Node\n", indent)
	fmt.Fprintf(w, "%s  HBA Port Number\t\t : %d\n", indent, h.HBAPortNumber)
	fmt.Fprintf(w, "%s  Port Multiplier Port Number\t : %d\n", indent, h.PortMultiplierPortNumber)
	fmt.Fprintf(w, "%s  Logical Unit Number\t\t : %d\n", indent, h.LogicalUnitNumber)
}

func parseSataMessagingNode(data []byte) (*SataMessagingNode, error) {
	if len(data) != 6 {
		return nil, fmt.Errorf(
			"invalid messaging SATA node payload size: got %d, want 6",
			len(data),
		)
	}

	return &SataMessagingNode{
		HBAPortNumber:            binary.LittleEndian.Uint16(data[0:2]),
		PortMultiplierPortNumber: binary.LittleEndian.Uint16(data[2:4]),
		LogicalUnitNumber:        binary.LittleEndian.Uint16(data[4:6]),
	}, nil
}
