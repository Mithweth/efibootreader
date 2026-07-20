package network

import (
	"fmt"
	"net/netip"
)

type IPv4AddressType uint8

const (
	IPv4AddressDHCP   IPv4AddressType = 0
	IPv4AddressStatic IPv4AddressType = 1
)

func (t IPv4AddressType) String() string {
	switch t {
	case IPv4AddressDHCP:
		return "DHCP"
	case IPv4AddressStatic:
		return "Static"
	default:
		return fmt.Sprintf("%d", uint8(t))
	}
}

func (t IPv4AddressType) GoString() string {
	switch t {
	case IPv4AddressDHCP:
		return "network.IPv4AddressDHCP"
	case IPv4AddressStatic:
		return "network.IPv4AddressStatic"
	default:
		return fmt.Sprintf("network.IPv4AddressType(%d)", uint8(t))
	}
}

func ParseIPv4AddressType(data byte) IPv4AddressType {
	return IPv4AddressType(data)
}

type IPv4Address [4]byte

func (a IPv4Address) String() string {
	return netip.AddrFrom4([4]byte(a)).String()
}

func (a IPv4Address) GoString() string {
	return fmt.Sprintf("network.MustParseIPv4Address(%q)", a.String())
}

func ParseIPv4Address(data []byte) (IPv4Address, error) {
	if len(data) != 4 {
		return IPv4Address{}, fmt.Errorf(
			"expected 4 bytes, got %d",
			len(data),
		)
	}

	return IPv4Address{data[0], data[1], data[2], data[3]}, nil
}

func MustParseIPv4Address(s string) IPv4Address {
	addr, err := netip.ParseAddr(s)
	if err != nil {
		panic(fmt.Sprintf("network.MustParseIPv4Address(%q): %v", s, err))
	}

	if !addr.Is4() {
		panic(fmt.Sprintf(
			"network.MustParseIPv4Address(%q): not an IPv4 address",
			s,
		))
	}

	return IPv4Address(addr.As4())
}
