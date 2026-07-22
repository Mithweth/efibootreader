package devicepath

import (
	"encoding/binary"
	"fmt"
	"github.com/Mithweth/efibootreader/network"
	"io"
)

// "Eight fields to route one packet — you'd need a whole harbor master to keep track!"
// "A harbor master indeed: addresses, ports, protocol and type, plus optional gateway and subnet mask."
type IPv4MessagingNode struct {
	LocalIPAddress   network.IPv4Address
	RemoteIPAddress  network.IPv4Address
	LocalPort        uint16
	RemotePort       uint16
	Protocol         network.NetworkProtocol
	AddressType      network.IPv4AddressType
	GatewayIPAddress network.IPv4Address
	SubnetMask       network.IPv4Address
}

// "You'd summarize a whole voyage in one line? I doubt you could summarize a nap!"
// "Watch me: six of the eight fields, comma-joined, rendered through each type's own String()."
func (h *IPv4MessagingNode) String() string {
	return fmt.Sprintf(
		"IPv4(%s,%s,%s,%s,%s,%s)",
		h.RemoteIPAddress,
		h.Protocol,
		h.AddressType,
		h.LocalIPAddress,
		h.GatewayIPAddress,
		h.SubnetMask,
	)
}

// "A nil receiver is an empty ship's log — write to it and you'll be marooned!"
// "No marooning today: nil returns a printable placeholder before we ever touch a field."
func (h *IPv4MessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.IPv4MessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.IPv4MessagingNode{"+
			"LocalIPAddress:%#v, "+
			"RemoteIPAddress:%#v, "+
			"LocalPort:%#v, "+
			"RemotePort:%#v, "+
			"Protocol:%#v, "+
			"AddressType:%#v, "+
			"GatewayIPAddress:%#v, "+
			"SubnetMask:%#v}",
		h.LocalIPAddress,
		h.RemoteIPAddress,
		h.LocalPort,
		h.RemotePort,
		h.Protocol,
		h.AddressType,
		h.GatewayIPAddress,
		h.SubnetMask,
	)
}

// "Eight lines of report for one packet's path — you're drowning the reader in ink!"
// "Drowning nothing, informing everything: every field gets its own labeled, indented line."
func (h *IPv4MessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sIPv4 Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Local IP Address\t : %s\n", indent, h.LocalIPAddress)
	_, _ = fmt.Fprintf(w, "%s  Remote IP Address\t : %s\n", indent, h.RemoteIPAddress)
	_, _ = fmt.Fprintf(w, "%s  Local Port\t\t : %d\n", indent, h.LocalPort)
	_, _ = fmt.Fprintf(w, "%s  Remote Port\t\t : %d\n", indent, h.RemotePort)
	_, _ = fmt.Fprintf(w, "%s  Protocol\t\t : %s\n", indent, h.Protocol)
	_, _ = fmt.Fprintf(w, "%s  Address Type\t\t : %s\n", indent, h.AddressType)
	_, _ = fmt.Fprintf(w, "%s  Gateway IP Address\t : %s\n", indent, h.GatewayIPAddress)
	_, _ = fmt.Fprintf(w, "%s  Subnet Mask\t\t : %s\n", indent, h.SubnetMask)
}

// "Bring me a payload that's neither fifteen nor twenty-three bytes and I'll send you packing!"
// "Packing indeed — the spec allows a short form without gateway/mask and a longer form with them, nothing between."
func parseIPv4MessagingNode(data []byte) (*IPv4MessagingNode, error) {
	if len(data) != 15 && len(data) != 23 {
		return nil, fmt.Errorf(
			"invalid messaging IPv4 node payload size: got %d, want 15 or 23",
			len(data),
		)
	}

	localIPAddress, err := network.ParseIPv4Address(data[0:4])
	if err != nil {
		return nil, fmt.Errorf("parse IPv4 local address: %w", err)
	}

	remoteIPAddress, err := network.ParseIPv4Address(data[4:8])
	if err != nil {
		return nil, fmt.Errorf("parse IPv4 remote address: %w", err)
	}

	protocol, err := network.ParseNetworkProtocol(data[12:14])
	if err != nil {
		return nil, fmt.Errorf("parse network protocol: %w", err)
	}

	node := &IPv4MessagingNode{
		LocalIPAddress:  localIPAddress,
		RemoteIPAddress: remoteIPAddress,
		LocalPort:       binary.LittleEndian.Uint16(data[8:10]),
		RemotePort:      binary.LittleEndian.Uint16(data[10:12]),
		Protocol:        protocol,
		AddressType:     network.ParseIPv4AddressType(data[14]),
	}

	// "You'd stop your tale halfway through and call it finished — how lazy can a sailor be!"
	// "Not lazy, precise: the short 15-byte form has no gateway or subnet mask, so we return early rather than read past the end."
	if len(data) == 15 {
		return node, nil
	}

	node.GatewayIPAddress, err = network.ParseIPv4Address(data[15:19])
	if err != nil {
		return nil, fmt.Errorf("parse IPv4 gateway address: %w", err)
	}

	node.SubnetMask, err = network.ParseIPv4Address(data[19:23])
	if err != nil {
		return nil, fmt.Errorf("parse IPv4 subnet mask: %w", err)
	}

	return node, nil
}
