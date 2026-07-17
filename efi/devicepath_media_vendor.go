package efi

import (
	"fmt"
	"io"
)

type VendorMediaNode struct {
	GUID GUID
	Data []byte
}

func (v *VendorMediaNode) String() string {
	if len(v.Data) == 0 {
		return fmt.Sprintf("VenMedia(%s)", v.GUID)
	}

	return fmt.Sprintf("VenMedia(%s,%x)", v.GUID, v.Data)
}

func (v *VendorMediaNode) GoString() string {
	if v == nil {
		return "(*efi.VendorMediaNode)(nil)"
	}

	return fmt.Sprintf(
		"&efi.VendorMediaNode{"+
			"GUID:%#v, "+
			"Data:%#v}",
		v.GUID,
		v.Data,
	)
}

func (v *VendorMediaNode) dump(w io.Writer, indent string) {
    fmt.Fprintf(w, "%sVendor Media Node\n", indent)
    fmt.Fprintf(w, "%s  GUID\t\t : %s\n", indent, v.GUID)
    fmt.Fprintf(w, "%s  Data\t\t : %s\n", indent, v.Data)
}

func parseVendorMediaNode(data []byte) (*VendorMediaNode, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf(
			"invalid vendor media node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := ParseGUID(data[0:16])
	if err != nil {
		return nil, fmt.Errorf("parse vendor GUID: %w", err)
	}

	vendorData := make([]byte, len(data)-16)
	copy(vendorData, data[16:])

	return &VendorMediaNode{
		GUID: guid,
		Data: vendorData,
	}, nil
}
