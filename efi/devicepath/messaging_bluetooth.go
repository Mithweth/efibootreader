package devicepath

import (
	"fmt"
	"io"
)

type BluetoothMessagingNode struct {
	DeviceAddress [6]byte
}

func (f *BluetoothMessagingNode) String() string {
	return fmt.Sprintf("Bluetooth(0x%x)", f.DeviceAddress)
}

func (f *BluetoothMessagingNode) GoString() string {
	if f == nil {
		return "(*devicepath.BluetoothMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.BluetoothMessagingNode{DeviceAddress:%#v}", f.DeviceAddress)
}

func (f *BluetoothMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sBluetooth Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Device Address\t : %x\n", indent, f.DeviceAddress)
}

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
