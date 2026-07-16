package efi

import (
	"fmt"
)

type ProtocolMediaNode struct {
	GUID GUID
}

func (p *ProtocolMediaNode) String() string {
	return fmt.Sprintf("Media(%x)", p.GUID)
}

func parseProtocolMediaNode(data []byte) (*ProtocolMediaNode, error) {
	if len(data) != 16 {
		return nil, fmt.Errorf(
			"invalid protocol media node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := ParseGUID(data)
	if err != nil {
		return nil, fmt.Errorf("parse protocol GUID: %w", err)
	}

	return &ProtocolMediaNode{GUID: guid}, nil
}
