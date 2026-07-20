package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

type FibreChannelExMessagingNode struct {
	Reserved          uint32
	WorldWideName     [8]byte
	LogicalUnitNumber [8]byte
}

func (h *FibreChannelExMessagingNode) String() string {
	return fmt.Sprintf("FibreEx(0x%x,0x%x)", h.WorldWideName, h.LogicalUnitNumber)
}

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

func (h *FibreChannelExMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sFibre Channel Ex Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Reserved\t\t : %#x\n", indent, h.Reserved)
	_, _ = fmt.Fprintf(w, "%s  World Wide Name\t : 0x%x\n", indent, h.WorldWideName)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : 0x%x\n", indent, h.LogicalUnitNumber)
}

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
