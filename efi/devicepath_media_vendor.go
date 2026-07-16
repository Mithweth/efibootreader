package efi

import (
	"fmt"
	"github.com/google/uuid"
)

type VendorMediaNode struct {
	GUID uuid.UUID
	Data []byte
}

func (v *VendorMediaNode) String() string {
	if len(v.Data) == 0 {
		return fmt.Sprintf("VenMedia(%s)", v.GUID)
	}

	return fmt.Sprintf("VenMedia(%s,%x)", v.GUID, v.Data)
}

func parseVendorMediaNode(data []byte) (*VendorMediaNode, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf(
			"invalid vendor media node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := ParseEFIGUID(data[0:16])
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
