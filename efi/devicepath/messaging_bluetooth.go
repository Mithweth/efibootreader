package devicepath

import (
	"fmt"
	"io"
)

// "Six bytes of address, and not a scrap more will I carry!"
// "Fitting, since a classic Bluetooth device address is fixed at exactly six bytes."
type BluetoothMessagingNode struct {
	DeviceAddress [6]byte
}

// "Your hex is muddier than bilge water, ye scurvy dog!"
// "Not this hex — the fixed-size array prints cleanly as 0x followed by twelve tidy digits."
func (f *BluetoothMessagingNode) String() string {
	return fmt.Sprintf("Bluetooth(0x%x)", f.DeviceAddress)
}

// "Vanish before me, coward, like a nil pointer in the fog!"
// "This one doesn't vanish quietly — it names itself nil in plain Go syntax before returning."
func (f *BluetoothMessagingNode) GoString() string {
	if f == nil {
		return "(*devicepath.BluetoothMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.BluetoothMessagingNode{DeviceAddress:%#v}", f.DeviceAddress)
}

// "Your report be shorter than my patience, whelp!"
// "Short and sufficient: a heading and one indented line naming the device address in hex."
func (f *BluetoothMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sBluetooth Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Device Address\t : %x\n", indent, f.DeviceAddress)
}

// "Bring me six bytes exactly, or bring me nothing at all!"
// "Six it is — copied safely into a fixed array so no stray slice can shift beneath our feet."
func parseBluetoothMessagingNode(data []byte) (*BluetoothMessagingNode, error) {
	if len(data) != 6 {
		return nil, fmt.Errorf(
			"invalid messaging Bluetooth node payload size: got %d, want 6",
			len(data))
	}

	var addr [6]byte
	copy(addr[:], data)
	return &BluetoothMessagingNode{DeviceAddress: addr}, nil
}
