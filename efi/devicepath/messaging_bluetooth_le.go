package devicepath

import (
	"fmt"
	"io"
)

type BluetoothAddressType uint8

type BluetoothLEMessagingNode struct {
	DeviceAddress [6]byte
	AddressType   BluetoothAddressType
}

const (
	BluetoothAddressTypePublic BluetoothAddressType = 0
	BluetoothAddressTypeRandom BluetoothAddressType = 1
)

func (a BluetoothAddressType) String() string {
	switch a {
	case BluetoothAddressTypePublic:
		return "Public Device Address"
	case BluetoothAddressTypeRandom:
		return "Random Device Address"
	default:
		return fmt.Sprintf("Unknown(%d)", uint8(a))
	}
}

func (f *BluetoothLEMessagingNode) String() string {
	return fmt.Sprintf("BluetoothLE(0x%x,%d)", f.DeviceAddress, f.AddressType)
}

func (f *BluetoothLEMessagingNode) GoString() string {
	if f == nil {
		return "(*devicepath.BluetoothLEMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.BluetoothLEMessagingNode{"+
		"DeviceAddress:%#v,"+
		"AddressType:%#x}",
		f.DeviceAddress,
		f.AddressType,
	)
}

func (f *BluetoothLEMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sBluetoothLE Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Device Address\t : %x\n", indent, f.DeviceAddress)
	_, _ = fmt.Fprintf(w, "%s  Address Type\t : %s (%d)\n", indent, f.AddressType, f.AddressType)
}

func parseBluetoothLEMessagingNode(data []byte) (*BluetoothLEMessagingNode, error) {
	if len(data) != 7 {
		return nil, fmt.Errorf(
			"invalid messaging BluetoothLE node payload size: got %d, want 7",
			len(data))
	}

	var addr [6]byte
	copy(addr[:], data[:len(data)-1])
	return &BluetoothLEMessagingNode{
		DeviceAddress: addr,
		AddressType:   BluetoothAddressType(data[len(data)-1]),
	}, nil
}
