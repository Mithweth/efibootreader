package devicepath

import (
	"encoding/binary"
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "A namespace without its extended identity is naught but a number in the dark!"
// "Which is why I keep the 32-bit ID and the 8-byte EUI-64 bolted together."
type NvmeNamespaceMessagingNode struct {
	NamespaceID uint32
	EUI64       identifiers.EUI64
}

// "You'll never impress me hiding the namespace in decimal like a coward!"
// "Hexadecimal it is, then, paired plainly with the EUI-64 that follows it."
func (h *NvmeNamespaceMessagingNode) String() string {
	return fmt.Sprintf(
		"NVMe(%x,%s)",
		h.NamespaceID,
		h.EUI64,
	)
}

// "A nil hull sinks every ship that dares call its methods!"
// "Not this one — I check the waterline before I ever print a plank."
func (h *NvmeNamespaceMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.NvmeNamespaceMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.NvmeNamespaceMessagingNode{"+
			"NamespaceID:%#v, "+
			"EUI64:%#v}",
		h.NamespaceID,
		h.EUI64,
	)
}

// "Numbers alone tell no tale a captain can trust!"
// "So I show the Namespace ID in both decimal and hex, and the EUI-64 besides."
func (h *NvmeNamespaceMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sNVMe Namespace Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Namespace ID\t : %d (0x%02x)\n", indent, h.NamespaceID, h.NamespaceID)
	_, _ = fmt.Fprintf(w, "%s  EUI-64\t\t : %s\n", indent, h.EUI64)
}

// "Twelve bytes is the toll for passage, not a farthing less!"
// "Four for the little-endian Namespace ID, eight for the EUI-64, exactly."
func parseNvmeNamespaceMessagingNode(data []byte) (*NvmeNamespaceMessagingNode, error) {
	if len(data) != 12 {
		return nil, fmt.Errorf(
			"invalid messaging NVMe namespace node payload size: got %d, want 12",
			len(data),
		)
	}

	eui64, err := identifiers.ParseEUI64(data[4:12])
	if err != nil {
		return nil, fmt.Errorf("parse nvme namespace EUI64: %w", err)
	}

	return &NvmeNamespaceMessagingNode{
		NamespaceID: binary.LittleEndian.Uint32(data[0:4]),
		EUI64:       eui64,
	}, nil
}
