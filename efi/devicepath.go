package efi

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func (n DevicePathNode) String() string {
	return fmt.Sprintf(
		"Unknown(%d,%d,%x)",
		n.Type,
		n.SubType,
		n.Data,
	)
}

func (p *DevicePath) String() string {
	var nodes []string
	for _, node := range p.Nodes {
		nodes = append(nodes, node.Details.String())
	}
	return strings.Join(nodes, "/")
}

func isEndEntireDevicePath(node DevicePathNode) bool {
	return node.Type == DevicePathEnd && EndDevicePathSubType(node.SubType) == EndEntireDevicePathSubType
}

func ParseDevicePath(data []byte) (*DevicePath, error) {
	nodes := make([]DevicePathNode, 0)

	for offset := 0; offset < len(data); {
		if offset+4 > len(data) {
			return nil, fmt.Errorf(
				"truncated device path header at offset %d",
				offset,
			)
		}

		nodeLength := int(binary.LittleEndian.Uint16(
			data[offset+2 : offset+4],
		))

		if nodeLength < 4 {
			return nil, fmt.Errorf(
				"invalid device path node length %d at offset %d",
				nodeLength,
				offset,
			)
		}

		nodeEnd := offset + nodeLength
		if nodeEnd > len(data) {
			return nil, fmt.Errorf(
				"device path node exceeds buffer: offset=%d length=%d total=%d",
				offset,
				nodeLength,
				len(data),
			)
		}

		node := DevicePathNode{
			Type:    DevicePathType(data[offset]),
			SubType: data[offset+1],
			Data:    data[offset+4 : nodeEnd],
		}

		offset = nodeEnd

		if isEndEntireDevicePath(node) {
			break
		}

		details, err := parseDevicePathNode(node)
		if err != nil {
			return nil, err
		}

		node.Details = details
		nodes = append(nodes, node)
	}

	return &DevicePath{Nodes: nodes}, nil
}

func parseDevicePathNode(node DevicePathNode) (fmt.Stringer, error) {
	switch node.Type {
	case DevicePathMedia:
		return parseMediaDevicePathNode(node)

	// case DevicePathHardware:
	// 	return parseHardwareDevicePathNode(node)

	// case DevicePathACPI:
	// 	return parseACPIDevicePathNode(node)

	// case DevicePathMessaging:
	// 	return parseMessagingDevicePathNode(node)

	// case DevicePathBBS:
	// 	return parseBBSDevicePathNode(node)

	default:
		return unknownDevicePathNode(node), nil
	}
}
