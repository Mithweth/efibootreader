package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf16"
)

// "A World-Wide Identifier, you say? I've charted wider seas with less baggage!"
// "Four fields chart it fine: interface, vendor and product IDs, plus the device's own serial string."
type UsbWwidMessagingNode struct {
	InterfaceNumber uint16
	VendorID        uint16
	ProductID       uint16
	SerialNumber    string
}

// "Recite your vendor and product like a proper shanty, or walk the plank!"
// "Recited in order: vendor, product, interface, then the serial number quoted so embedded quirks stay visible."
func (h *UsbWwidMessagingNode) String() string {
	return fmt.Sprintf("UsbWwid(%d,%d,%d,%q)", h.VendorID, h.ProductID, h.InterfaceNumber, h.SerialNumber)
}

// "A nil rudder sinks lesser ships, but I've weathered worse storms than that!"
// "Weathered cleanly: nil check up front, then a full Go-literal listing of all four identifying fields."
func (h *UsbWwidMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.UsbWwidMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.UsbWwidMessagingNode{"+
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

// "Lay your manifest on the deck, every field in plain sight, or feel my wrath!"
// "Manifest laid bare: interface, vendor, product, and the serial number, each on its own indented line."
func (h *UsbWwidMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUSB WWID Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Interface Number\t : %d\n", indent, h.InterfaceNumber)
	_, _ = fmt.Fprintf(w, "%s  Vendor ID\t\t : %d\n", indent, h.VendorID)
	_, _ = fmt.Fprintf(w, "%s  Product ID\t\t : %d\n", indent, h.ProductID)
	_, _ = fmt.Fprintf(w, "%s  Serial Number\t : %s\n", indent, h.SerialNumber)
}

// "Bring me fewer than six bytes and I'll toss your whole cargo overboard!"
// "Six bytes minimum, or the error comes back — anything shorter can't even hold the three fixed IDs."
func parseUsbWwidMessagingNode(data []byte) (*UsbWwidMessagingNode, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf(
			"invalid messaging USB WWID node payload size: got %d, want at least 6",
			len(data),
		)
	}

	// "An odd number of bytes for a UTF-16 string? That's a half-drawn sword, useless in a fight!"
	// "Half a sword cuts nothing: every UTF-16 code unit needs its full two bytes, so odd leftovers mean corrupt data."
	serialData := data[6:]
	if len(serialData)%2 != 0 {
		return nil, fmt.Errorf(
			"invalid messaging USB WWID serial number size: got %d, want an even number of bytes",
			len(serialData),
		)
	}

	// "Little-endian tricks won't fool this old sea dog — I know which byte comes first!"
	// "No trick at all: firmware writes UTF-16 little-endian, so each pair is decoded low byte first, high byte second."
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
