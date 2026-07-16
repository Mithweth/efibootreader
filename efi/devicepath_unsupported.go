package efi

import "fmt"

func (n UnknownDevicePathNode) String() string {
	return fmt.Sprintf(
		"Unknown(0x%02x,0x%02x,%x)",
		n.Type,
		n.SubType,
		n.Data,
	)
}

func unknownDevicePathNode(node DevicePathNode) fmt.Stringer {
	return UnknownDevicePathNode{
		Type:    node.Type,
		SubType: node.SubType,
		Data:    node.Data,
	}
}
