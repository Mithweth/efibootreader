package network

import (
	"encoding/binary"
	"fmt"
)

type NetworkProtocol uint16

const (
	NetworkProtocolNil NetworkProtocol = 0
	NetworkProtocolTCP NetworkProtocol = 6
	NetworkProtocolUDP NetworkProtocol = 17
)

func (p NetworkProtocol) String() string {
	switch p {
	case NetworkProtocolTCP:
		return "TCP"
	case NetworkProtocolUDP:
		return "UDP"
	default:
		return fmt.Sprintf("%d", uint16(p))
	}
}

func (p NetworkProtocol) GoString() string {
	switch p {
	case NetworkProtocolTCP:
		return "efi.NetworkProtocolTCP"
	case NetworkProtocolUDP:
		return "efi.NetworkProtocolUDP"
	default:
		return fmt.Sprintf("efi.NetworkProtocol(%d)", uint16(p))
	}
}

func ParseNetworkProtocol(data []byte) (NetworkProtocol, error) {
	if len(data) != 2 {
		return NetworkProtocolNil, fmt.Errorf("expected 2 bytes, got %d", len(data))
	}
	return NetworkProtocol(binary.LittleEndian.Uint16(data)), nil
}
