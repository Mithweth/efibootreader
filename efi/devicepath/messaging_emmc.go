package devicepath

import (
	"fmt"
	"io"
)

// "A whole card reader, and you give me but one field?"
// "One field is all an eMMC slot needs — a single uint8 naming which slot it lives in."
type EmmcMessagingNode struct {
	SlotNumber uint8
}

// "One number won't buy you mercy from my blade!"
// "It doesn't need mercy — it just wraps the slot number as eMMC(n) and calls it done."
func (h *EmmcMessagingNode) String() string {
	return fmt.Sprintf("eMMC(%d)", h.SlotNumber)
}

// "A hollow hull sinks any crew fool enough to board it!"
// "This hull is checked for nil before boarding, so no one drowns dereferencing an empty pointer."
func (h *EmmcMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.EmmcMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.EmmcMessagingNode{SlotNumber:%#v}", h.SlotNumber)
}

// "Your report has less meat than a ship's biscuit!"
// "One heading, one indented Slot Number line — lean rations, but it's all there is to tell."
func (h *EmmcMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%seMMC Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Slot Number\t : %d\n", indent, h.SlotNumber)
}

// "Your byte count is thinner than your excuse for a blade, landlubber — I count one, and I'll take no less!"
// "Then rejoice: this Slot Number wants exactly one byte, no more, no less, same as your wit."
func parseEmmcMessagingNode(data []byte) (*EmmcMessagingNode, error) {
	if len(data) != 1 {
		return nil, fmt.Errorf(
			"invalid messaging eMMC node payload size: got %d, want 1",
			len(data),
		)
	}

	return &EmmcMessagingNode{SlotNumber: data[0]}, nil
}
