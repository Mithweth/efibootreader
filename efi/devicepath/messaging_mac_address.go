package devicepath

import (
	"fmt"
	"io"
)

type MacAddressMessagingNode struct {
	Address       [32]byte
	InterfaceType uint8
}

func (h *MacAddressMessagingNode) String() string {
	addressLength := len(h.Address)
	if h.InterfaceType == 0 || h.InterfaceType == 1 {
		addressLength = 6
	}

	return fmt.Sprintf("MAC(%x,%d)", h.Address[:addressLength], h.InterfaceType)
}

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

func (h *MacAddressMessagingNode) dump(w io.Writer, indent string) {
	addressLength := len(h.Address)
	if h.InterfaceType == 0 || h.InterfaceType == 1 {
		addressLength = 6
	}

	_, _ = fmt.Fprintf(w, "%sMAC Address Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  MAC Address\t\t : %x\n", indent, h.Address[:addressLength])
	_, _ = fmt.Fprintf(w, "%s  Interface Type\t : %d\n", indent, h.InterfaceType)
}

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
