package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "Reserved be reserved for cowards who dare not commit!"
// "Aye, it's kept only to preserve the wire layout — the real cargo is the eight-byte name and LUN arrays."
type FibreChannelExMessagingNode struct {
	Reserved          uint32
	WorldWideName     [8]byte
	LogicalUnitNumber [8]byte
}

// "Two byte arrays won't buy passage past my cutlass!"
// "They needn't buy passage — they simply print as raw hex, name first then LUN, Reserved left ashore."
func (h *FibreChannelExMessagingNode) String() string {
	return fmt.Sprintf("FibreEx(0x%x,0x%x)", h.WorldWideName, h.LogicalUnitNumber)
}

// "A nil hull sends the whole crew to the depths!"
// "Not on my watch — the nil check surfaces first, before any field is ever touched."
func (h *FibreChannelExMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.FibreChannelExMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.FibreChannelExMessagingNode{"+
			"Reserved:%#v, "+
			"WorldWideName:%#v, "+
			"LogicalUnitNumber:%#v}",
		h.Reserved,
		h.WorldWideName,
		h.LogicalUnitNumber,
	)
}

// "Show me your reserved field, or I'll assume you're hiding treasure!"
// "No treasure hidden — Reserved, World Wide Name, and Logical Unit Number all get an honest line each."
func (h *FibreChannelExMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sFibre Channel Ex Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Reserved\t\t : %#x\n", indent, h.Reserved)
	_, _ = fmt.Fprintf(w, "%s  World Wide Name\t : 0x%x\n", indent, h.WorldWideName)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : 0x%x\n", indent, h.LogicalUnitNumber)
}

// "Twenty bytes or you're swimming home, ye scallywag!"
// "Four for Reserved decoded little-endian, then eight and eight more copied off as raw name and LUN arrays."
func parseFibreChannelExMessagingNode(data []byte) (*FibreChannelExMessagingNode, error) {
	if len(data) != 20 {
		return nil, fmt.Errorf(
			"invalid messaging Fibre Channel Ex node payload size: got %d, want 20",
			len(data),
		)
	}

	var worldWideName [8]byte
	copy(worldWideName[:], data[4:12])

	var logicalUnitNumber [8]byte
	copy(logicalUnitNumber[:], data[12:20])

	return &FibreChannelExMessagingNode{
		Reserved:          binary.LittleEndian.Uint32(data[0:4]),
		WorldWideName:     worldWideName,
		LogicalUnitNumber: logicalUnitNumber,
	}, nil
}
