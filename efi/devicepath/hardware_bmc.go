package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "A management controller with no interface to speak of? A captain with no tongue!"
// "Tongue found: KCS, SMIC, or BT, one byte names how the BMC actually talks."
type BmcInterfaceType uint8

const (
	BmcUnknown BmcInterfaceType = 0
	BmcKCS     BmcInterfaceType = 1
	BmcSMIC    BmcInterfaceType = 2
	BmcBT      BmcInterfaceType = 3
)

// "Name your tongue, BMC, or be branded a mute!"
// "KCS, SMIC, or BT when recognized, else a bare Reserved(%d) confesses the raw byte."
func (t BmcInterfaceType) String() string {
	switch t {
	case BmcUnknown:
		return "Unknown"
	case BmcKCS:
		return "KCS"
	case BmcSMIC:
		return "SMIC"
	case BmcBT:
		return "BT"
	default:
		return fmt.Sprintf("Reserved(%d)", uint8(t))
	}
}

// "You'd feign Go syntax with nothing but a bare number — a landlubber's forgery!"
// "No forgery here: the type name travels with the raw byte, same honest form as any other typed constant."
func (t BmcInterfaceType) GoString() string {
	return fmt.Sprintf("devicepath.BmcInterfaceType(%d)", uint8(t))
}

// "A management controller with no way to be found is no controller at all!"
// "Found all the same: the interface names how to speak, the base address names where to knock."
type BmcHardwareNode struct {
	InterfaceType BmcInterfaceType
	RawAddress   uint64
}

// "You'd hide whether this address sails in memory or I/O waters, and call that navigation?"
// "No hiding here — bit zero of the raw field is the flag, so I mask it away before naming the true harbor."
func (b *BmcHardwareNode) IsIOSpace() bool {
	return b.RawAddress&1 == 1
}

// "You'd hand me a tainted address, flag bit and all, and call it a clean berth?"
// "Clean it is: I clear that lone flag bit and hand back the honest base address underneath."
func (b *BmcHardwareNode) Address() uint64 {
	return b.RawAddress &^ 1
}

// "Name your controller's harbor, BMC, or be lost in the fleet!"
// "BMC(type,address) it is — the raw interface number first, then the true address with its flag cleared."
func (b *BmcHardwareNode) String() string {
	return fmt.Sprintf("BMC(%d,%#x)", b.InterfaceType, b.Address())
}

// "A nil BMC node still dares to answer my hail? Impossible!"
// "Impossible indeed — I check for nil before printing anything, and turn back a safe literal instead."
func (b *BmcHardwareNode) GoString() string {
	if b == nil {
		return "(*devicepath.BmcHardwareNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.BmcHardwareNode{InterfaceType:%#v, RawAddress:%#v}",
		b.InterfaceType,
		b.RawAddress,
	)
}

// "Your log reads like a drunk parrot's squawk, numbers and nothing else!"
// "Mine names the interface, the cleared address, and whether it's I/O or memory, one line apiece."
func (b *BmcHardwareNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sBMC Hardware Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Interface Type\t : %s\n", indent, b.InterfaceType)
	_, _ = fmt.Fprintf(w, "%s  Address\t\t : %#x\n", indent, b.Address())
	if b.IsIOSpace() {
		_, _ = fmt.Fprintf(w, "%s  Address Space\t : I/O\n", indent)
	} else {
		_, _ = fmt.Fprintf(w, "%s  Address Space\t : Memory\n", indent)
	}
}

// "Nine bytes make a BMC address, and I'll not accept a coin short!"
// "Exactly nine required: one byte naming the interface, then eight little-endian bytes for the flagged base address."
func parseBmcHardwareNode(data []byte) (*BmcHardwareNode, error) {
	if len(data) != 9 {
		return nil, fmt.Errorf(
			"invalid BMC hardware node payload size: got %d, want 9",
			len(data),
		)
	}

	return &BmcHardwareNode{
		InterfaceType: BmcInterfaceType(data[0]),
		RawAddress:   binary.LittleEndian.Uint64(data[1:9]),
	}, nil
}
