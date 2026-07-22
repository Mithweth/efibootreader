package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "You call that a vendor node? I've seen better cargo manifests nailed to a plank!"
// "This one carries a GUID for a name tag and a free-form Data hold for whatever the vendor smuggled aboard."
type VendorHardwareNode struct {
	GUID identifiers.GUID
	Data []byte
}

// "Speak plainly, or I'll toss your words overboard!"
// "Plainly then: an empty cargo hold prints just the GUID, a full one adds the payload in hex."
func (v *VendorHardwareNode) String() string {
	if len(v.Data) == 0 {
		return fmt.Sprintf("VenHw(%s)", v.GUID)
	}

	return fmt.Sprintf("VenHw(%s,%x)", v.GUID, v.Data)
}

// "A ghost ship has more substance than your nil pointers!"
// "Which is why this lookout calls out the empty hold by name before daring to read its cargo."
func (v *VendorHardwareNode) GoString() string {
	if v == nil {
		return "(*devicepath.VendorHardwareNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.VendorHardwareNode{"+
			"GUID:%#v, "+
			"Data:%#v}",
		v.GUID,
		v.Data,
	)
}

// "Your logbook is as dull as a rusted cutlass!"
// "Mine names the GUID's known identity too, when the registry recognizes the mark."
func (v *VendorHardwareNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sVendor Hardware Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t : %s", indent, v.GUID)
	if description, ok := identifiers.LookupGUID(v.GUID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n%s  Data\t : %x\n", indent, v.Data)
}

// "Sixteen bytes or you don't sail with my GUID, landlubber!"
// "Then the crew below decks — whatever remains past those sixteen — is copied off as the vendor's own cargo."
func parseVendorHardwareNode(data []byte) (*VendorHardwareNode, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf(
			"invalid vendor hardware node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := identifiers.ParseGUID(data[0:16])
	if err != nil {
		return nil, fmt.Errorf("parse vendor GUID: %w", err)
	}

	vendorData := make([]byte, len(data)-16)
	copy(vendorData, data[16:])

	return &VendorHardwareNode{
		GUID: guid,
		Data: vendorData,
	}, nil
}
