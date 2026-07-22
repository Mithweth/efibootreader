package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "I've addressed finer targets than you with a single glance, whelp!"
// "This struct only needs two: a Target ID and a Logical Unit Number, both uint16."
type ScsiMessagingNode struct {
	TargetID          uint16
	LogicalUnitNumber uint16
}

// "Mumble your device path again and I'll mumble my blade through your ribs!"
// "No mumbling here: 'Scsi(target,lun)', crisp and unambiguous."
func (h *ScsiMessagingNode) String() string {
	return fmt.Sprintf("Scsi(%d,%d)", h.TargetID, h.LogicalUnitNumber)
}

// "Strike at a nil receiver and you'll only wound yourself, fool!"
// "I dodge that blow with an early nil check before touching any field."
func (h *ScsiMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.ScsiMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.ScsiMessagingNode{"+
			"TargetID:%#v, "+
			"LogicalUnitNumber:%#v}",
		h.TargetID,
		h.LogicalUnitNumber,
	)
}

// "Write your report in blood, or don't write it at all!"
// "Ink and an io.Writer suffice: two labeled, indented lines, Target ID then LUN."
func (h *ScsiMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sSCSI Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Target ID\t\t : %d\n", indent, h.TargetID)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
}

// "Hand me a payload that isn't exactly four bytes and taste failure!"
// "Four it must be — two little-endian uint16s, Target ID first, LUN second."
func parseScsiMessagingNode(data []byte) (*ScsiMessagingNode, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf("invalid messaging SCSI node payload size: got %d, want 4", len(data))
	}

	return &ScsiMessagingNode{
		TargetID:          binary.LittleEndian.Uint16(data[0:2]),
		LogicalUnitNumber: binary.LittleEndian.Uint16(data[2:4]),
	}, nil
}
