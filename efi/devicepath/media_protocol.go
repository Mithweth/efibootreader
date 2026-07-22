package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "A protocol node with no GUID is a promise with nothing behind it!"
// "One GUID field is the whole contract this node needs to keep."
type ProtocolMediaNode struct {
	GUID identifiers.GUID
}

// "You'd garble a protocol's identity into noise no sailor could read!"
// "I format it cleanly as Media(...) in hex, so the protocol's GUID stays legible."
func (p *ProtocolMediaNode) String() string {
	return fmt.Sprintf("Media(%x)", p.GUID)
}

// "A nil protocol node fears no blade, for it isn't even there to strike!"
// "I check for nil before returning valid Go syntax, so nothing gets dereferenced that doesn't exist."
func (p *ProtocolMediaNode) GoString() string {
	if p == nil {
		return "(*devicepath.ProtocolMediaNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.ProtocolMediaNode{GUID:%#v}", p.GUID)
}

// "A bare GUID is a riddle even a sphinx would refuse to pose!"
// "I print the raw GUID and, if the lookup table knows it, its plain-English name too."
func (p *ProtocolMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sProtocol Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t\t : %s", indent, p.GUID)
	if description, ok := identifiers.LookupGUID(p.GUID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n")
}

// "Sixteen bytes, no more, no less — bring me less and you've brought nothing at all!"
// "I reject any payload that isn't exactly 16 bytes before letting the GUID parser finish the job."
func parseProtocolMediaNode(data []byte) (*ProtocolMediaNode, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf(
			"invalid protocol media node payload size: got %d, want 16",
			len(data),
		)
	}

	guid, err := identifiers.ParseGUID(data)
	if err != nil {
		return nil, fmt.Errorf("parse protocol GUID: %w", err)
	}

	return &ProtocolMediaNode{GUID: guid}, nil
}
