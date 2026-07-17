package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

type I2OMessagingNode struct {
	TargetID uint32
}

func (h *I2OMessagingNode) String() string {
	return fmt.Sprintf("I2O(%d)", h.TargetID)
}

func (h *I2OMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.I2OMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.I2OMessagingNode{"+
			"TargetID:%#v}",
		h.TargetID,
	)
}

func (h *I2OMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sI2O Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Target ID\t : %d\n", indent, h.TargetID)
}

func parseI2OMessagingNode(data []byte) (*I2OMessagingNode, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf(
			"invalid messaging I2O node payload size: got %d, want 4",
			len(data),
		)
	}

	return &I2OMessagingNode{
		TargetID: binary.LittleEndian.Uint32(data[0:4]),
	}, nil
}
