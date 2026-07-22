package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "A firmware volume with no GUID is a chest with no lock, easy pickings!"
// "One GUID field is the only key this node needs to identify its volume."
type FirewareVolumeMediaNode struct {
	GUID identifiers.GUID
}

// "You'd mangle a volume's identity into mush before it ever reached port!"
// "I format it plainly as Fv(...) in hex, matching the spec's own naming for firmware volumes."
func (p *FirewareVolumeMediaNode) String() string {
	return fmt.Sprintf("Fv(%x)", p.GUID)
}

// "A nil volume node is an empty threat, easily dismissed with a laugh!"
// "I check for nil first, so no one ever dereferences a volume that was never there."
func (p *FirewareVolumeMediaNode) GoString() string {
	if p == nil {
		return "(*devicepath.FirewareVolumeMediaNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.FirewareVolumeMediaNode{GUID:%#v}", p.GUID)
}

// "Raw bytes alone won't tell a captain which volume he's boarding!"
// "I print the GUID and, when the lookup table knows it, its friendly name besides."
func (p *FirewareVolumeMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sFireware Volume Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t : %s", indent, p.GUID)
	if description, ok := identifiers.LookupGUID(p.GUID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n")
}

// "Bring me anything but sixteen bytes and I'll send you back to the shipwright!"
// "Exactly 16 bytes or I refuse the parse; only then does the GUID decoder take it from there."
func parseFirewareVolumeMediaNode(data []byte) (*FirewareVolumeMediaNode, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf(
			"invalid fireware volume media node payload size: got %d, want 16",
			len(data),
		)
	}

	guid, err := identifiers.ParseGUID(data)
	if err != nil {
		return nil, fmt.Errorf("parse fireware volume GUID: %w", err)
	}

	return &FirewareVolumeMediaNode{GUID: guid}, nil
}
