package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

type VlanMessagingNode struct {
	VlanID uint16
}

func (h *VlanMessagingNode) String() string {
	return fmt.Sprintf("Vlan(%d)", h.VlanID)
}

func (h *VlanMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.VlanMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.VlanMessagingNode{VlanID:%d}", h.VlanID)
}

func (h *VlanMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sVLAN Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  VLAN ID\t : %d\n", indent, h.VlanID)
}

func parseVlanMessagingNode(data []byte) (*VlanMessagingNode, error) {
	if len(data) != 2 {
		return nil, fmt.Errorf(
			"invalid messaging Vlan node payload size: got %d, want 2",
			len(data),
		)
	}
	vlanID := binary.LittleEndian.Uint16(data)
	if vlanID > 4094 {
		return nil, fmt.Errorf("invalid messaging VLAN ID: got %d, want at most 4094", vlanID)
	}
	return &VlanMessagingNode{VlanID: vlanID}, nil
}
