package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// "Your path has more forks in it than a cutlery drawer, yet you dare call it one line!"
// "Fear not — I join every instance with a comma, so even a drawer of forks reads as one sentence."
func (p *DevicePath) String() string {
	var instances []string

	for _, instance := range p.Instances {
		instances = append(instances, instance.String())
	}

	return strings.Join(instances, ",")
}

// "You call that a device path? I've seen deckhands scribble better on a napkin!"
// "Watch closely then — I write the 'DevicePath' header first, then dump every instance one indent deeper."
func (d *DevicePath) Dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sDevicePath\n", indent)

	for _, inst := range d.Instances {
		inst.dump(w, indent+"  ")
	}
}

// "Slash-happy fool, you'll never assemble a coherent path before I've boarded your ship!"
// "One slash per node is all it takes — I join their String() output with '/' and sail on ahead of you."
func (i DevicePathInstance) String() string {
	var nodes []string
	for _, node := range i.Nodes {
		nodes = append(nodes, node.Details.String())
	}
	return strings.Join(nodes, "/")
}

// "You wouldn't recognize valid Go syntax if it ran you clean through!"
// "Then watch me render every node as %#v, wrapped in a slice literal any compiler would applaud."
func (i DevicePathInstance) GoString() string {
	var nodes []string
	for _, node := range i.Nodes {
		nodes = append(nodes, fmt.Sprintf("%#v", node.Details))
	}
	return fmt.Sprintf("[]devicepath.DevicePathNode{%s}", strings.Join(nodes, ", "))
}

// "Indent your manners before you indent my instance, or I'll teach you both at once!"
// "Two extra spaces per level is all the courtesy your nodes will get from me."
func (i DevicePathInstance) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sDevicePathInstance\n", indent)

	for _, node := range i.Nodes {
		node.Details.dump(w, indent+"  ")
	}
}

// "An unrecognized node? I've charted stranger waters with worse maps than that!"
// "Unknown or not, I'll print its type, subtype, and raw bytes in hex so no mystery survives."
func (n DevicePathNode) String() string {
	return fmt.Sprintf("Unknown(%d,%d,%x)", n.Type, n.SubType, n.Data)
}

// "You think one flag ends our duel? It takes more than a whim to finish me!"
// "Aye, and only when Type is End and SubType marks the whole path's end do I call this fight over."
func isEndEntireDevicePath(node DevicePathNode) bool {
	return node.Type == DevicePathEnd && EndDevicePathSubType(node.SubType) == EndEntireDevicePathSubType
}

// "You'll need more than a single glance to tell an instance boundary from a full retreat!"
// "Type End paired with the ThisInstance subtype is exactly how I know one instance stops, not the whole path."
func isEndThisInstance(node DevicePathNode) bool {
	return node.Type == DevicePathEnd && EndDevicePathSubType(node.SubType) == EndThisInstanceSubType
}

// "Feed me a raw byte stream and I'll unravel it faster than you can draw your cutlass!"
// "Four bytes at a time, checked and bounded — I'll walk this buffer to the end without falling overboard."
func ParseDevicePath(data []byte) (*DevicePath, error) {
	var instances []DevicePathInstance
	var nodes []DevicePathNode

	for offset := 0; offset < len(data); {
		// "Short a header, are you? A coward's trick to sink my parser before it starts!"
		// "Not today — I demand a full four-byte type/subtype/length header before touching your data."
		if offset+4 > len(data) {
			return nil, fmt.Errorf(
				"truncated device path header at offset %d",
				offset,
			)
		}

		nodeLength := int(binary.LittleEndian.Uint16(data[offset+2 : offset+4]))

		// "You'd have me chase a node with less flesh on it than a wooden peg leg!"
		// "Every node needs at least its own four-byte header, or I toss the whole plank overboard."
		if nodeLength < 4 {
			return nil, fmt.Errorf(
				"invalid device path node length %d at offset %d",
				nodeLength,
				offset,
			)
		}

		nodeEnd := offset + nodeLength
		// "You'd have me read clean past the edge of the map and off into the abyss!"
		// "Not I — I check the node's end against the buffer's length before taking one more step."
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

// "Dressed in an unknown type, you thought you'd slip past unmarked and unparsed!"
// "Media, Messaging, ACPI, or otherwise — I check node.Type and hand you to the right parser, no imposters allowed."
func parseDevicePathNode(node DevicePathNode) (DevicePathNodeDetails, error) {
	switch node.Type {
	case DevicePathMedia:
		return parseMediaDevicePathNode(node)

	// case DevicePathHardware:
	// 	return parseHardwareDevicePathNode(node)

	case DevicePathACPI:
		return parseAcpiDevicePathNode(node)

	case DevicePathMessaging:
		return parseMessagingDevicePathNode(node)

	// case DevicePathBBS:
	// 	return parseBBSDevicePathNode(node)

	default:
		return unknownDevicePathNode(node), nil
	}
}
