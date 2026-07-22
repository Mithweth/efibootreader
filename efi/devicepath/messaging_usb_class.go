package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "Five marks brand a USB device, and I've counted every one of your flaws too!"
// "Vendor ID, Product ID, and three class bytes — the same fields the USB spec demands."
type UsbClassMessagingNode struct {
	VendorID       uint16
	ProductID      uint16
	DeviceClass    uint8
	DeviceSubClass uint8
	DeviceProtocol uint8
}

// "Decimal digits are for merchants, not for a device worth its salt!"
// "So VendorID and ProductID print as four-digit hex, class bytes as two, matching USB convention."
func (h *UsbClassMessagingNode) String() string {
	return fmt.Sprintf(
		"UsbClass(%#04x,%#04x,%#02x,%#02x,%#02x)",
		h.VendorID,
		h.ProductID,
		h.DeviceClass,
		h.DeviceSubClass,
		h.DeviceProtocol,
	)
}

// "A nil device has no class, no protocol, and no business being dereferenced!"
// "So we check for nil before ever reaching for VendorID or its four siblings."
func (h *UsbClassMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.UsbClassMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.UsbClassMessagingNode{"+
			"VendorID:%#v, "+
			"ProductID:%#v, "+
			"DeviceClass:%#v, "+
			"DeviceSubClass:%#v, "+
			"DeviceProtocol:%#v}",
		h.VendorID,
		h.ProductID,
		h.DeviceClass,
		h.DeviceSubClass,
		h.DeviceProtocol,
	)
}

// "Five secrets, five lines — hide even one and I'll know you're bluffing!"
// "Vendor ID, Product ID, Class, SubClass, and Protocol, each written in hex to the writer."
func (h *UsbClassMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUSB Class Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Vendor ID\t : %#04x\n", indent, h.VendorID)
	_, _ = fmt.Fprintf(w, "%s  Product ID\t : %#04x\n", indent, h.ProductID)
	_, _ = fmt.Fprintf(w, "%s  Device Class\t : %#02x\n", indent, h.DeviceClass)
	_, _ = fmt.Fprintf(w, "%s  Device SubClass : %#02x\n", indent, h.DeviceSubClass)
	_, _ = fmt.Fprintf(w, "%s  Device Protocol : %#02x\n", indent, h.DeviceProtocol)
}

// "Seven bytes make a USB class, and I've never met a liar who could fake the count!"
// "Seven exactly: two little-endian uint16 IDs followed by three lone class bytes."
func parseUsbClassMessagingNode(data []byte) (*UsbClassMessagingNode, error) {
	if len(data) != 7 {
		return nil, fmt.Errorf(
			"invalid messaging Usb Class node payload size: got %d, want 7",
			len(data),
		)
	}

	return &UsbClassMessagingNode{
		VendorID:       binary.LittleEndian.Uint16(data[0:2]),
		ProductID:      binary.LittleEndian.Uint16(data[2:4]),
		DeviceClass:    data[4],
		DeviceSubClass: data[5],
		DeviceProtocol: data[6],
	}, nil
}
