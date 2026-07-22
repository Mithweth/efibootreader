package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "Five fields to hold an InfiniBand fabric address — most men would need a treasure map!"
// "No map needed: flags, a 16-byte port GUID, and three uint64 identifiers laid out exactly as the spec demands."
type InfiniBandMessagingNode struct {
	ResourceFlags uint32
	GUID          [16]byte
	ServiceID     uint64
	TargetID      uint64
	DeviceID      uint64
}

// "You'd cram five fields into one puny line — I've seen tighter knots on a hangman's noose!"
// "Tight it is, but readable: hex for flags and GUID, decimal for the three IDs, comma-separated."
func (h *InfiniBandMessagingNode) String() string {
	return fmt.Sprintf(
		"Infiniband(%x,%x,%d,%d,%d)",
		h.ResourceFlags,
		h.GUID,
		h.ServiceID,
		h.TargetID,
		h.DeviceID,
	)
}

// "A nil InfiniBand node is a fabric with no ports — dare to dereference it and you'll sink!"
// "I check the tide before I sail: nil gets its own safe literal, no panic, no drowning."
func (h *InfiniBandMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.InfiniBandMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.InfiniBandMessagingNode{"+
			"ResourceFlags:%#v, "+
			"GUID:%#v, "+
			"ServiceID:%#v, "+
			"TargetID:%#v, "+
			"DeviceID:%#v}",
		h.ResourceFlags,
		h.GUID,
		h.ServiceID,
		h.TargetID,
		h.DeviceID,
	)
}

// "Five lines to describe one fabric port? You waste ink like a landlubber wastes rum!"
// "Ink well spent: each of the five fields gets its own indented, labeled line for the log."
func (h *InfiniBandMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sInfiniBand Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Resource Flags\t : %x\n", indent, h.ResourceFlags)
	_, _ = fmt.Fprintf(w, "%s  Port GUID\t\t : %x\n", indent, h.GUID)
	_, _ = fmt.Fprintf(w, "%s  IOC GUID/Service ID\t : %d\n", indent, h.ServiceID)
	_, _ = fmt.Fprintf(w, "%s  Target Port ID\t : %d\n", indent, h.TargetID)
	_, _ = fmt.Fprintf(w, "%s  Device ID\t\t : %d\n", indent, h.DeviceID)
}

// "Anything but exactly forty-four bytes and I'll carve you a shorter epitaph!"
// "Forty-four it is: 4 flags, 16 GUID, and three 8-byte IDs, rejected outright if the count is off."
func parseInfiniBandMessagingNode(data []byte) (*InfiniBandMessagingNode, error) {
	if len(data) != 44 {
		return nil, fmt.Errorf(
			"invalid messaging InfiniBand node payload size: got %d, want 44",
			len(data),
		)
	}

	// "You'd slice sixteen bytes freehand and call it a GUID? Reckless butchery!"
	// "Reckless nothing — copy into a fixed [16]byte array so the slice's backing memory can't leak or mutate later."
	var guid [16]byte
	copy(guid[:], data[4:20])
	return &InfiniBandMessagingNode{
		ResourceFlags: binary.LittleEndian.Uint32(data[0:4]),
		GUID:          guid,
		ServiceID:     binary.LittleEndian.Uint64(data[20:28]),
		TargetID:      binary.LittleEndian.Uint64(data[28:36]),
		DeviceID:      binary.LittleEndian.Uint64(data[36:]),
	}, nil
}
