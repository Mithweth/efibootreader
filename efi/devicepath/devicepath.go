package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

func (p *DevicePath) String() string {
	var instances []string

	for _, instance := range p.Instances {
		instances = append(instances, instance.String())
	}

	return strings.Join(instances, ",")
}

func (d *DevicePath) Dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sDevicePath\n", indent)

	for _, inst := range d.Instances {
		inst.dump(w, indent+"  ")
	}
}

func (i DevicePathInstance) String() string {
	var nodes []string
	for _, node := range i.Nodes {
		nodes = append(nodes, node.Details.String())
	}
	return strings.Join(nodes, "/")
}

func (i DevicePathInstance) GoString() string {
	var nodes []string
	for _, node := range i.Nodes {
		nodes = append(nodes, fmt.Sprintf("%#v", node.Details))
	}
	return fmt.Sprintf("[]efi.DevicePathNode{%s}", strings.Join(nodes, ", "))
}

func (i DevicePathInstance) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sDevicePathInstance\n", indent)

	for _, node := range i.Nodes {
		node.Details.dump(w, indent+"  ")
	}
}

func (n DevicePathNode) String() string {
	return fmt.Sprintf("Unknown(%d,%d,%x)", n.Type, n.SubType, n.Data)
}

func isEndEntireDevicePath(node DevicePathNode) bool {
	return node.Type == DevicePathEnd && EndDevicePathSubType(node.SubType) == EndEntireDevicePathSubType
}

func isEndThisInstance(node DevicePathNode) bool {
	return node.Type == DevicePathEnd && EndDevicePathSubType(node.SubType) == EndThisInstanceSubType
}

func ParseDevicePath(data []byte) (*DevicePath, error) {
	var instances []DevicePathInstance
	var nodes []DevicePathNode

	for offset := 0; offset < len(data); {
		if offset+4 > len(data) {
			return nil, fmt.Errorf(
				"truncated device path header at offset %d",
				offset,
			)
		}

		nodeLength := int(binary.LittleEndian.Uint16(data[offset+2 : offset+4]))

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

		if isEndThisInstance(node) {
			instances = append(instances, DevicePathInstance{Nodes: nodes})
			nodes = []DevicePathNode{}
			continue
		}

		if isEndEntireDevicePath(node) {
			instances = append(instances, DevicePathInstance{Nodes: nodes})
			break
		}

		details, err := parseDevicePathNode(node)
		if err != nil {
			return nil, err
		}

		node.Details = details
		nodes = append(nodes, node)
	}

	return &DevicePath{Instances: instances}, nil
}

func parseDevicePathNode(node DevicePathNode) (DevicePathNodeDetails, error) {
	switch node.Type {
	case DevicePathMedia:
		return parseMediaDevicePathNode(node)

	// case DevicePathHardware:
	// 	return parseHardwareDevicePathNode(node)

	// case DevicePathACPI:
	// 	return parseACPIDevicePathNode(node)

	case DevicePathMessaging:
		return parseMessagingDevicePathNode(node)

	// case DevicePathBBS:
	// 	return parseBBSDevicePathNode(node)

	default:
		return unknownDevicePathNode(node), nil
	}
}
