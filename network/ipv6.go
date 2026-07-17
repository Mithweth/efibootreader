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
	return fmt.Sprintf("efi.IPv6Address{%s}", a)
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
