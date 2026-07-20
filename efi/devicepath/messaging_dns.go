package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/network"
	"io"
	"strings"
)

type DnsMessagingNode struct {
	IsIPv6        bool
	IPv4Addresses []network.IPv4Address
	IPv6Addresses []network.IPv6Address
}

func (h *DnsMessagingNode) String() string {
	if h == nil {
		return "<nil>"
	}

	var addresses []string
	if h.IsIPv6 {
		for _, address := range h.IPv6Addresses {
			addresses = append(addresses, address.String())
		}
	} else {
		for _, address := range h.IPv4Addresses {
			addresses = append(addresses, address.String())
		}
	}
	return fmt.Sprintf("Dns(%s)", strings.Join(addresses, ","))
}

func (h *DnsMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.DnsMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.DnsMessagingNode{"+
			"IsIPv6:%#v, "+
			"IPv4Addresses:%#v, "+
			"IPv6Addresses:%#v}",
		h.IsIPv6,
		h.IPv4Addresses,
		h.IPv6Addresses,
	)
}

func (h *DnsMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sDNS Messaging Node\n", indent)
	if h.IsIPv6 {
		_, _ = fmt.Fprintf(w, "%s  Address Type\t : IPv6\n", indent)
		for _, address := range h.IPv6Addresses {
			_, _ = fmt.Fprintf(w, "%s  DNS Server\t : %s\n", indent, address)
		}
	} else {
		_, _ = fmt.Fprintf(w, "%s  Address Type\t : IPv4\n", indent)
		for _, address := range h.IPv4Addresses {
			_, _ = fmt.Fprintf(w, "%s  DNS Server\t : %s\n", indent, address)
		}
	}
}

func parseDnsMessagingNode(data []byte) (*DnsMessagingNode, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf(
			"invalid messaging DNS node payload size: got %d, want at least 1",
			len(data),
		)
	}
	if len(data) == 1 {
		return nil, fmt.Errorf(
			"invalid messaging DNS node payload size: no DNS address",
		)
	}

	switch data[0] {
	case 0:
		if (len(data)-1)%4 != 0 {
			return nil, fmt.Errorf(
				"invalid messaging DNS IPv4 node payload size: got %d",
				len(data),
			)
		}

		var addresses []network.IPv4Address

		for offset := 1; offset < len(data); offset += 4 {
			address, err := network.ParseIPv4Address(data[offset : offset+4])
			if err != nil {
				return nil, fmt.Errorf(
					"parse DNS IPv4 address: %w",
					err,
				)
			}

			addresses = append(addresses, address)
		}

		return &DnsMessagingNode{IsIPv6: false, IPv4Addresses: addresses}, nil

	case 1:
		if (len(data)-1)%16 != 0 {
			return nil, fmt.Errorf(
				"invalid messaging DNS IPv6 node payload size: got %d",
				len(data),
			)
		}

		var addresses []network.IPv6Address

		for offset := 1; offset < len(data); offset += 16 {
			address, err := network.ParseIPv6Address(data[offset : offset+16])
			if err != nil {
				return nil, fmt.Errorf(
					"parse DNS IPv6 address: %w",
					err,
				)
			}

			addresses = append(addresses, address)
		}

		return &DnsMessagingNode{IsIPv6: true, IPv6Addresses: addresses}, nil

	default:
		return nil, fmt.Errorf(
			"invalid messaging DNS address type: got %d, want 0 or 1",
			data[0],
		)
	}
}
