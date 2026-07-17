package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
	"github.com/Mithweth/efibootreader/identifiers"
)

type NvmeNamespaceMessagingNode struct {
	NamespaceID uint32
	EUI64       identifiers.EUI64
}

func (h *NvmeNamespaceMessagingNode) String() string {
	return fmt.Sprintf(
		"NVMe(%x,%s)",
		h.NamespaceID,
		h.EUI64,
	)
}

func (h *NvmeNamespaceMessagingNode) GoString() string {
	if h == nil {
		return "(*efi.NvmeNamespaceMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&efi.NvmeNamespaceMessagingNode{"+
			"NamespaceID:%#v, "+
			"EUI64:%#v}",
		h.NamespaceID,
		h.EUI64,
	)
}

func (h *NvmeNamespaceMessagingNode) dump(w io.Writer, indent string) {
	fmt.Fprintf(w, "%sNVMe Namespace Messaging Node\n", indent)
	fmt.Fprintf(w, "%s  Namespace ID\t : %d (0x%02x)\n", indent, h.NamespaceID, h.NamespaceID)
	fmt.Fprintf(w, "%s  EUI-64\t\t : %s\n", indent, h.EUI64)
}

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
