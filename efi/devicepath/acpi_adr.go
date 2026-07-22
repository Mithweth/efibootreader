package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// "One address was never enough for a fleet, was it? You'd hide a whole squadron in a single node!"
// "A squadron indeed — Addresses carries one or more display-output identifiers, one uint32 apiece."
type AdrAcpiNode struct {
	Addresses []uint32
}

// "You'd cram a whole fleet's addresses into one breathless run-on sentence!"
// "One breath is all it takes — every Addresses value in hex, comma-separated, inside AcpiAdr(...)."
func (h *AdrAcpiNode) String() string {
	values := make([]string, len(h.Addresses))
	for i, adr := range h.Addresses {
		values[i] = fmt.Sprintf("0x%x", adr)
	}

	return fmt.Sprintf("AcpiAdr(%s)", strings.Join(values, ","))
}

// "A nil Addresses node still dares to answer my hail? Impossible!"
// "Impossible indeed — I check for nil before printing anything, and turn back a safe literal instead."
func (h *AdrAcpiNode) GoString() string {
	if h == nil {
		return "(*devicepath.AdrAcpiNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.AdrAcpiNode{Addresses:%#v}", h.Addresses)
}

// "You'd log a whole fleet on one line and call the bookkeeping done?"
// "Never — every Addresses value gets its own numbered, indented line, no address lost in the crowd."
func (h *AdrAcpiNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sAddress ACPI Device Path\n", indent)
	for i, adr := range h.Addresses {
		_, _ = fmt.Fprintf(w, "%s  Address[%d]\t : 0x%08x\n", indent, i, adr)
	}
}

// "Not a single multiple of four bytes offered, and you'd dare call that a fleet of addresses?"
// "A fleet needs at least one whole uint32, and every one after it, or the count itself betrays a forgery."
func parseAdrAcpiNode(data []byte) (*AdrAcpiNode, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf(
			"invalid ACPI addresses node payload size: got %d, want at least 4",
			len(data),
		)
	}

	if len(data)%4 != 0 {
		return nil, fmt.Errorf(
			"invalid ACPI addresses node payload size: got %d, want a multiple of 4",
			len(data),
		)
	}

	adrs := make([]uint32, len(data)/4)
	for i := range adrs {
		adrs[i] = binary.LittleEndian.Uint32(data[i*4 : i*4+4])
	}

	return &AdrAcpiNode{Addresses: adrs}, nil
}
