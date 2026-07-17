package efi

import (
	"fmt"
	"io"
	"encoding/binary"
)

type AtaController uint8
type AtaDrive uint8


const (
    AtaPrimary AtaController = 0
    AtaSecondary AtaController = 1
    AtaMaster AtaDrive = 0
    AtaSlave  AtaDrive = 1
)

type AtapiMessagingNode struct {
    Controller AtaController
    Drive      AtaDrive
    LogicalUnitNumber uint16
}

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

func (h *AtapiMessagingNode) String() string {
	return fmt.Sprintf("Ata(%d,%d,%d)", h.Controller, h.Drive, h.LogicalUnitNumber)
}

func (h *AtapiMessagingNode) GoString() string {
	if h == nil {
		return "(*efi.AtapiMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&efi.AtapiMessagingNode{"+
			"Controller:%#v, "+
			"Drive:%#v, "+
			"LogicalUnitNumber:%#v}",
		h.Controller,
		h.Drive,
		h.LogicalUnitNumber,
	)
}

func (h *AtapiMessagingNode) dump(w io.Writer, indent string) {

    fmt.Fprintf(w, "%sAtapi Messaging Node\n", indent)
    fmt.Fprintf(w, "%s  Controller\t : %d\n", indent, h.Controller)
    fmt.Fprintf(w, "%s  Drive\t\t : %d\n", indent, h.Drive)
    fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
}

func parseAtapiMessagingNode(data []byte) (*AtapiMessagingNode, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf(
			"invalid messaging atapi node payload size: got %d, want 4",
			len(data),
		)
	}
	return &AtapiMessagingNode{
		Controller: AtaController(data[0]),
		Drive: AtaDrive(data[1]),
		LogicalUnitNumber: binary.LittleEndian.Uint16(data[2:4]),
	}, nil
}
