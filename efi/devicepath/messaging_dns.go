package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/network"
	"io"
	"strings"
)

// "One flag can't hold both oceans of address, ye bilge rat!"
// "It doesn't try — IsIPv6 picks which of the two address lists actually holds water."
type DnsMessagingNode struct {
	IsIPv6        bool
	IPv4Addresses []network.IPv4Address
	IPv6Addresses []network.IPv6Address
}

// "Address a ghost and it'll haunt your call stack forever!"
// "Not this ghost — the nil receiver is caught up front and answers with a plain <nil> instead of a panic."
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

// "Nil or not, you'll answer for your Go syntax!"
// "It answers honestly, printing the literal nil pointer form before touching any field."
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

// "Pick a side of the ocean, DNS, or drown in your own ambiguity!"
// "It picks by IsIPv6, then walks only the matching list of servers, IPv6 or IPv4."
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

// "Empty holds and lone marker bytes both earn ye the depths!"
// "So both are rejected outright: at least one byte for the type flag, and at least one more for an actual address."
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
		// "Four bytes to a fleet, or the whole convoy sinks!"
		// "Each IPv4 address is exactly four bytes, so the remainder must divide evenly."
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
		// "Sixteen bytes make an IPv6 galleon, not a dinghy!"
		// "So the leftover after the type flag must split cleanly into sixteen-byte chunks."
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
