package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

type NvdimmNamespaceMessagingNode struct {
	UUID identifiers.GUID
}

func (f *NvdimmNamespaceMessagingNode) String() string {
	return fmt.Sprintf("NVDIMM(%s)", f.UUID)
}

func (f *NvdimmNamespaceMessagingNode) GoString() string {
	if f == nil {
		return "(*devicepath.NvdimmNamespaceMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.NvdimmNamespaceMessagingNode{UUID:%#v}", f.UUID)
}

func (f *NvdimmNamespaceMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sNVDIMM Namespace Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  UUID\t : %s\n", indent, f.UUID)
}

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
