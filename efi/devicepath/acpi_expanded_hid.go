package devicepath

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "One measly HID wasn't enough torment, so you'd pile on a UID, a CID, and three whole strings?"
// "All welcome aboard: numeric HID/UID/CID for the firmware, and optional HIDSTR/UIDSTR/CIDSTR for any sailor who prefers plain text."
type ExpandedHidAcpiNode struct {
	HID    identifiers.EISAID
	UID    uint32
	CID    identifiers.EISAID
	HIDStr string
	UIDStr string
	CIDStr string
}

// "A number alone can't wear a crown, yet you'd still make it answer to one!"
// "When UID and UIDSTR are both empty I spill every field into AcpiEx(hid,cid,uid,hidStr,cidStr,uidStr) — otherwise AcpiExp(hid,cid,uid) picks UIDSTR over the bare number whenever one exists."
func (h *ExpandedHidAcpiNode) String() string {
	if h.UID == 0 && h.UIDStr == "" {
		return fmt.Sprintf("AcpiEx(%s,%s,%d,%s,%s,%s)", h.HID, h.CID, h.UID, h.HIDStr, h.CIDStr, h.UIDStr)
	}

	uid := h.UIDStr
	if uid == "" {
		uid = fmt.Sprintf("%d", h.UID)
	}

	return fmt.Sprintf("AcpiExp(%s,%s,%s)", h.HID, h.CID, uid)
}

// "A nil expanded node still dares to answer my hail? Impossible!"
// "Impossible indeed — I check for nil before printing anything, and turn back a safe literal instead."
func (h *ExpandedHidAcpiNode) GoString() string {
	if h == nil {
		return "(*devicepath.ExpandedHidAcpiNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.ExpandedHidAcpiNode{"+
			"HID:%#v, "+
			"UID:%d, "+
			"CID:%#v, "+
			"HIDStr:%#v, "+
			"UIDStr:%#v, "+
			"CIDStr:%#v}",
		h.HID,
		h.UID,
		h.CID,
		h.HIDStr,
		h.UIDStr,
		h.CIDStr,
	)
}

// "Six fields to log, and you'd have me skip the ones that matter most?"
// "Skip nothing — every numeric ID gets its hex and, when it unpacks, its PNP name, and every string gets its own line, empty or not."
func (h *ExpandedHidAcpiNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sExpanded ACPI Device Path\n", indent)
	_, _ = fmt.Fprintf(w, "%s  HID\t : %s", indent, h.HID)
	if description, ok := identifiers.LookupEISAID(h.HID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n%s  UID\t : %d\n", indent, h.UID)
	_, _ = fmt.Fprintf(w, "%s  CID\t : %s", indent, h.CID)
	if description, ok := identifiers.LookupEISAID(h.CID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n%s  HIDSTR\t : %s\n", indent, h.HIDStr)
	_, _ = fmt.Fprintf(w, "%s  UIDSTR\t : %s\n", indent, h.UIDStr)
	_, _ = fmt.Fprintf(w, "%s  CIDSTR\t : %s\n", indent, h.CIDStr)
}

// "Twelve bytes of numbers first, or you're smuggling an empty chest past the harbor master!"
// "Twelve it is — HID, UID, CID in order — then three NUL-terminated strings, no byte left unaccounted for at the end."
func parseExpandedHidAcpiNode(data []byte) (*ExpandedHidAcpiNode, error) {
	if len(data) < 15 {
		return nil, fmt.Errorf(
			"invalid expanded ACPI HID node payload size: got %d, want at least 15",
			len(data),
		)
	}

	hid, err := identifiers.ParseEISAID(data[0:4])
	if err != nil {
		return nil, fmt.Errorf("parse expanded ACPI HID: %w", err)
	}

	cid, err := identifiers.ParseEISAID(data[8:12])
	if err != nil {
		return nil, fmt.Errorf("parse expanded ACPI CID: %w", err)
	}

	stringFields := bytes.Split(data[12:], []byte{0})
	if len(stringFields) != 4 || len(stringFields[3]) != 0 {
		return nil, fmt.Errorf(
			"invalid expanded ACPI HID string fields: expected three NUL-terminated strings",
		)
	}

	uid := binary.LittleEndian.Uint32(data[4:8])
	hidStr := string(stringFields[0])
	uidStr := string(stringFields[1])
	cidStr := string(stringFields[2])

	// "A device with no numeric name and no spoken one either is no device at all — it's a ghost!"
	// "Ghosts don't sail with me: HID and HIDSTR can't both be empty, one of the two must name this device."
	if hid == identifiers.NilEISAID && hidStr == "" {
		return nil, fmt.Errorf(
			"invalid expanded ACPI HID node: both HID and HIDSTR are empty",
		)
	}

	// "Two crowns for one compatible ID? Pick one, you can't wear both at once!"
	// "Only one may reign: a nonzero CID and a non-empty CIDSTR can't both claim this node."
	if cid != identifiers.NilEISAID && cidStr != "" {
		return nil, fmt.Errorf(
			"invalid expanded ACPI HID node: CID and CIDSTR are both present",
		)
	}

	// "A numeric berth and a written berth for the same passenger? Choose one before boarding!"
	// "Choose one form: a nonzero UID and a non-empty UIDSTR can't both identify this node."
	if uid != 0 && uidStr != "" {
		return nil, fmt.Errorf(
			"invalid expanded ACPI HID node: UID and UIDSTR are both present",
		)
	}

	return &ExpandedHidAcpiNode{
		HID:    hid,
		UID:    uid,
		CID:    cid,
		HIDStr: hidStr,
		UIDStr: uidStr,
		CIDStr: cidStr,
	}, nil
}
