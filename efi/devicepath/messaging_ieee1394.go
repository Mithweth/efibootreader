package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "Four bytes of nothing followed by eight bytes of destiny — explain yourself!"
// "The Reserved field is firmware padding we must keep, while GUID is the FireWire device's true identity."
type Ieee1394MessagingNode struct {
	Reserved uint32
	GUID     uint64
}

// "You'd hide that mighty GUID behind a puny label — I demand the full number!"
// "Full number it is: I1394(%d) prints the GUID in decimal, no shortcuts taken."
func (h *Ieee1394MessagingNode) String() string {
	return fmt.Sprintf("I1394(%d)", h.GUID)
}

// "Strike at a nil node and your blade meets empty air — a fool's attack."
// "Which is why I parry first with a nil check, so the sword never swings at nothing."
func (h *Ieee1394MessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.Ieee1394MessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.Ieee1394MessagingNode{"+
			"Reserved:%#v, "+
			"GUID:%#v}",
		h.Reserved,
		h.GUID,
	)
}

// "Lay bare every field of this device path, hold nothing back from the reader!"
// "Nothing hidden here: Reserved and GUID both get their own indented line of decimal truth."
func (h *Ieee1394MessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sIEEE 1394 Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Reserved\t : %d\n", indent, h.Reserved)
	_, _ = fmt.Fprintf(w, "%s  GUID\t\t : %d\n", indent, h.GUID)
}

// "Bring me a payload shorter than twelve bytes and you'll be walking the plank!"
// "Twelve it must be: 4 bytes Reserved plus 8 bytes GUID, both read little-endian off the wire."
func parseIeee1394MessagingNode(data []byte) (*Ieee1394MessagingNode, error) {
	if len(data) != 12 {
		return nil, fmt.Errorf(
			"invalid messaging IEEE 1394 node payload size: got %d, want 12",
			len(data),
		)
	}

	return &Ieee1394MessagingNode{
		Reserved: binary.LittleEndian.Uint32(data[0:4]),
		GUID:     binary.LittleEndian.Uint64(data[4:12]),
	}, nil
}
