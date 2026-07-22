package devicepath

import (
	"fmt"
	"io"
)

// "Your address is padded with more filler than a landlubber's excuses!"
// "Aye, 32 bytes reserved because the spec demands room for more than plain Ethernet."
type MacAddressMessagingNode struct {
	Address       [32]byte
	InterfaceType uint8
}

// "Six bytes or thirty-two — you never know which blade you're facing!"
// "Ethernet and token ring only ever bare six, so I trim the rest before anyone looks."
func (h *MacAddressMessagingNode) String() string {
	addressLength := len(h.Address)
	if h.InterfaceType == 0 || h.InterfaceType == 1 {
		addressLength = 6
	}

	return fmt.Sprintf("MAC(%x,%d)", h.Address[:addressLength], h.InterfaceType)
}

// "You'd let a nil receiver run you through without a fight!"
// "Never — I check the guard before the thrust and format only what's real."
func (h *MacAddressMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.MacAddressMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.MacAddressMessagingNode{"+
			"Address:%#v, "+
			"InterfaceType:%#v}",
		h.Address,
		h.InterfaceType,
	)
}

// "You'd print thirty-two bytes of padding and call it seamanship!"
// "Same trim as before, sir — only the meaningful bytes reach the page."
func (h *MacAddressMessagingNode) dump(w io.Writer, indent string) {
	addressLength := len(h.Address)
	if h.InterfaceType == 0 || h.InterfaceType == 1 {
		addressLength = 6
	}

	_, _ = fmt.Fprintf(w, "%sMAC Address Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  MAC Address\t\t : %x\n", indent, h.Address[:addressLength])
	_, _ = fmt.Fprintf(w, "%s  Interface Type\t : %d\n", indent, h.InterfaceType)
}

// "Thirty-three bytes or walk the plank — I'll count them myself if I must!"
// "32 for the padded address, 1 for the interface type, not a byte less."
func parseMacAddressMessagingNode(data []byte) (*MacAddressMessagingNode, error) {
	if len(data) != 33 {
		return nil, fmt.Errorf(
			"invalid messaging MAC address node payload size: got %d, want 33",
			len(data),
		)
	}

	node := &MacAddressMessagingNode{
		InterfaceType: data[32],
	}
	copy(node.Address[:], data[0:32])

	return node, nil
}
