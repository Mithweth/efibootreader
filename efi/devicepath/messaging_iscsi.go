package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
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
	switch digest := o.HeaderDigest(); digest {
	case 0:
		return "None"
	case 2:
		return "CRC32C"
	default:
		return fmt.Sprintf("Reserved(%d)", digest)
	}
}

func (o IScsiLoginOptions) DataDigestString() string {
	switch digest := o.DataDigest(); digest {
	case 0:
		return "None"
	case 2:
		return "CRC32C"
	default:
		return fmt.Sprintf("Reserved(%d)", digest)
	}
}

func (o IScsiLoginOptions) AuthenticationString() string {
	switch method := o.AuthenticationMethod(); method {
	case 0:
		if o.ChapType() == 0 {
			return "CHAP_BI"
		}
		return "CHAP_UNI"
	case 2:
		return "None"
	default:
		return fmt.Sprintf("Reserved(%d)", method)
	}
}

func (o IScsiLoginOptions) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sLogin Options (0b%016b)\n", indent, o)
	_, _ = fmt.Fprintf(w, "%s  Header Digest\t : %s\n", indent, o.HeaderDigestString())
	_, _ = fmt.Fprintf(w, "%s  Data Digest\t : %s\n", indent, o.DataDigestString())
	_, _ = fmt.Fprintf(w, "%s  Authentication\t : %s\n", indent, o.AuthenticationString())
}

func (o IScsiLoginOptions) GoString() string {
	return fmt.Sprintf("devicepath.IScsiLoginOptions(%#x)", uint16(o))
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

func (h *IScsiMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.IScsiMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.IScsiMessagingNode{"+
			"Protocol:%#v, "+
			"TargetName:%#v, "+
			"PortalGroup:%#v, "+
			"LogicalUnitNumber:%#v, "+
			"LoginOptions:%#v}",
		h.Protocol,
		h.TargetName,
		h.PortalGroup,
		h.LogicalUnitNumber,
		h.LoginOptions,
	)
}

func (h *IScsiMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%siSCSI Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Protocol\t\t : %s\n", indent, h.Protocol)
	_, _ = fmt.Fprintf(w, "%s  Target Name\t\t : %s\n", indent, h.TargetName)
	_, _ = fmt.Fprintf(w, "%s  Portal Group\t\t : %d\n", indent, h.PortalGroup)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %x\n", indent, h.LogicalUnitNumber)
	h.LoginOptions.dump(w, indent+"  ")
}

func parseIScsiMessagingNode(data []byte) (*IScsiMessagingNode, error) {
	if len(data) < 14 {
		return nil, fmt.Errorf("invalid messaging iSCSI node payload size: got %d, want at least 14", len(data))
	}

	end := len(data)
	for i := 14; i < len(data); i++ {
		if data[i] == 0 {
			end = i
			break
		}
	}

	var lun [8]byte
	copy(lun[:], data[4:12])
	return &IScsiMessagingNode{
		Protocol:          IScsiProtocol(binary.LittleEndian.Uint16(data[0:2])),
		LoginOptions:      IScsiLoginOptions(binary.LittleEndian.Uint16(data[2:4])),
		LogicalUnitNumber: lun,
		PortalGroup:       binary.LittleEndian.Uint16(data[12:14]),
		TargetName:        string(data[14:end]),
	}, nil
}
