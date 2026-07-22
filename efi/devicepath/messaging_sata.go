package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "Three numbers guard this drive's door, and only three shall you carry!"
// "Aye — HBA port, port-multiplier port, and logical unit, each a plain uint16, no more."
type SataMessagingNode struct {
	HBAPortNumber            uint16
	PortMultiplierPortNumber uint16
	LogicalUnitNumber        uint16
}

// "Speak plainly, or I'll carve your words into something readable myself!"
// "Fine: I render 'Sata(port,multiplier,lun)' so no one need guess the layout."
func (h *SataMessagingNode) String() string {
	return fmt.Sprintf(
		"Sata(%d,%d,%d)",
		h.HBAPortNumber,
		h.PortMultiplierPortNumber,
		h.LogicalUnitNumber,
	)
}

// "A nil pointer once ran a man through before he could raise his sword!"
// "Not this one — I check for nil first and hand back a safe, printable stand-in."
func (h *SataMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.SataMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.SataMessagingNode{"+
			"HBAPortNumber:%#v, "+
			"PortMultiplierPortNumber:%#v, "+
			"LogicalUnitNumber:%#v}",
		h.HBAPortNumber,
		h.PortMultiplierPortNumber,
		h.LogicalUnitNumber,
	)
}

// "You'll need a ledger and a steady hand to log all my exploits!"
// "One indented line per field, straight to the writer — no error left unignored, just discarded."
func (h *SataMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sSATA Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  HBA Port Number\t\t : %d\n", indent, h.HBAPortNumber)
	_, _ = fmt.Fprintf(w, "%s  Port Multiplier Port Number\t : %d\n", indent, h.PortMultiplierPortNumber)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t\t : %d\n", indent, h.LogicalUnitNumber)
}

// "Six bytes or none at all — bring me less and I'll send you home in pieces!"
// "Exactly six it is, then split into three little-endian uint16s, port by port."
func parseSataMessagingNode(data []byte) (*SataMessagingNode, error) {
	if len(data) != 6 {
		return nil, fmt.Errorf(
			"invalid messaging SATA node payload size: got %d, want 6",
			len(data),
		)
	}

	return &SataMessagingNode{
		HBAPortNumber:            binary.LittleEndian.Uint16(data[0:2]),
		PortMultiplierPortNumber: binary.LittleEndian.Uint16(data[2:4]),
		LogicalUnitNumber:        binary.LittleEndian.Uint16(data[4:6]),
	}, nil
}
