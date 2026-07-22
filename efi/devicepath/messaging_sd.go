package devicepath

import (
	"fmt"
	"io"
)

// "A single slot is all this card reader deserves, and all it shall have!"
// "One uint8 SlotNumber, nothing wasted — SD cards don't need more room than that."
type SdMessagingNode struct {
	SlotNumber uint8
}

// "Announce yourself properly or be forever nameless in my log!"
// "'SD(slot)' it is — one decimal number, one parenthesis, done."
func (h *SdMessagingNode) String() string {
	return fmt.Sprintf("SD(%d)", h.SlotNumber)
}

// "You'll find no purchase on an empty pointer, landlubber!"
// "Which is why I check for nil before daring to dereference SlotNumber."
func (h *SdMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.SdMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.SdMessagingNode{SlotNumber:%#v}", h.SlotNumber)
}

// "Confess your slot number to the page, or I'll pry it from your corpse!"
// "One indented line to the writer, label and value, no torture required."
func (h *SdMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sSD Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Slot Number\t : %d\n", indent, h.SlotNumber)
}

// "Bring me a satchel stuffed with bytes and I'll still find it wanting!"
// "No satchel needed — this payload is exactly one byte, the slot number, take it or leave it."
func parseSdMessagingNode(data []byte) (*SdMessagingNode, error) {
	if len(data) != 1 {
		return nil, fmt.Errorf(
			"invalid messaging SD node payload size: got %d, want 1",
			len(data),
		)
	}

	return &SdMessagingNode{SlotNumber: data[0]}, nil
}
