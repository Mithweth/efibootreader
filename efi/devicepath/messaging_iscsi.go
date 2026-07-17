package devicepath

import (
	//	"encoding/binary"
	"fmt"
	// "io"
)

type IScsiProtocol uint16
type IScsiLoginOptions uint16

const (
	IScsiProtocolTCP IScsiProtocol = 0
)

func (p IScsiProtocol) String() string {
	switch p {
	case IScsiProtocolTCP:
		return "TCP"
	default:
		return fmt.Sprintf("%#x", uint16(p))
	}
}

func (p IScsiProtocol) GoString() string {
	switch p {
	case IScsiProtocolTCP:
		return "devicepath.IScsiProtocolTCP"
	default:
		return fmt.Sprintf("devicepath.IScsiProtocol(%#x)", uint16(p))
	}
}

func (o IScsiLoginOptions) HeaderDigest() uint8 {
	return uint8(o & 0x3)
}

func (o IScsiLoginOptions) DataDigest() uint8 {
	return uint8((o >> 2) & 0x3)
}

func (o IScsiLoginOptions) AuthenticationMethod() uint8 {
	return uint8((o >> 10) & 0x3)
}

func (o IScsiLoginOptions) ChapType() uint8 {
	return uint8((o >> 12) & 0x1)
}

func (o IScsiLoginOptions) HeaderDigestString() string {
	switch o.HeaderDigest() {
	case 0:
		return "None"
	case 2:
		return "CRC32C"
	default:
		return fmt.Sprintf("Reserved(%d)", o.HeaderDigest())
	}
}

func (o IScsiLoginOptions) DataDigestString() string {
	switch o.DataDigest() {
	case 0:
		return "None"
	case 2:
		return "CRC32C"
	default:
		return fmt.Sprintf("Reserved(%d)", o.DataDigest())
	}
}

func (o IScsiLoginOptions) AuthenticationString() string {
	switch o.AuthenticationMethod() {
	case 0:
		if o.ChapType() == 0 {
			return "CHAP_BI"
		}
		return "CHAP_UNI"
	case 2:
		return "None"
	default:
		return fmt.Sprintf("Reserved(%d)", o.AuthenticationMethod())
	}
}

type IScsiMessagingNode struct {
	Protocol          IScsiProtocol
	LoginOptions      IScsiLoginOptions
	LogicalUnitNumber [8]byte
	PortalGroup       uint16
	TargetName        string
}

func (h *IScsiMessagingNode) String() string {
	return fmt.Sprintf(
		"iSCSI(%s,%d,%x,%s,%s,%s,%s)",
		h.TargetName,
		h.PortalGroup,
		h.LogicalUnitNumber,
		h.LoginOptions.HeaderDigestString(),
		h.LoginOptions.DataDigestString(),
		h.LoginOptions.AuthenticationString(),
		h.Protocol,
	)
}

// func (h *IScsiMessagingNode) GoString() string {
// 	if h == nil {
// 		return "(*devicepath.IScsiMessagingNode)(nil)"
// 	}

// 	return fmt.Sprintf(
// 		"&devicepath.IScsiMessagingNode{"+
// 			"TargetID:%#v, "+
// 			"LogicalUnitNumber:%#v}",
// 		h.TargetID,
// 		h.LogicalUnitNumber,
// 	)
// }

// func (h *IScsiMessagingNode) dump(w io.Writer, indent string) {
// 	_, _ = fmt.Fprintf(w, "%sSCSI Messaging Node\n", indent)
// 	_, _ = fmt.Fprintf(w, "%s  Target ID\t\t : %d\n", indent, h.TargetID)
// 	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %d\n", indent, h.LogicalUnitNumber)
// }

// func parseIScsiMessagingNode(data []byte) (*IScsiMessagingNode, error) {
// 	if len(data) < 19 {
// 		return nil, fmt.Errorf("invalid messaging ISCSI node payload size: got %d, want at least 19", len(data))
// 	}

// 	return &IScsiMessagingNode{
// 		TargetID:          binary.LittleEndian.Uint16(data[0:2]),
// 		LogicalUnitNumber: binary.LittleEndian.Uint16(data[2:4]),
// 	}, nil
// }
