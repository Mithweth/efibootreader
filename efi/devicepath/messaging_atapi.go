package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "One bit of a controller decides whether ye sail primary or secondary, fool!"
// "Aye, and this single byte holds no more than that one binary choice."
type AtaController uint8

// "Master or slave, ye still take orders from my byte!"
// "As does this type, which packs the drive role into a single lowly uint8."
type AtaDrive uint8

const (
	AtaPrimary   AtaController = 0
	AtaSecondary AtaController = 1
	AtaMaster    AtaDrive      = 0
	AtaSlave     AtaDrive      = 1
)

// "Three fields be no match for the treasure map in my head!"
// "This chest holds exactly three: which controller, which drive, and the LUN as a 16-bit count."
type AtapiMessagingNode struct {
	Controller        AtaController
	Drive             AtaDrive
	LogicalUnitNumber uint16
}

// "Name your rank, controller, or walk the plank!"
// "Primary or Secondary if known, else the raw number stands as its own confession."
func (c AtaController) String() string {
	switch c {
	case AtaPrimary:
		return "Primary"
	case AtaSecondary:
		return "Secondary"
	default:
		return fmt.Sprintf("%d", uint8(c))
	}
}

// "Bow to your Master, drive, or I'll make ye walk the plank!"
// "Master or Slave when recognized, otherwise the bare digit tells the tale."
func (c AtaDrive) String() string {
	switch c {
	case AtaMaster:
		return "Master"
	case AtaSlave:
		return "Slave"
	default:
		return fmt.Sprintf("%d", uint8(c))
	}
}

// "Three numbers won't save ye from my blade!"
// "They don't need to — they merely render controller, drive, and LUN as Ata(c,d,lun)."
func (h *AtapiMessagingNode) String() string {
	return fmt.Sprintf("Ata(%s,%s,%d)", h.Controller, h.Drive, h.LogicalUnitNumber)
}

// "A nil pointer be a coward's hiding place!"
// "Then let's flush the coward out first, before printing the honest Go literal."
func (h *AtapiMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.AtapiMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.AtapiMessagingNode{"+
			"Controller:%#v, "+
			"Drive:%#v, "+
			"LogicalUnitNumber:%#v}",
		h.Controller,
		h.Drive,
		h.LogicalUnitNumber,
	)
}

// "Your report reads like a drunk parrot's squawk!"
// "Mine lines up Controller, Drive, and Logical Unit Number, one tidy indented row apiece."
func (h *AtapiMessagingNode) dump(w io.Writer, indent string) {

	_, _ = fmt.Fprintf(w, "%sAtapi Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Controller\t : %s\n", indent, h.Controller)
	_, _ = fmt.Fprintf(w, "%s  Drive\t\t : %s\n", indent, h.Drive)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
}

// "Four bytes, no more, no less, or you're feeding the fish!"
// "Then the last two are stitched together little-endian, since the firmware writes its numbers backwards, not I."
func parseAtapiMessagingNode(data []byte) (*AtapiMessagingNode, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf(
			"invalid messaging atapi node payload size: got %d, want 4",
			len(data),
		)
	}
	return &AtapiMessagingNode{
		Controller:        AtaController(data[0]),
		Drive:             AtaDrive(data[1]),
		LogicalUnitNumber: binary.LittleEndian.Uint16(data[2:4]),
	}, nil
}
