package devicepath

import (
	"fmt"
	"io"
	"github.com/Mithweth/efibootreader/identifiers"
)

type FirewareFileMediaNode struct {
	GUID identifiers.GUID
}

func (p *FirewareFileMediaNode) String() string {
	return fmt.Sprintf("FvFile(%x)", p.GUID)
}

func (p *FirewareFileMediaNode) GoString() string {
	if p == nil {
		return "(*efi.FirewareFileMediaNode)(nil)"
	}

	return fmt.Sprintf("&efi.FirewareFileMediaNode{GUID:%#v}", p.GUID)
}

func (p *FirewareFileMediaNode) dump(w io.Writer, indent string) {
    fmt.Fprintf(w, "%sFireware File Media Node\n", indent)
    fmt.Fprintf(w, "%s  GUID\t : %s\n", indent, p.GUID)
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
