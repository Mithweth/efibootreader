package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

type ProtocolMediaNode struct {
	GUID identifiers.GUID
}

func (p *ProtocolMediaNode) String() string {
	return fmt.Sprintf("Media(%x)", p.GUID)
}

func (p *ProtocolMediaNode) GoString() string {
	if p == nil {
		return "(*devicepath.ProtocolMediaNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.ProtocolMediaNode{GUID:%#v}", p.GUID)
}

func (p *ProtocolMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sProtocol Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t\t : %s", indent, p.GUID)
	if description, ok := identifiers.LookupGUID(p.GUID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n")
}

func parseProtocolMediaNode(data []byte) (*ProtocolMediaNode, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf(
			"invalid protocol media node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := identifiers.ParseGUID(data)
	if err != nil {
		return nil, fmt.Errorf("parse protocol GUID: %w", err)
	}

	return &ProtocolMediaNode{GUID: guid}, nil
}
