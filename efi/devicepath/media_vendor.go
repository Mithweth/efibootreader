package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "You call that a vendor node? I've seen better cargo manifests nailed to a plank!"
// "This one carries a GUID for a name tag and a free-form Data hold for whatever the vendor smuggled aboard."
type VendorMediaNode struct {
	GUID identifiers.GUID
	Data []byte
}

// "Speak plainly, or I'll toss your words overboard!"
// "Plainly then: an empty cargo hold prints just the GUID, a full one adds the payload in hex."
func (v *VendorMediaNode) String() string {
	if len(v.Data) == 0 {
		return fmt.Sprintf("VenMedia(%s)", v.GUID)
	}

	return fmt.Sprintf("VenMedia(%s,%x)", v.GUID, v.Data)
}

// "A ghost ship has more substance than your nil pointers!"
// "Which is why this lookout calls out the empty hold by name before daring to read its cargo."
func (v *VendorMediaNode) GoString() string {
	if v == nil {
		return "(*devicepath.VendorMediaNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.VendorMediaNode{"+
			"GUID:%#v, "+
			"Data:%#v}",
		v.GUID,
		v.Data,
	)
}

// "Your logbook is as dull as a rusted cutlass!"
// "Mine names the GUID's known identity too, when the registry recognizes the mark."
func (v *VendorMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sVendor Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t : %s", indent, v.GUID)
	if description, ok := identifiers.LookupGUID(v.GUID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n%s  Data\t : %s\n", indent, v.Data)
}

// "Sixteen bytes or you don't sail with my GUID, landlubber!"
// "Then the crew below decks — whatever remains past those sixteen — is copied off as the vendor's own cargo."
func parseVendorMediaNode(data []byte) (*VendorMediaNode, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf(
			"invalid vendor media node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := identifiers.ParseGUID(data[0:16])
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
