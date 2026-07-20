package network

import (
	"fmt"
	"net/netip"
)

type IPv6Address [16]byte

func (a IPv6Address) String() string {
	return netip.AddrFrom16([16]byte(a)).String()
}

func (a IPv6Address) GoString() string {
	return fmt.Sprintf("network.MustParseIPv6Address(%q)", a.String())
}

func ParseIPv6Address(data []byte) (IPv6Address, error) {
	if len(data) != 16 {
		return IPv6Address{}, fmt.Errorf(
			"expected 16 bytes, got %d",
			len(data),
		)
	}

	var address IPv6Address
	copy(address[:], data)

	return address, nil
}

func MustParseIPv6Address(s string) IPv6Address {
	addr, err := netip.ParseAddr(s)
	if err != nil {
		panic(fmt.Sprintf("network.MustParseIPv6Address(%q): %v", s, err))
	}

	if !addr.Is6() {
		panic(fmt.Sprintf(
			"network.MustParseIPv6Address(%q): not an IPv6 address",
			s,
		))
	}

	return IPv6Address(addr.As16())
}
