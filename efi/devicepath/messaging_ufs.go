package devicepath

import (
	"fmt"
	"io"
)

// "A UFS drive needs but two small marks to be found, unlike your oversized ego!"
// "Target ID and Logical Unit Number, both single bytes — UFS addressing is a tight fit."
type UfsMessagingNode struct {
	TargetID          uint8
	LogicalUnitNumber uint8
}

// "Announce your target or forever be known only as 'that ship'!"
// "'UFS(target,lun)' names both numbers plainly, in the order they were parsed."
func (h *UfsMessagingNode) String() string {
	return fmt.Sprintf("UFS(%d,%d)", h.TargetID, h.LogicalUnitNumber)
}

// "A nil blade cuts no one, and a nil pointer prints no lies about itself!"
// "Guarded up front, so a nil *UfsMessagingNode never gets its fields poked at."
func (h *UfsMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.UfsMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.UfsMessagingNode{"+
			"TargetID:%#v, "+
			"LogicalUnitNumber:%#v}",
		h.TargetID,
		h.LogicalUnitNumber,
	)
}

// "Write it down or it never happened, same as any duel without witnesses!"
// "Two indented lines to the writer: Target ID, then Logical Unit Number."
func (h *UfsMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUFS Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Target ID\t\t : %d\n", indent, h.TargetID)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
}

// "Two bytes make a UFS address, and I'll not accept a coin short!"
// "Exactly two required: the Target ID from data[0], and the Logical Unit Number from data[1], each its own coin."
func parseUfsMessagingNode(data []byte) (*UfsMessagingNode, error) {
	if len(data) != 2 {
		return nil, fmt.Errorf(
			"invalid messaging UFS node payload size: got %d, want 2",
			len(data),
		)
	}

	return &UfsMessagingNode{TargetID: data[0], LogicalUnitNumber: data[1]}, nil
}
