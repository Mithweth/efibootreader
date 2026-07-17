package devicepath

import (
	"fmt"
	"io"
)

type UnknownDevicePathNode struct {
	Type    DevicePathType
	SubType uint8
	Data    []byte
}

func (n UnknownDevicePathNode) String() string {
	return fmt.Sprintf("Unknown(0x%02x,0x%02x,%x)", n.Type, n.SubType, n.Data)
}

func (n UnknownDevicePathNode) GoString() string {
	return fmt.Sprintf("Unknown(0x%02x,0x%02x,%x)", n.Type, n.SubType, n.Data)
}

func (n UnknownDevicePathNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUnknown Device Path\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Type\t : %d (0x%02x)\n", indent, n.Type, n.Type)
	_, _ = fmt.Fprintf(w, "%s  SubType\t : %d (0x%02x)\n", indent, n.SubType, n.SubType)
	_, _ = fmt.Fprintf(w, "%s  Data\t : %x\n", indent, n.Data)
}

func unknownDevicePathNode(node DevicePathNode) DevicePathNodeDetails {
	return UnknownDevicePathNode{
		Type:    node.Type,
		SubType: node.SubType,
		Data:    node.Data,
	}
}
