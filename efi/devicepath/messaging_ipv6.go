package devicepath

import (
	"encoding/binary"
	"fmt"
	"github.com/Mithweth/efibootreader/network"
	"io"
)

type IPv6AddressOrigin uint8

const (
	IPv6AddressManual                 IPv6AddressOrigin = 0
	IPv6AddressStatelessAutoConfigure IPv6AddressOrigin = 1
	IPv6AddressStatefulAutoConfigure  IPv6AddressOrigin = 2
)

func (o IPv6AddressOrigin) String() string {
	switch o {
	case IPv6AddressManual:
		return "Static"
	case IPv6AddressStatelessAutoConfigure:
		return "StatelessAutoConfigure"
	case IPv6AddressStatefulAutoConfigure:
		return "StatefulAutoConfigure"
	default:
		return fmt.Sprintf("%d", uint8(o))
	}
}

func (o IPv6AddressOrigin) GoString() string {
	switch o {
	case IPv6AddressManual:
		return "devicepath.IPv6AddressManual"
	case IPv6AddressStatelessAutoConfigure:
		return "devicepath.IPv6AddressStatelessAutoConfigure"
	case IPv6AddressStatefulAutoConfigure:
		return "devicepath.IPv6AddressStatefulAutoConfigure"
	default:
		return fmt.Sprintf("devicepath.IPv6AddressOrigin(%d)", uint8(o))
	}
}

type IPv6MessagingNode struct {
	LocalIPAddress   network.IPv6Address
	RemoteIPAddress  network.IPv6Address
	LocalPort        uint16
	RemotePort       uint16
	Protocol         network.NetworkProtocol
	AddressOrigin    IPv6AddressOrigin
	PrefixLength     uint8
	GatewayIPAddress network.IPv6Address
}

func (h *IPv6MessagingNode) String() string {
	return fmt.Sprintf(
		"IPv6(%s,%s,%s,%s,%d,%s)",
		h.RemoteIPAddress,
		h.Protocol,
		h.AddressOrigin,
		h.LocalIPAddress,
		h.PrefixLength,
		h.GatewayIPAddress,
	)
}

func (h *IPv6MessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.IPv6MessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.IPv6MessagingNode{"+
			"LocalIPAddress:%#v, "+
			"RemoteIPAddress:%#v, "+
			"LocalPort:%#v, "+
			"RemotePort:%#v, "+
			"Protocol:%#v, "+
			"AddressOrigin:%#v, "+
			"PrefixLength:%#v, "+
			"GatewayIPAddress:%#v}",
		h.LocalIPAddress,
		h.RemoteIPAddress,
		h.LocalPort,
		h.RemotePort,
		h.Protocol,
		h.AddressOrigin,
		h.PrefixLength,
		h.GatewayIPAddress,
	)
}

func (h *IPv6MessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sIPv6 Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Local IP Address\t : %s\n", indent, h.LocalIPAddress)
	_, _ = fmt.Fprintf(w, "%s  Remote IP Address\t : %s\n", indent, h.RemoteIPAddress)
	_, _ = fmt.Fprintf(w, "%s  Local Port\t\t : %d\n", indent, h.LocalPort)
	_, _ = fmt.Fprintf(w, "%s  Remote Port\t\t : %d\n", indent, h.RemotePort)
	_, _ = fmt.Fprintf(w, "%s  Protocol\t\t : %s\n", indent, h.Protocol)
	_, _ = fmt.Fprintf(w, "%s  Address Origin\t : %s\n", indent, h.AddressOrigin)
	_, _ = fmt.Fprintf(w, "%s  Prefix Length\t : %d\n", indent, h.PrefixLength)
	_, _ = fmt.Fprintf(w, "%s  Gateway IP Address\t : %s\n", indent, h.GatewayIPAddress)
}

func parseIPv6MessagingNode(data []byte) (*IPv6MessagingNode, error) {
	if len(data) != 56 {
		return nil, fmt.Errorf(
			"invalid messaging IPv6 node payload size: got %d, want 56",
			len(data),
		)
	}

	localIPAddress, err := network.ParseIPv6Address(data[0:16])
	if err != nil {
		return nil, fmt.Errorf("parse IPv6 local address: %w", err)
	}

	remoteIPAddress, err := network.ParseIPv6Address(data[16:32])
	if err != nil {
		return nil, fmt.Errorf("parse IPv6 remote address: %w", err)
	}

	gatewayIPAddress, err := network.ParseIPv6Address(data[40:56])
	if err != nil {
		return nil, fmt.Errorf("parse IPv6 gateway address: %w", err)
	}

	protocol, err := network.ParseNetworkProtocol(data[36:38])
	if err != nil {
		return nil, fmt.Errorf("parse network protocol: %w", err)
	}

	return &IPv6MessagingNode{
		LocalIPAddress:   localIPAddress,
		RemoteIPAddress:  remoteIPAddress,
		LocalPort:        binary.LittleEndian.Uint16(data[32:34]),
		RemotePort:       binary.LittleEndian.Uint16(data[34:36]),
		Protocol:         protocol,
		AddressOrigin:    IPv6AddressOrigin(data[38]),
		PrefixLength:     data[39],
		GatewayIPAddress: gatewayIPAddress,
	}, nil
}
