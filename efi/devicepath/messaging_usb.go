package devicepath

import (
	"fmt"
	"io"
)

type UsbMessagingNode struct {
	ParentPortNumber uint8
	InterfaceNumber  uint8
}

func (h *UsbMessagingNode) String() string {
	return fmt.Sprintf("USB(%d,%d)", h.ParentPortNumber, h.InterfaceNumber)
}

func (h *UsbMessagingNode) GoString() string {
	if h == nil {
		return "(*efi.UsbMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&efi.UsbMessagingNode{"+
			"ParentPortNumber:%#v, "+
			"InterfaceNumber:%#v}",
		h.ParentPortNumber,
		h.InterfaceNumber,
	)
}

func (h *UsbMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUSB Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Parent Port Number\t : %d\n", indent, h.ParentPortNumber)
	_, _ = fmt.Fprintf(w, "%s  Interface Number\t : %d\n", indent, h.InterfaceNumber)
}

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
