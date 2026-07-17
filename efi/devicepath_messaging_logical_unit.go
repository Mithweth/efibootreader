package efi

import (
	"fmt"
	"io"
)

type LogicalUnitMessagingNode struct {
	LogicalUnitNumber uint8
}

func (h *LogicalUnitMessagingNode) String() string {
	return fmt.Sprintf("Unit(%d)", h.LogicalUnitNumber)
}

func (h *LogicalUnitMessagingNode) GoString() string {
	if h == nil {
		return "(*efi.LogicalUnitMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&efi.LogicalUnitMessagingNode{"+
			"LogicalUnitNumber:%#v}",
		h.LogicalUnitNumber,
	)
}

func (h *LogicalUnitMessagingNode) dump(w io.Writer, indent string) {
	fmt.Fprintf(w, "%sLogical Unit Messaging Node\n", indent)
	fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
}

func parseLogicalUnitMessagingNode(data []byte) (*LogicalUnitMessagingNode, error) {
	if len(data) != 1 {
		return nil, fmt.Errorf(
			"invalid messaging logical unit node payload size: got %d, want 1",
			len(data),
		)
	}

	return &LogicalUnitMessagingNode{
		LogicalUnitNumber: data[0],
	}, nil
}
