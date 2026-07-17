package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ScsiMessagingNode struct {
	TargetID          uint16
	LogicalUnitNumber uint16
}

func (h *ScsiMessagingNode) String() string {
	return fmt.Sprintf("Scsi(%d,%d)", h.TargetID, h.LogicalUnitNumber)
}

func (h *ScsiMessagingNode) GoString() string {
	if h == nil {
		return "(*efi.ScsiMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&efi.ScsiMessagingNode{"+
			"TargetID:%#v, "+
			"LogicalUnitNumber:%#v}",
		h.TargetID,
		h.LogicalUnitNumber,
	)
}

func (h *ScsiMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sSCSI Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Target ID\t\t : %d\n", indent, h.TargetID)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
}

func parseScsiMessagingNode(data []byte) (*ScsiMessagingNode, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf("invalid messaging SCSI node payload size: got %d, want 4", len(data))
	}

	return &ScsiMessagingNode{
		TargetID:          binary.LittleEndian.Uint16(data[0:2]),
		LogicalUnitNumber: binary.LittleEndian.Uint16(data[2:4]),
	}, nil
}
