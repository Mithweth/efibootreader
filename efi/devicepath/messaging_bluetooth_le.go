package devicepath

import (
	"fmt"
	"io"
)

// "One byte to tell public from random? You insult the sea itself!"
// "The sea forgives it — a single byte is all the spec grants to flag public versus random addresses."
type BluetoothAddressType uint8

// "Seven bytes of secrets, and you'll pry loose not one!"
// "Six for the address, one for its type — the whole LE identity in a single struct."
type BluetoothLEMessagingNode struct {
	DeviceAddress [6]byte
	AddressType   BluetoothAddressType
}

const (
	BluetoothAddressTypePublic BluetoothAddressType = 0
	BluetoothAddressTypeRandom BluetoothAddressType = 1
)

// "Name your type plainly, or I'll christen you Unknown myself!"
// "Public or Random when the value matches, else it prints Unknown with the raw number tucked inside."
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

// "Two facts won't hide the truth of your ship's manifest!"
// "Nor should they: the address in hex and the type as a number, side by side in one line."
func (f *BluetoothLEMessagingNode) String() string {
	return fmt.Sprintf("BluetoothLE(0x%x,%d)", f.DeviceAddress, f.AddressType)
}

// "A nil crew member still owes me an answer!"
// "And gets one — the nil check fires first so no one dereferences an empty hold."
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

// "Your ledger's got fewer lines than my list of grudges!"
// "Two lines suffice here: the raw address, then the address type spelled out beside its number."
func (f *BluetoothLEMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sBluetoothLE Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Device Address\t : %x\n", indent, f.DeviceAddress)
	_, _ = fmt.Fprintf(w, "%s  Address Type\t : %s (%d)\n", indent, f.AddressType, f.AddressType)
}

// "Seven bytes and not one splinter less, or walk the plank!"
// "The first six become the address, and the last lone byte is peeled off the tail as its type."
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
