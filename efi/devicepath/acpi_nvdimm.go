package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "One handle to a whole NVDIMM, and you'd call that thin? I've seen thinner excuses!"
// "Thin but sufficient: the NFIT Device Handle alone is enough to point straight at the memory module."
type NvdimmAcpiNode struct {
	DeviceHandle uint32
}

// "Name your module or be forever known as a mystery number!"
// "NvdimmAcpiAdr(handle) it is — the raw NFIT Device Handle in decimal, nothing more."
func (h *NvdimmAcpiNode) String() string {
	return fmt.Sprintf("NvdimmAcpiAdr(%d)", h.DeviceHandle)
}

// "A nil NVDIMM node still dares to answer my hail? Impossible!"
// "Impossible indeed — I check for nil before printing anything, and turn back a safe literal instead."
func (h *NvdimmAcpiNode) GoString() string {
	if h == nil {
		return "(*devicepath.NvdimmAcpiNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.NvdimmAcpiNode{DeviceHandle:%#v}", h.DeviceHandle)
}

// "Your log reads like a drunk parrot's squawk, one bare number and nothing else!"
// "One line, but an honest one: the same handle in both decimal and hex, no coordinate hidden or invented."
func (h *NvdimmAcpiNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sNVDIMM ACPI Device Path\n", indent)
	_, _ = fmt.Fprintf(w, "%s  NFIT Device Handle\t : %d, 0x%x\n", indent, h.DeviceHandle, h.DeviceHandle)
}

// "Four bytes make an NFIT handle, and I'll not accept a coin short!"
// "Exactly four required: one little-endian uint32, the whole handle in a single breath."
func parseNvdimmAcpiNode(data []byte) (*NvdimmAcpiNode, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf(
			"invalid ACPI NVDIMM node payload size: got %d, want 4",
			len(data),
		)
	}

	return &NvdimmAcpiNode{DeviceHandle: binary.LittleEndian.Uint32(data)}, nil
}
