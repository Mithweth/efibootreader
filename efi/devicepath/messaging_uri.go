package devicepath

import (
	"fmt"
	"io"
)

type UriMessagingNode struct {
	URI string
}

func (h *UriMessagingNode) String() string {
	return fmt.Sprintf("Uri(%s)", h.URI)
}

func (h *UriMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.UriMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.UriMessagingNode{URI:%q}", h.URI)
}

func (h *UriMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sURI Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  URI\t : %s\n", indent, h.URI)
}

func parseUriMessagingNode(data []byte) (*UriMessagingNode, error) {
	return &UriMessagingNode{URI: string(data)}, nil
}
