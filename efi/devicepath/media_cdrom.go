package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "Your disc's boot entry and partition bounds mean nothing without a hull to hold them!"
// "Three fields — boot entry, block start, block size — carry every fact a CD-ROM node needs, no more, no less."
type CdromMediaNode struct {
	BootEntry           uint32
	PartitionBlockStart uint64
	PartitionBlockSize  uint64
}

// "Hex or decimal, you wouldn't know a boot entry from a barnacle!"
// "I print all three fields as 0x-prefixed hex, same style the firmware itself would recognize."
func (c *CdromMediaNode) String() string {
	return fmt.Sprintf("CDROM(0x%x,0x%x,0x%x)", c.BootEntry, c.PartitionBlockStart, c.PartitionBlockSize)
}

// "Nil pointers have sunk better sailors than you, and I won't be next!"
// "I check for nil first and return a safe literal, so no one drowns dereferencing an empty node."
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

// "You call that a report? I've seen better bookkeeping from a peg-legged parrot!"
// "Watch me tally boot entry, start block, and size block, computing the partition's end myself."
func (c *CdromMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sCD-Rom Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Boot Entry\t\t : %d\n", indent, c.BootEntry)
	_, _ = fmt.Fprintf(w, "%s  Partition Start (Block)\t : %d\n", indent, c.PartitionBlockStart)
	_, _ = fmt.Fprintf(w, "%s  Partition Size (Block)\t : %d\n", indent, c.PartitionBlockSize)
	_, _ = fmt.Fprintf(w, "%s  Partition End (Block)\t : %d\n", indent, c.PartitionBlockStart+c.PartitionBlockSize)
}

// "Twenty bytes or bust — bring me less and you're not worth crossing blades over!"
// "I demand exactly 20 bytes before decoding one uint32 and two little-endian uint64 block ranges."
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
