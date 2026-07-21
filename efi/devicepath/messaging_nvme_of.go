package devicepath

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"unicode/utf8"
)

type NvmeOfNamespaceIdentifierType uint8

const (
	NvmeOfNIDTypeEUI64 NvmeOfNamespaceIdentifierType = 0x01
	NvmeOfNIDTypeNGUID NvmeOfNamespaceIdentifierType = 0x02
	NvmeOfNIDTypeUUID  NvmeOfNamespaceIdentifierType = 0x03
)

func (r NvmeOfNamespaceIdentifierType) String() string {
	switch r {
	case NvmeOfNIDTypeEUI64:
		return "EUI64"
	case NvmeOfNIDTypeNGUID:
		return "NGUID"
	case NvmeOfNIDTypeUUID:
		return "UUID"
	default:
		return fmt.Sprintf("Reserved(%d)", uint8(r))
	}
}

func (r NvmeOfNamespaceIdentifierType) GoString() string {
	return fmt.Sprintf("devicepath.NvmeOfNamespaceIdentifierType(%#v)", uint8(r))
}

type NvmeOfNamespaceMessagingNode struct {
	NIDT         NvmeOfNamespaceIdentifierType
	NID          [16]byte
	SubsystemNQN string
}

func (h *NvmeOfNamespaceMessagingNode) NIDString() string {
	switch h.NIDT {
	case NvmeOfNIDTypeEUI64:
		return fmt.Sprintf("%x", h.NID[:8])

	case NvmeOfNIDTypeNGUID:
		return fmt.Sprintf("%x", h.NID)

	case NvmeOfNIDTypeUUID:
		return fmt.Sprintf("urn:uuid:%s", uuid.UUID(h.NID))

	default:
		return fmt.Sprintf("%x", h.NID)
	}
}

func (h *NvmeOfNamespaceMessagingNode) String() string {
	return fmt.Sprintf("NVMEoF(%s,%s)", h.SubsystemNQN, h.NIDString())
}

func (h *NvmeOfNamespaceMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.NvmeOfNamespaceMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.NvmeOfNamespaceMessagingNode{"+
			"NIDT:%#v, "+
			"NID:%#v, "+
			"SubsystemNQN:%#v}",
		h.NIDT,
		h.NID,
		h.SubsystemNQN,
	)
}

func (h *NvmeOfNamespaceMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sNVMe-oF Namespace Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  NIDT\t\t : %d (%#x)\n", indent, uint8(h.NIDT), uint8(h.NIDT))
	_, _ = fmt.Fprintf(w, "%s  NID\t\t\t : %s\n", indent, h.NIDString())
	_, _ = fmt.Fprintf(w, "%s  Subsystem NQN\t : %s\n", indent, h.SubsystemNQN)
}

func parseNvmeOfNamespaceMessagingNode(data []byte) (*NvmeOfNamespaceMessagingNode, error) {
	if len(data) < 18 {
		return nil, fmt.Errorf(
			"invalid messaging NVMe-oF namespace node payload size: got %d, want at least 18",
			len(data),
		)
	}

	end := -1

	for i := 17; i < len(data); i++ {
		if data[i] == 0 {
			end = i
			break
		}
	}

	if end == -1 {
		return nil, fmt.Errorf(
			"invalid messaging NVMe-oF namespace node: SubsystemNQN is not null-terminated")
	}

	if end-16 > 224 {
		return nil, fmt.Errorf(
			"invalid messaging NVMe-oF namespace node SubsystemNQN size: got %d, want at most 224",
			end-16,
		)
	}

	if !utf8.Valid(data[17:end]) {
		return nil, fmt.Errorf(
			"invalid messaging NVMe-oF namespace node: SubsystemNQN is not valid UTF-8")
	}

	var nid [16]byte
	copy(nid[:], data[1:17])

	return &NvmeOfNamespaceMessagingNode{
		NIDT:         NvmeOfNamespaceIdentifierType(data[0]),
		NID:          nid,
		SubsystemNQN: string(data[17:end]),
	}, nil
}
