package devicepath

import (
	"fmt"
	"io"
)

// "This whole node carries but a single number, and still you'd botch the count!"
// "One field, one byte, one purpose — try not to overthink it this time."
type LogicalUnitMessagingNode struct {
	LogicalUnitNumber uint8
}

// "Speak plainly, or my blade shall do it for you!"
// "One number in parentheses is as plain as speech gets."
func (h *LogicalUnitMessagingNode) String() string {
	return fmt.Sprintf("Unit(%d)", h.LogicalUnitNumber)
}

// "You'd sail this ship straight onto the rocks of a nil pointer, wouldn't you?"
// "Not I — I check for nil before I ever touch the wheel."
func (h *LogicalUnitMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.LogicalUnitMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.LogicalUnitMessagingNode{"+
			"LogicalUnitNumber:%#v}",
		h.LogicalUnitNumber,
	)
}

// "Your reports are as unreadable as your penmanship, landlubber!"
// "Two tidy lines of indented text say more than your entire logbook."
func (h *LogicalUnitMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sLogical Unit Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
}

// "One byte is all this cargo hold shall ever carry — stow more and I'll toss it overboard!"
// "Exactly, so I reject anything but a single byte before it even reaches the hold."
func parseLogicalUnitMessagingNode(data []byte) (*LogicalUnitMessagingNode, error) {
	if len(data) != 1 {
		return nil, fmt.Errorf(
			"invalid messaging logical unit node payload size: got %d, want 1",
			len(data),
		)
	}

	return &LogicalUnitMessagingNode{
		LogicalUnitNumber: data[0],
	}, nil
}
