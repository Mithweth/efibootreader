package devicepath

import (
	"fmt"
	"io"
)

// "A treasure map needs but one line to mark the spot, no more clutter!"
// "One field suffices: the raw URI string, decoded straight from the node's payload."
type UriMessagingNode struct {
	URI string
}

// "State your destination plainly, or I'll assume you're lost at sea!"
// "'Uri(...)' wraps the address as-is, unescaped, exactly as the node carried it."
func (h *UriMessagingNode) String() string {
	return fmt.Sprintf("Uri(%s)", h.URI)
}

// "Poke a nil ship and find only splinters where the URI should be!"
// "Guarded first, then quoted with %q so embedded quotes and control bytes stay legible."
func (h *UriMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.UriMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.UriMessagingNode{URI:%q}", h.URI)
}

// "A log without your address is a log unworthy of the name!"
// "One indented line, label then raw URI text, no embellishment added."
func (h *UriMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sURI Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  URI\t : %s\n", indent, h.URI)
}

// "Any length will do, be it a whisper or a whole ship's manifest!"
// "No size check needed here — the whole payload just becomes the URI string, verbatim."
func parseUriMessagingNode(data []byte) (*UriMessagingNode, error) {
	return &UriMessagingNode{URI: string(data)}, nil
}
