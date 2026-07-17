package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf16"
)

type UsbWwidMessagingNode struct {
	InterfaceNumber uint16
	VendorID        uint16
	ProductID       uint16
	SerialNumber    string
}

func (h *UsbWwidMessagingNode) String() string {
	return fmt.Sprintf("UsbWwid(%d,%d,%d,%q)", h.VendorID, h.ProductID, h.InterfaceNumber, h.SerialNumber)
}

func (h *UsbWwidMessagingNode) GoString() string {
	if h == nil {
		return "(*efi.UsbWwidMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&efi.UsbWwidMessagingNode{"+
			"InterfaceNumber:%#v, "+
			"VendorID:%#v, "+
			"ProductID:%#v, "+
			"SerialNumber:%#v}",
		h.InterfaceNumber,
		h.VendorID,
		h.ProductID,
		h.SerialNumber,
	)
}

func (h *UsbWwidMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUSB WWID Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Interface Number\t : %d\n", indent, h.InterfaceNumber)
	_, _ = fmt.Fprintf(w, "%s  Vendor ID\t\t : %d\n", indent, h.VendorID)
	_, _ = fmt.Fprintf(w, "%s  Product ID\t\t : %d\n", indent, h.ProductID)
	_, _ = fmt.Fprintf(w, "%s  Serial Number\t : %s\n", indent, h.SerialNumber)
}

func parseUsbWwidMessagingNode(data []byte) (*UsbWwidMessagingNode, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf(
			"invalid messaging USB WWID node payload size: got %d, want at least 6",
			len(data),
		)
	}

	serialData := data[6:]
	if len(serialData)%2 != 0 {
		return nil, fmt.Errorf(
			"invalid messaging USB WWID serial number size: got %d, want an even number of bytes",
			len(serialData),
		)
	}

	serialUTF16 := make([]uint16, len(serialData)/2)
	for i := range serialUTF16 {
		serialUTF16[i] = binary.LittleEndian.Uint16(serialData[i*2 : i*2+2])
	}

	return &UsbWwidMessagingNode{
		InterfaceNumber: binary.LittleEndian.Uint16(data[0:2]),
		VendorID:        binary.LittleEndian.Uint16(data[2:4]),
		ProductID:       binary.LittleEndian.Uint16(data[4:6]),
		SerialNumber:    string(utf16.Decode(serialUTF16)),
	}, nil
}
