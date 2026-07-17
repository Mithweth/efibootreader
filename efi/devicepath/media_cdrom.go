package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

type CdromMediaNode struct {
	BootEntry           uint32
	PartitionBlockStart uint64
	PartitionBlockSize  uint64
}

func (c *CdromMediaNode) String() string {
	return fmt.Sprintf("CDROM(0x%x,0x%x,0x%x)", c.BootEntry, c.PartitionBlockStart, c.PartitionBlockSize)
}

func (c *CdromMediaNode) GoString() string {
	if c == nil {
		return "(*devicepath.CdromMediaNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.CdromMediaNode{"+
			"BootEntry:%#v, "+
			"PartitionBlockStart:%#v, "+
			"PartitionBlockSize:%#v}",
		c.BootEntry,
		c.PartitionBlockStart,
		c.PartitionBlockSize,
	)
}

func (c *CdromMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sCD-Rom Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Boot Entry\t\t : %d\n", indent, c.BootEntry)
	_, _ = fmt.Fprintf(w, "%s  Partition Start (Block)\t : %d\n", indent, c.PartitionBlockStart)
	_, _ = fmt.Fprintf(w, "%s  Partition Size (Block)\t : %d\n", indent, c.PartitionBlockSize)
	_, _ = fmt.Fprintf(w, "%s  Partition End (Block)\t : %d\n", indent, c.PartitionBlockStart+c.PartitionBlockSize)
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
