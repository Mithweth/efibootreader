package devicepath

import (
	"encoding/binary"
	"fmt"
	"github.com/Mithweth/efibootreader/network"
	"io"
)

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
