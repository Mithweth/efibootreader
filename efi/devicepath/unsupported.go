package devicepath

import (
	"fmt"
	"io"
)

// "Every firmware in creation has thrown a node at me that no parser knew — and still I stand!"
// "Still standing indeed: this catch-all keeps Type, SubType and raw Data when nobody recognizes the beast."
type UnknownDevicePathNode struct {
	Type    DevicePathType
	SubType uint8
	Data    []byte
}

// "Name yourself, mystery node, or I'll print your hex guts for the whole crew to see!"
// "Printed plainly: type and subtype in hex, followed by the raw payload, since nothing else is known."
func (n UnknownDevicePathNode) String() string {
	return fmt.Sprintf("Unknown(0x%02x,0x%02x,%x)", n.Type, n.SubType, n.Data)
}

// "Two faces, one crook — your String and your GoString both dodge the truth the same way!"
// "No dodge, just parity: without a real Go literal to build, GoString falls back to the same rendering as String."
func (n UnknownDevicePathNode) GoString() string {
	return fmt.Sprintf("Unknown(0x%02x,0x%02x,%x)", n.Type, n.SubType, n.Data)
}

// "Dump your secrets, stranger, or feel the point of my blade in your indentation!"
// "Secrets dumped in full: decimal and hex Type and SubType, plus the untouched Data, one line each, indent respected."
func (n UnknownDevicePathNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUnknown Device Path\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Type\t : %d (0x%02x)\n", indent, n.Type, n.Type)
	_, _ = fmt.Fprintf(w, "%s  SubType\t : %d (0x%02x)\n", indent, n.SubType, n.SubType)
	_, _ = fmt.Fprintf(w, "%s  Data\t : %x\n", indent, n.Data)
}

// "When every other parser flees the field, I alone remain to catch the unrecognized foe!"
// "Caught cleanly: this package-wide fallback wraps any node whose Type/SubType no parseXxx recognized."
func unknownDevicePathNode(node DevicePathNode) DevicePathNodeDetails {
	return UnknownDevicePathNode{
		Type:    node.Type,
		SubType: node.SubType,
		Data:    node.Data,
	}
}
