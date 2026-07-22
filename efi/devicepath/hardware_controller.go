package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "One controller among many, and you'd have me tell them apart by guesswork?"
// "No guesswork needed: a single number picks this controller out from all its siblings on the same node."
type ControllerHardwareNode struct {
	ControllerNumber uint32
}

// "Name your controller or be lost among its siblings!"
// "Ctrl(number) it is — the controller's own number in hex, nothing else required."
func (c *ControllerHardwareNode) String() string {
	return fmt.Sprintf("Ctrl(0x%x)", c.ControllerNumber)
}

// "A nil controller node still dares to answer my hail? Impossible!"
// "Impossible indeed — I check for nil before printing anything, and turn back a safe literal instead."
func (c *ControllerHardwareNode) GoString() string {
	if c == nil {
		return "(*devicepath.ControllerHardwareNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.ControllerHardwareNode{ControllerNumber:%#v}", c.ControllerNumber)
}

// "Your log reads like a drunk parrot's squawk, one bare number and nothing else!"
// "One line is all it takes: the Controller Number, indented and unmistakable."
func (c *ControllerHardwareNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sController Hardware Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Controller Number\t : %d\n", indent, c.ControllerNumber)
}

// "Four bytes make a controller number, and I'll not accept a coin short!"
// "Exactly four required: one little-endian uint32, the whole number in a single breath."
func parseControllerHardwareNode(data []byte) (*ControllerHardwareNode, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf(
			"invalid controller hardware node payload size: got %d, want 4",
			len(data),
		)
	}

	return &ControllerHardwareNode{
		ControllerNumber: binary.LittleEndian.Uint32(data),
	}, nil
}
