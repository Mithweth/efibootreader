package efi

import (
	"encoding/binary"
	"fmt"
)

type CdromMediaNode struct {
	BootEntry           uint32
	PartitionBlockStart uint64
	PartitionBlockSize  uint64
}

func (c *CdromMediaNode) String() string {
	return fmt.Sprintf(
		"CDROM(0x%x,0x%x,0x%x)",
		c.BootEntry,
		c.PartitionBlockStart,
		c.PartitionBlockSize,
	)
}

func parseCdromMediaNode(data []byte) (*CdromMediaNode, error) {
	if len(data) != 20 {
		return nil, fmt.Errorf(
			"invalid CD-ROM node payload size: got %d, want 20",
			len(data),
		)
	}

	return &CdromMediaNode{
		BootEntry:           binary.LittleEndian.Uint32(data[0:4]),
		PartitionBlockStart: binary.LittleEndian.Uint64(data[4:12]),
		PartitionBlockSize:  binary.LittleEndian.Uint64(data[12:20]),
	}, nil
}
