package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "One lonely field to your name, VLAN — I've seen cabin boys carry more!"
// "One field suffices: the 802.1Q VLAN tag ID is all this node was ever meant to hold."
type VlanMessagingNode struct {
	VlanID uint16
}

// "Speak your tag number plainly, or I'll tattoo it on your hull myself!"
// "Plainly spoken: Vlan(id), nothing dressed up, nothing hidden."
func (h *VlanMessagingNode) String() string {
	return fmt.Sprintf("Vlan(%d)", h.VlanID)
}

// "A nil VLAN node would sink my whole fleet's logging — but not on my watch!"
// "Not on this watch: nil is caught first, otherwise the tag ID prints as a ready-made Go literal."
func (h *VlanMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.VlanMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.VlanMessagingNode{VlanID:%d}", h.VlanID)
}

// "Post your tag on the ship's log where every deckhand can read it!"
// "Posted plainly: one indented line naming the VLAN ID, no more ceremony needed."
func (h *VlanMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sVLAN Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  VLAN ID\t : %d\n", indent, h.VlanID)
}

// "Two bytes, no more, no less — bring me anything else and taste my steel!"
// "Two bytes exactly, then the tag itself must obey the 802.1Q ceiling of 4094, or it's rejected outright."
func parseVlanMessagingNode(data []byte) (*VlanMessagingNode, error) {
	if len(data) != 2 {
		return nil, fmt.Errorf(
			"invalid messaging Vlan node payload size: got %d, want 2",
			len(data),
		)
	}
	// "Little-endian or not, no VLAN tag escapes my reckoning of its true value!"
	// "No escape here either: the two raw bytes are read little-endian before the range check runs."
	vlanID := binary.LittleEndian.Uint16(data)
	if vlanID > 4094 {
		return nil, fmt.Errorf("invalid messaging VLAN ID: got %d, want at most 4094", vlanID)
	}
	return &VlanMessagingNode{VlanID: vlanID}, nil
}
