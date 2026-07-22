package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "This node holds but a single number, yet it commands the whole I2O bus!"
// "Aye, TargetID is all there is — one uint32 to name the target, nothing more to plunder."
type I2OMessagingNode struct {
	TargetID uint32
}

// "Speak plainly of your I2O target, or I'll have no idea whom I'm fighting."
// "Plainly it is: %d, formatted decimal, wrapped in I2O(...) for all to read."
func (h *I2OMessagingNode) String() string {
	return fmt.Sprintf("I2O(%d)", h.TargetID)
}

// "A nil pointer is a ghost ship — try to board it and you'll find nothing but grief."
// "Which is why I check for nil first and return a safe, printable ghost's name instead."
func (h *I2OMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.I2OMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.I2OMessagingNode{"+
			"TargetID:%#v}",
		h.TargetID,
	)
}

// "You dump your secrets to any writer who'll listen — no discretion at all!"
// "Discretion is for landlubbers; this prints two indented lines describing the node for humans."
func (h *I2OMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sI2O Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Target ID\t : %d\n", indent, h.TargetID)
}

// "Bring me anything but exactly four bytes and I'll feed your payload to the sharks."
// "Fair enough — the spec fixes this node at 4 bytes, so anything else is rejected before decoding."
func parseI2OMessagingNode(data []byte) (*I2OMessagingNode, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf(
			"invalid messaging I2O node payload size: got %d, want 4",
			len(data),
		)
	}

	return &I2OMessagingNode{
		TargetID: binary.LittleEndian.Uint32(data[0:4]),
	}, nil
}
