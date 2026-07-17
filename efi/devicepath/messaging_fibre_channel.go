package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

type FibreChannelMessagingNode struct {
	Reserved          uint32
	WorldWideName     uint64
	LogicalUnitNumber uint64
}

func (h *FibreChannelMessagingNode) String() string {
	return fmt.Sprintf("Fibre(%d,%d)", h.WorldWideName, h.LogicalUnitNumber)
}

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

func (h *FibreChannelMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sFibre Channel Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Reserved\t\t : %d\n", indent, h.Reserved)
	_, _ = fmt.Fprintf(w, "%s  World Wide Name\t : %d\n", indent, h.WorldWideName)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
}

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
