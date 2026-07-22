package devicepath

import (
	"fmt"
	"io"
)

// "A PC Card slot named by anything less than its own socket number? Unthinkable!"
// "Unthinkable indeed — one byte, the socket's Function Number, is all this card needs to be found."
type PccardHardwareNode struct {
	FunctionNumber uint8
}

// "Name your socket, PC Card, or be forgotten in the chassis!"
// "PcCard(function) it is — the socket number in hex, plain as the slot itself."
func (p *PccardHardwareNode) String() string {
	return fmt.Sprintf("PcCard(%d)", p.FunctionNumber)
}

// "A nil PC Card node still dares to answer my hail? Impossible!"
// "Impossible indeed — I check for nil before printing anything, and turn back a safe literal instead."
func (p *PccardHardwareNode) GoString() string {
	if p == nil {
		return "(*devicepath.PccardHardwareNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.PccardHardwareNode{FunctionNumber:%#v}", p.FunctionNumber)
}

// "Your log reads like a drunk parrot's squawk, one bare number and nothing else!"
// "One line is all it takes: the Function Number, indented and unmistakable."
func (p *PccardHardwareNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sPCCARD Hardware Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Function Number\t : %d\n", indent, p.FunctionNumber)
}

// "One byte makes a PC Card socket, and I'll not accept an empty slot!"
// "Exactly one required: the Function Number, and nothing more, fills the whole payload."
func parsePccardHardwareNode(data []byte) (*PccardHardwareNode, error) {
	if len(data) != 1 {
		return nil, fmt.Errorf(
			"invalid PCCARD hardware node payload size: got %d, want 1",
			len(data),
		)
	}

	return &PccardHardwareNode{FunctionNumber: data[0]}, nil
}
