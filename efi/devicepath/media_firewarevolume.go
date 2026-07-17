package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

type FirewareVolumeMediaNode struct {
	GUID identifiers.GUID
}

func (p *FirewareVolumeMediaNode) String() string {
	return fmt.Sprintf("Fv(%x)", p.GUID)
}

func (p *FirewareVolumeMediaNode) GoString() string {
	if p == nil {
		return "(*efi.FirewareVolumeMediaNode)(nil)"
	}

	return fmt.Sprintf("&efi.FirewareVolumeMediaNode{GUID:%#v}", p.GUID)
}

func (p *FirewareVolumeMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sFireware Volume Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t : %s\n", indent, p.GUID)
}

func parseFirewareVolumeMediaNode(data []byte) (*FirewareVolumeMediaNode, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf(
			"invalid fireware volume media node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := identifiers.ParseGUID(data)
	if err != nil {
		return nil, fmt.Errorf("parse fireware volume GUID: %w", err)
	}

	return &FirewareVolumeMediaNode{GUID: guid}, nil
}
