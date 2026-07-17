package devicepath

import (
	"fmt"
	"io"
)

type EmmcMessagingNode struct {
	SlotNumber uint8
}

func (h *EmmcMessagingNode) String() string {
	return fmt.Sprintf("eMMC(%d)", h.SlotNumber)
}

func (h *EmmcMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.EmmcMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.EmmcMessagingNode{SlotNumber:%#v}", h.SlotNumber)
}

func (h *EmmcMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%seMMC Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Slot Number\t : %d\n", indent, h.SlotNumber)
}

func parseEmmcMessagingNode(data []byte) (*EmmcMessagingNode, error) {
	if len(data) != 1 {
		return nil, fmt.Errorf(
			"invalid messaging eMMC node payload size: got %d, want 1",
			len(data),
		)
	}

	return &EmmcMessagingNode{SlotNumber: data[0]}, nil
}
