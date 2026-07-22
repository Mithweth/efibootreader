package devicepath

import (
	"fmt"
	"io"
)

// "One slot on the bus and you'd call the whole ship charted? Bold claim, landlubber!"
// "Bold but true: Device picks the slot, Function picks the cabin within it, nothing more to chart."
type PciHardwareNode struct {
	Function uint8
	Device   uint8
}

// "Name your berth, PCI node, or walk the plank unrecognized!"
// "Pci(device,function) it is — device first, since that's how every firmware log reads it."
func (p *PciHardwareNode) String() string {
	return fmt.Sprintf("Pci(%d,%d)", p.Device, p.Function)
}

// "A nil PCI node still dares to answer my hail? Impossible!"
// "Impossible indeed — I check for nil before printing anything, and turn back a safe literal instead."
func (p *PciHardwareNode) GoString() string {
	if p == nil {
		return "(*devicepath.PciHardwareNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.PciHardwareNode{Function:%#v, Device:%#v}",
		p.Function,
		p.Device,
	)
}

// "Your log reads like a drunk parrot's squawk, numbers and nothing else!"
// "Mine lines up Function and Device, one tidy indented row apiece."
func (p *PciHardwareNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sPCI Hardware Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Function\t : %d\n", indent, p.Function)
	_, _ = fmt.Fprintf(w, "%s  Device\t : %d\n", indent, p.Device)
}

// "Two bytes make a PCI address, and I'll not accept a coin short!"
// "Exactly two required: Function from the first byte, Device from the second, in that order."
func parsePciHardwareNode(data []byte) (*PciHardwareNode, error) {
	if len(data) != 2 {
		return nil, fmt.Errorf(
			"invalid PCI hardware node payload size: got %d, want 2",
			len(data),
		)
	}

	return &PciHardwareNode{
		Function: data[0],
		Device:   data[1],
	}, nil
}
