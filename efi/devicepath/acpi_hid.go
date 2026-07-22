package devicepath

import (
	"encoding/binary"
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "A HID and a UID, and you'd have me believe that's a whole device's life story?"
// "Enough of a story: the compressed hardware ID names the device, the unique ID tells two of a kind apart."
type HidAcpiNode struct {
	HID identifiers.EISAID
	UID uint32
}

// "Name your device or be forever known as a mystery number!"
// "Acpi(HID,UID) it is — HID with its readable PNP name."
func (h *HidAcpiNode) String() string {
	return fmt.Sprintf("Acpi(%s,%d)", h.HID, h.UID)
}

// "A nil ACPI node still dares to answer my hail? Impossible!"
// "Impossible indeed — I check for nil before printing anything, and turn back a safe literal instead."
func (h *HidAcpiNode) GoString() string {
	if h == nil {
		return "(*devicepath.HidAcpiNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.HidAcpiNode{HID:%#v, UID:%d}", h.HID, h.UID)
}

// "Your log reads like a drunk parrot's squawk, numbers and nothing else!"
// "Mine prints the PNP name, then the UID plain."
func (h *HidAcpiNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sACPI Device Path\n", indent)
	_, _ = fmt.Fprintf(w, "%s  HID\t : %s", indent, h.HID)
	if description, ok := identifiers.LookupEISAID(h.HID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n%s  UID\t : %d\n", indent, h.UID)
}

// "Eight bytes make an ACPI address, and I'll not accept a coin short!"
// "Exactly eight required: HID from the first four bytes, UID from the last four, both little-endian."
func parseHidAcpiNode(data []byte) (*HidAcpiNode, error) {
	if len(data) != 8 {
		return nil, fmt.Errorf(
			"invalid ACPI HID node payload size: got %d, want 8",
			len(data),
		)
	}
	hid, err := identifiers.ParseEISAID(data[0:4])
	if err != nil {
		return nil, fmt.Errorf("parse vendor EISAID: %w", err)
	}

	return &HidAcpiNode{
		HID: hid,
		UID: binary.LittleEndian.Uint32(data[4:8]),
	}, nil
}
