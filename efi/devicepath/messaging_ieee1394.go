package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Ieee1394MessagingNode struct {
	Reserved uint32
	GUID     uint64
}

func (h *Ieee1394MessagingNode) String() string {
	return fmt.Sprintf("I1394(%d)", h.GUID)
}

func (h *Ieee1394MessagingNode) GoString() string {
	if h == nil {
		return "(*efi.Ieee1394MessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&efi.Ieee1394MessagingNode{"+
			"Reserved:%#v, "+
			"GUID:%#v}",
		h.Reserved,
		h.GUID,
	)
}

func (h *Ieee1394MessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sIEEE 1394 Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Reserved\t : %d\n", indent, h.Reserved)
	_, _ = fmt.Fprintf(w, "%s  GUID\t\t : %d\n", indent, h.GUID)
}

func parseIeee1394MessagingNode(data []byte) (*Ieee1394MessagingNode, error) {
	if len(data) != 12 {
		return nil, fmt.Errorf(
			"invalid messaging IEEE 1394 node payload size: got %d, want 12",
			len(data),
		)
	}

	return &Ieee1394MessagingNode{
		Reserved: binary.LittleEndian.Uint32(data[0:4]),
		GUID:     binary.LittleEndian.Uint64(data[4:12]),
	}, nil
}
