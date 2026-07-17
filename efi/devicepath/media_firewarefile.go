package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

type FirewareFileMediaNode struct {
	GUID identifiers.GUID
}

func (p *FirewareFileMediaNode) String() string {
	return fmt.Sprintf("FvFile(%x)", p.GUID)
}

func (p *FirewareFileMediaNode) GoString() string {
	if p == nil {
		return "(*devicepath.FirewareFileMediaNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.FirewareFileMediaNode{GUID:%#v}", p.GUID)
}

func (p *FirewareFileMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sFireware File Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t : %s", indent, p.GUID)
	if description, ok := identifiers.LookupGUID(p.GUID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n")
}

func parseFirewareFileMediaNode(data []byte) (*FirewareFileMediaNode, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf(
			"invalid fireware file media node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := identifiers.ParseGUID(data)
	if err != nil {
		return nil, fmt.Errorf("parse fireware file GUID: %w", err)
	}

	return &FirewareFileMediaNode{GUID: guid}, nil
}
