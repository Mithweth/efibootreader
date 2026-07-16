package efi

import (
	"fmt"
)

type FirewareFileMediaNode struct {
	GUID GUID
}

func (p *FirewareFileMediaNode) String() string {
	return fmt.Sprintf("FvFile(%x)", p.GUID)
}

func parseFirewareFileMediaNode(data []byte) (*FirewareFileMediaNode, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf(
			"invalid fireware file media node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := ParseGUID(data)
	if err != nil {
		return nil, fmt.Errorf("parse fireware file GUID: %w", err)
	}

	return &FirewareFileMediaNode{GUID: guid}, nil
}
