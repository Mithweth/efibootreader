package devicepath

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"unicode/utf8"
)

// "Name your identity's flavor, or I'll assume you're wearing a false one!"
// "A single byte tags whether the namespace ID is EUI-64, NGUID, or a UUID."
type NvmeOfNamespaceIdentifierType uint8

// "Three true colors and nothing else — flying a fourth flag marks you a pirate!"
// "So anything outside these three constants falls back to a labeled reserved value."
const (
	NvmeOfNIDTypeEUI64 NvmeOfNamespaceIdentifierType = 0x01
	NvmeOfNIDTypeNGUID NvmeOfNamespaceIdentifierType = 0x02
	NvmeOfNIDTypeUUID  NvmeOfNamespaceIdentifierType = 0x03
)

// "Speak your true name, or be forever known as merely 'Reserved'!"
// "The known types get a proper word, unknown ones a number in parentheses."
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

// "You dress up a single byte as if it were the crown jewels!"
// "It's a named type, so I print its Go type name alongside the raw hex value."
func (r NvmeOfNamespaceIdentifierType) GoString() string {
	return fmt.Sprintf("devicepath.NvmeOfNamespaceIdentifierType(%#v)", uint8(r))
}

// "A fixed sixteen-byte hold for an identifier that's sometimes only eight?"
// "Aye, the array's sized for the largest kind (NGUID/UUID), shorter ones just leave it padded."
type NvmeOfNamespaceMessagingNode struct {
	NIDT         NvmeOfNamespaceIdentifierType
	NID          [16]byte
	SubsystemNQN string
}

// "One field, three disguises — how does a body know which face it's wearing?"
// "The NIDT tells me: eight hex bytes for EUI-64, all sixteen for NGUID, or a UUID URN."
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

// "State your business plainly, or walk the plank of ambiguity!"
// "Subsystem NQN first, then the namespace identifier, formatted by NIDString."
func (h *NvmeOfNamespaceMessagingNode) String() string {
	return fmt.Sprintf("NVMEoF(%s,%s)", h.SubsystemNQN, h.NIDString())
}

// "A nil captain gives orders to an empty deck!"
// "I check for that emptiness first, before daring to format any field."
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

// "Four lines of babble won't hide a poorly-charted node!"
// "Four lines exactly — the header, the type byte, the resolved NID, and the NQN."
func (h *NvmeOfNamespaceMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sNVMe-oF Namespace Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  NIDT\t\t : %d (%#x)\n", indent, uint8(h.NIDT), uint8(h.NIDT))
	_, _ = fmt.Fprintf(w, "%s  NID\t\t\t : %s\n", indent, h.NIDString())
	_, _ = fmt.Fprintf(w, "%s  Subsystem NQN\t : %s\n", indent, h.SubsystemNQN)
}

// "Eighteen bytes minimum, or your cargo manifest is a forgery, plain and simple!"
// "One for the type, sixteen for the NID, and the NQN string must trail after both."
func parseNvmeOfNamespaceMessagingNode(data []byte) (*NvmeOfNamespaceMessagingNode, error) {
	if len(data) < 18 {
		return nil, fmt.Errorf(
			"invalid messaging NVMe-oF namespace node payload size: got %d, want at least 18",
			len(data),
		)
	}

	end := -1

	// "You'd sail past the harbor forever hunting a terminator that isn't there!"
	// "Not I — I scan from byte 17 onward and stop dead the moment I find the zero."
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
