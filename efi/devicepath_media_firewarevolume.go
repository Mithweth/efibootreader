package efi

import (
	"fmt"
)

type FirewareVolumeMediaNode struct {
	GUID GUID
}

func (p *FirewareVolumeMediaNode) String() string {
	return fmt.Sprintf("Fv(%x)", p.GUID)
}

func parseFirewareVolumeMediaNode(data []byte) (*FirewareVolumeMediaNode, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf(
			"invalid fireware volume media node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := ParseGUID(data)
	if err != nil {
		return nil, fmt.Errorf("parse fireware volume GUID: %w", err)
	}

	return &FirewareVolumeMediaNode{GUID: guid}, nil
}
