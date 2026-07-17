package devicepath

import (
	"fmt"
	"io"
)

type UfsMessagingNode struct {
	TargetID          uint8
	LogicalUnitNumber uint8
}

func (h *UfsMessagingNode) String() string {
	return fmt.Sprintf("UFS(%d,%d)", h.TargetID, h.LogicalUnitNumber)
}

func (h *UfsMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.UfsMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.ScsiMessagingNode{"+
			"TargetID:%#v, "+
			"LogicalUnitNumber:%#v}",
		h.TargetID,
		h.LogicalUnitNumber,
	)
}

func (h *UfsMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUFS Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Target ID\t\t : %d\n", indent, h.TargetID)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
}

func parseUfsMessagingNode(data []byte) (*UfsMessagingNode, error) {
	if len(data) != 2 {
		return nil, fmt.Errorf(
			"invalid messaging UFS node payload size: got %d, want 2",
			len(data),
		)
	}

	return &UfsMessagingNode{TargetID: data[0], LogicalUnitNumber: data[0]}, nil
}
