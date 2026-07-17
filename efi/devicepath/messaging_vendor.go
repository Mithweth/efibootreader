package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

type VendorMessagingNode struct {
	GUID identifiers.GUID
	Data []byte
}

func (v *VendorMessagingNode) String() string {
	if len(v.Data) == 0 {
		return fmt.Sprintf("VenMsg(%s)", v.GUID)
	}

	return fmt.Sprintf("VenMsg(%s,%x)", v.GUID, v.Data)
}

func (v *VendorMessagingNode) GoString() string {
	if v == nil {
		return "(*devicepath.VendorMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.VendorMessagingNode{"+
			"GUID:%#v, "+
			"Data:%#v}",
		v.GUID,
		v.Data,
	)
}

func (v *VendorMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sVendor Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t : %s", indent, v.GUID)
	if description, ok := identifiers.LookupGUID(v.GUID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n%s  Data\t : %x\n", indent, v.Data)
}

func parseVendorMessagingNode(data []byte) (*VendorMessagingNode, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf(
			"invalid vendor messaging node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := identifiers.ParseGUID(data[0:16])
	if err != nil {
		return nil, fmt.Errorf("parse vendor GUID: %w", err)
	}

	vendorData := make([]byte, len(data)-16)
	copy(vendorData, data[16:])

	return &VendorMessagingNode{
		GUID: guid,
		Data: vendorData,
	}, nil
}
