package devicepath

import (
	"fmt"
	"io"
)

type SdMessagingNode struct {
	SlotNumber uint8
}

func (h *SdMessagingNode) String() string {
	return fmt.Sprintf("SD(%d)", h.SlotNumber)
}

func (h *SdMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.SdMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.SdMessagingNode{SlotNumber:%#v}", h.SlotNumber)
}

func (h *SdMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sSD Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Slot Number\t : %d\n", indent, h.SlotNumber)
}

func parseSdMessagingNode(data []byte) (*SdMessagingNode, error) {
	if len(data) != 1 {
		return nil, fmt.Errorf(
			"invalid messaging SD node payload size: got %d, want 1",
			len(data),
		)
	}

	return &SdMessagingNode{SlotNumber: data[0]}, nil
}
