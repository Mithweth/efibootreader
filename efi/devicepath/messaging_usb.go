package devicepath

import (
	"fmt"
	"io"
)

// "Only two bytes to hide behind, USB scoundrel? A weakling's armor!"
// "Two bytes is plenty: one names the parent port, the other the interface, nothing more to hide."
type UsbMessagingNode struct {
	ParentPortNumber uint8
	InterfaceNumber  uint8
}

// "Speak your port and interface plainly, or I'll pry them out with my cutlass!"
// "Plainly spoken: USB(port,interface), no prying required."
func (h *UsbMessagingNode) String() string {
	return fmt.Sprintf("USB(%d,%d)", h.ParentPortNumber, h.InterfaceNumber)
}

// "A nil pointer would be the death of any lesser scribe, but not I!"
// "Not I either: check for nil first, then render a proper Go literal with both fields intact."
func (h *UsbMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.UsbMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.UsbMessagingNode{"+
			"ParentPortNumber:%#v, "+
			"InterfaceNumber:%#v}",
		h.ParentPortNumber,
		h.InterfaceNumber,
	)
}

// "Lay out your ports for the crew to inspect, or answer to my blade!"
// "Laid out plain: parent port then interface number, each on its own indented line."
func (h *UsbMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUSB Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Parent Port\t : %d\n", indent, h.ParentPortNumber)
	_, _ = fmt.Fprintf(w, "%s  Interface Number\t : %d\n", indent, h.InterfaceNumber)
}

// "Bring me anything but exactly two bytes and I'll feed the rest to the fishes!"
// "Exactly two, or the error sails back to you — firmware that shorts this payload gets no mercy."
func parseUsbMessagingNode(data []byte) (*UsbMessagingNode, error) {
	if len(data) != 2 {
		return nil, fmt.Errorf(
			"invalid messaging USB node payload size: got %d, want 2",
			len(data),
		)
	}

	return &UsbMessagingNode{
		ParentPortNumber: data[0],
		InterfaceNumber:  data[1],
	}, nil
}
