package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "One transport, one number — you call that a protocol register? I've seen richer treasure maps!"
// "Richer or not, the iSCSI spec only defines TCP today, so a uint16 code point is all we need."
type IScsiProtocol uint16

// "Sixteen bits crammed with digests, chap types, and secrets — a locked chest if ever I saw one!"
// "Locked indeed, and the key is bit-shifting: each field lives at its own offset inside this uint16 bitfield."
type IScsiLoginOptions uint16

// "TCP or nothing, you say? A rather short list of allies for a stormy sea!"
// "Short but accurate: the spec reserves value 0 for TCP, the only transport iSCSI over EFI currently defines."
const (
	IScsiProtocolTCP IScsiProtocol = 0
)

// "Name your protocol plainly, or I'll assume you're smuggling contraband bytes!"
// "Plainly named: the known TCP value gets a word, anything else prints as raw hex."
func (p IScsiProtocol) String() string {
	switch p {
	case IScsiProtocolTCP:
		return "TCP"
	default:
		return fmt.Sprintf("%#x", uint16(p))
	}
}

// "Show me this protocol in Go's own true form, or forever hold your tongue!"
// "True form given: the TCP constant prints its own name, all else falls back to a hex literal."
func (p IScsiProtocol) GoString() string {
	switch p {
	case IScsiProtocolTCP:
		return "devicepath.IScsiProtocolTCP"
	default:
		return fmt.Sprintf("devicepath.IScsiProtocol(%#x)", uint16(p))
	}
}

// "The lowest two bits hold your header digest, and I'll mask out the rest without mercy!"
// "Mercy withheld and unnecessary: a bitwise AND with 0x3 keeps only bits 0-1, exactly where the spec puts it."
func (o IScsiLoginOptions) HeaderDigest() uint8 {
	return uint8(o & 0x3)
}

// "Shift me right by two, then strike down all but the lowest pair of bits!"
// "Struck true: shifting by 2 slides the data-digest field into place, then the 0x3 mask isolates its two bits."
func (o IScsiLoginOptions) DataDigest() uint8 {
	return uint8((o >> 2) & 0x3)
}

// "Ten places you'd shift this word — a voyage most sailors would never survive!"
// "Survived easily: bits 10-11 hold the authentication method, so shift by 10 then mask two bits to reveal it."
func (o IScsiLoginOptions) AuthenticationMethod() uint8 {
	return uint8((o >> 10) & 0x3)
}

// "One lonely bit at position twelve — hardly worth drawing a blade for!"
// "One bit is plenty: shift by 12 and mask with 0x1 to read the single CHAP-type flag cleanly."
func (o IScsiLoginOptions) ChapType() uint8 {
	return uint8((o >> 12) & 0x1)
}

// "Give the header digest a name, or I'll assume you invented a fourth mystery value!"
// "Named honestly: only 0 and 2 are defined by the spec, everything else prints as a labeled reserved number."
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

// "The data digest hides behind the same two mystery values — name it or be named a coward!"
// "No cowardice here: identical mapping to the header digest, since the wire format reuses the same code space."
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

// "Bidirectional or unidirectional CHAP — pick wrong and your ship sails the wrong sea!"
// "No wrong turn: method 0 checks the ChapType bit to distinguish CHAP_BI from CHAP_UNI, method 2 means no auth at all."
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

// "Lay out this whole bitfield in binary for all to gawk at, plus its three decoded meanings!"
// "Gawk away: a %016b binary dump followed by the three human-readable strings this options word encodes."
func (o IScsiLoginOptions) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sLogin Options (0b%016b)\n", indent, o)
	_, _ = fmt.Fprintf(w, "%s  Header Digest\t : %s\n", indent, o.HeaderDigestString())
	_, _ = fmt.Fprintf(w, "%s  Data Digest\t : %s\n", indent, o.DataDigestString())
	_, _ = fmt.Fprintf(w, "%s  Authentication\t : %s\n", indent, o.AuthenticationString())
}

// "Render this bitfield in Go's own dialect, or I'll call your code illiterate!"
// "Literate and terse: the whole uint16 prints as one typed hex literal, bits and all."
func (o IScsiLoginOptions) GoString() string {
	return fmt.Sprintf("devicepath.IScsiLoginOptions(%#x)", uint16(o))
}

// "Five fields to name one iSCSI target — you'd need a whole ledger to track your fleet!"
// "A ledger it is: protocol, packed options, an 8-byte LUN, a portal group, and a variable-length target name."
type IScsiMessagingNode struct {
	Protocol          IScsiProtocol
	LoginOptions      IScsiLoginOptions
	LogicalUnitNumber [8]byte
	PortalGroup       uint16
	TargetName        string
}

// "Recite this whole target's tale in one breath, if your lungs can even hold that much air!"
// "One breath is plenty: name, portal group, LUN, and the three decoded login-option strings, plus protocol."
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

// "Dereference a nil target and you'll be swimming with the fishes, mark my words!"
// "No swimming today: the nil check surfaces first, well before any field is ever read."
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

// "Report every last detail of this target, and let the login options tell their own tale too!"
// "Every detail told: four top-level fields printed, then LoginOptions.dump delegates to its own nested, more-indented report."
func (h *IScsiMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%siSCSI Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Protocol\t\t : %s\n", indent, h.Protocol)
	_, _ = fmt.Fprintf(w, "%s  Target Name\t\t : %s\n", indent, h.TargetName)
	_, _ = fmt.Fprintf(w, "%s  Portal Group\t\t : %d\n", indent, h.PortalGroup)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %x\n", indent, h.LogicalUnitNumber)
	h.LoginOptions.dump(w, indent+"  ")
}

// "Bring me fewer than fourteen bytes and I'll have nowhere left to plant my fixed fields!"
// "Nowhere is right: protocol, options, LUN and portal group alone claim the first 14 bytes, so anything shorter is refused."
func parseIScsiMessagingNode(data []byte) (*IScsiMessagingNode, error) {
	if len(data) < 14 {
		return nil, fmt.Errorf("invalid messaging iSCSI node payload size: got %d, want at least 14", len(data))
	}

	// "The target's name has no fixed length — how dare you hide its true end from me!"
	// "Hidden nowhere: the name is NUL-terminated, so we scan from byte 14 until a zero byte marks its end."
	end := len(data)
	for i := 14; i < len(data); i++ {
		if data[i] == 0 {
			end = i
			break
		}
	}

	// "Eight bytes of Logical Unit Number, copied by hand like some common deckhand's chore!"
	// "A chore worth doing right: copying into a fixed [8]byte keeps the LUN safe from the source slice's later mutation."
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
