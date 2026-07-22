package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "A firmware volume file without a GUID is like a ship without a name!"
// "One GUID field is all this node carries, and it's enough to identify any file in the volume."
type FirewareFileMediaNode struct {
	GUID identifiers.GUID
}

// "You'd print a GUID as gibberish and call it seaworthy!"
// "I format it as FvFile(...) in hex, same convention the firmware spec itself demands."
func (p *FirewareFileMediaNode) String() string {
	return fmt.Sprintf("FvFile(%x)", p.GUID)
}

// "Strike at a nil pointer and you'll only wound yourself, fool!"
// "I check for nil before building the literal, so your blade never hits an empty hull."
func (p *FirewareFileMediaNode) GoString() string {
	if p == nil {
		return "(*devicepath.FirewareFileMediaNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.FirewareFileMediaNode{GUID:%#v}", p.GUID)
}

// "A bare GUID tells you nothing of what treasure it points to!"
// "I print the raw GUID, then consult the identifier lookup table to name the treasure if it's known."
func (p *FirewareFileMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sFireware File Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t : %s", indent, p.GUID)
	if description, ok := identifiers.LookupGUID(p.GUID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n")
}

// "Sixteen bytes is the price of admission, and you've brought less than that to the fight!"
// "I reject anything but exactly 16 bytes before letting the GUID parser take the wheel."
func parseFirewareFileMediaNode(data []byte) (*FirewareFileMediaNode, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf(
			"invalid fireware file media node payload size: got %d, want 16",
			len(data),
		)
	}

	guid, err := identifiers.ParseGUID(data)
	if err != nil {
		return nil, fmt.Errorf("parse fireware file GUID: %w", err)
	}

	return &FirewareFileMediaNode{GUID: guid}, nil
}
