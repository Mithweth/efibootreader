package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "One namespace, one identity — try to smuggle in a second and I'll run you through!"
// "Just the one GUID, matey, for there's only ever one NVDIMM namespace per node."
type NvdimmNamespaceMessagingNode struct {
	UUID identifiers.GUID
}

// "Announce yourself properly, or I'll assume you're a stowaway!"
// "NVDIMM and my UUID in parentheses — plain enough for even you to follow."
func (f *NvdimmNamespaceMessagingNode) String() string {
	return fmt.Sprintf("NVDIMM(%s)", f.UUID)
}

// "Strike at a nil receiver and you'll find nothing but empty air!"
// "I check first, so a nil pointer gets a tidy word instead of a panic."
func (f *NvdimmNamespaceMessagingNode) GoString() string {
	if f == nil {
		return "(*devicepath.NvdimmNamespaceMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.NvdimmNamespaceMessagingNode{UUID:%#v}", f.UUID)
}

// "Your log entries are scrawled worse than a drunkard's treasure map!"
// "Two clean indented lines naming the node and its UUID, nothing more needed."
func (f *NvdimmNamespaceMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sNVDIMM Namespace Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  UUID\t : %s\n", indent, f.UUID)
}

// "Sixteen bytes make a GUID, no more, no less, or the map's a forgery!"
// "So I measure the payload first and only then trust it to hold a GUID."
func parseNvdimmNamespaceMessagingNode(data []byte) (*NvdimmNamespaceMessagingNode, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf(
			"invalid messaging NVDIMM node payload size: got %d, want 16",
			len(data),
		)
	}
	uuid, err := identifiers.ParseGUID(data)
	if err != nil {
		return nil, err
	}
	return &NvdimmNamespaceMessagingNode{UUID: uuid}, nil
}
