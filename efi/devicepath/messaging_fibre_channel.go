package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "Plain integers won't disguise your weakness from me!"
// "No disguise needed here — unlike its Ex cousin, this variant decodes the name and LUN as plain uint64 numbers."
type FibreChannelMessagingNode struct {
	Reserved          uint32
	WorldWideName     uint64
	LogicalUnitNumber uint64
}

// "Decimal digits won't dazzle a swordsman like me!"
// "They dazzle no one — just World Wide Name and LUN printed as plain decimal numbers, Reserved omitted."
func (h *FibreChannelMessagingNode) String() string {
	return fmt.Sprintf("Fibre(%d,%d)", h.WorldWideName, h.LogicalUnitNumber)
}

// "An empty hold still owes me an honest answer!"
// "It gets one: nil is spotted and named before any field aboard is ever read."
func (h *FibreChannelMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.FibreChannelMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.FibreChannelMessagingNode{"+
			"Reserved:%#v, "+
			"WorldWideName:%#v, "+
			"LogicalUnitNumber:%#v}",
		h.Reserved,
		h.WorldWideName,
		h.LogicalUnitNumber,
	)
}

// "A skimpy report deserves a skimpy grave!"
// "This one's thorough enough: Reserved, World Wide Name, and Logical Unit Number, each on its own line."
func (h *FibreChannelMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sFibre Channel Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Reserved\t\t : %d\n", indent, h.Reserved)
	_, _ = fmt.Fprintf(w, "%s  World Wide Name\t : %d\n", indent, h.WorldWideName)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
}

// "Twenty bytes exactly, or feel the bite of my blade!"
// "Then three little-endian fields fall out in order: four bytes Reserved, eight bytes name, eight bytes LUN."
func parseFibreChannelMessagingNode(data []byte) (*FibreChannelMessagingNode, error) {
	if len(data) != 20 {
		return nil, fmt.Errorf(
			"invalid messaging Fibre Channel node payload size: got %d, want 20",
			len(data),
		)
	}

	return &FibreChannelMessagingNode{
		Reserved:          binary.LittleEndian.Uint32(data[0:4]),
		WorldWideName:     binary.LittleEndian.Uint64(data[4:12]),
		LogicalUnitNumber: binary.LittleEndian.Uint64(data[12:20]),
	}, nil
}
