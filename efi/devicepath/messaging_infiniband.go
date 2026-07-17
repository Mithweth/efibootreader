package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

type InfiniBandMessagingNode struct {
	ResourceFlags uint32
	GUID          [16]byte
	ServiceID     uint64
	TargetID      uint64
	DeviceID      uint64
}

func (h *InfiniBandMessagingNode) String() string {
	return fmt.Sprintf(
		"Infiniband(%x,%x,%d,%d,%d)",
		h.ResourceFlags,
		h.GUID,
		h.ServiceID,
		h.TargetID,
		h.DeviceID,
	)
}

func (h *InfiniBandMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.InfiniBandMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.InfiniBandMessagingNode{"+
			"ResourceFlags:%#v, "+
			"GUID:%#v, "+
			"ServiceID:%#v, "+
			"TargetID:%#v, "+
			"DeviceID:%#v}",
		h.ResourceFlags,
		h.GUID,
		h.ServiceID,
		h.TargetID,
		h.DeviceID,
	)
}

func (h *InfiniBandMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sInfiniBand Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Resource Flags\t : %x\n", indent, h.ResourceFlags)
	_, _ = fmt.Fprintf(w, "%s  Port GUID\t\t : %x\n", indent, h.GUID)
	_, _ = fmt.Fprintf(w, "%s  IOC GUID/Service ID\t : %d\n", indent, h.ServiceID)
	_, _ = fmt.Fprintf(w, "%s  Target Port ID\t : %d\n", indent, h.TargetID)
	_, _ = fmt.Fprintf(w, "%s  Service ID\t\t : %d\n", indent, h.DeviceID)
}

func parseInfiniBandMessagingNode(data []byte) (*InfiniBandMessagingNode, error) {
	if len(data) != 44 {
		return nil, fmt.Errorf(
			"invalid messaging InfiniBand node payload size: got %d, want 44",
			len(data),
		)
	}

	var guid [16]byte
	copy(guid[:], data[4:20])
	return &InfiniBandMessagingNode{
		ResourceFlags: binary.LittleEndian.Uint32(data[0:4]),
		GUID:          guid,
		ServiceID:     binary.LittleEndian.Uint64(data[20:28]),
		TargetID:      binary.LittleEndian.Uint64(data[28:36]),
		DeviceID:      binary.LittleEndian.Uint64(data[36:]),
	}, nil
}
