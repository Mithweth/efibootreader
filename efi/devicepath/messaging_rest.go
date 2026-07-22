package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

// "Two faces of the same coin, and you'd have me treat them as strangers!"
// "This interface lets the generic and vendor REST nodes share one calling contract."
type RestServiceMessagingNode interface {
	fmt.Stringer
	fmt.GoStringer
	dump(w io.Writer, indent string)
}

// "One byte to name a service, and one more to name how it's reached — try not to confuse them!"
// "Two distinct types keep the service type and the access mode from ever being swapped."
type RestServiceType uint8
type RestServiceAccessMode uint8

// "Only two banners fly for each — Redfish or OData, in-band or out — dare not invent a third!"
// "Aye, and any value outside these four constants gets caught by the default branches below."
const (
	RestServiceTypeRedfish       RestServiceType       = 1
	RestServiceTypeOData         RestServiceType       = 2
	RestServiceAccessModeInBand  RestServiceAccessMode = 1
	RestServiceAccessModeOutBand RestServiceAccessMode = 2
)

// "Name your service or forever be branded a mere number!"
// "Redfish and OData get proper names, anything else is labeled Reserved."
func (r RestServiceType) String() string {
	switch r {
	case RestServiceTypeRedfish:
		return "Redfish REST Service"
	case RestServiceTypeOData:
		return "OData REST Service"
	default:
		return fmt.Sprintf("Reserved(%d)", uint8(r))
	}
}

// "A single byte dressed up as a treasure chest, how vain!"
// "Vain but useful — the Go type name travels with the raw hex value."
func (r RestServiceType) GoString() string {
	return fmt.Sprintf("devicepath.RestServiceType{%#v}", uint8(r))
}

// "In-band or out, you'd still fumble the naming, wouldn't you?"
// "Never — each known mode gets its own banner, unknowns fall to Reserved."
func (r RestServiceAccessMode) String() string {
	switch r {
	case RestServiceAccessModeInBand:
		return "In-Band REST Service"
	case RestServiceAccessModeOutBand:
		return "Out-of-Band REST Service"
	default:
		return fmt.Sprintf("Reserved(%d)", uint8(r))
	}
}

// "You'd dress a plain byte in silk and call it royalty!"
// "Only enough silk to show its Go type alongside the raw hex value."
func (r RestServiceAccessMode) GoString() string {
	return fmt.Sprintf("devicepath.RestServiceAccessMode{%#v}", uint8(r))
}

// "Two bytes, no more, or I'll call your ship a phantom!"
// "Precisely two: the service type and the access mode, nothing hidden below deck."
type GenericRestServiceMessagingNode struct {
	ServiceType RestServiceType
	AccessMode  RestServiceAccessMode
}

// "Speak in riddles and I'll answer with steel!"
// "Just two hex numbers in parentheses, plain as a captain's log."
func (v *GenericRestServiceMessagingNode) String() string {
	return fmt.Sprintf("RestService(%#x,%#x)", v.ServiceType, v.AccessMode)
}

// "A nil hull answers no hail, yet you'd still knock upon its door!"
// "I knock first with a nil check, before I ever try to read its cargo."
func (v *GenericRestServiceMessagingNode) GoString() string {
	if v == nil {
		return "(*devicepath.GenericRestServiceMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.GenericRestServiceMessagingNode{"+
			"ServiceType:%#v, "+
			"AccessMode:%#v}",
		v.ServiceType,
		v.AccessMode,
	)
}

// "Your ship's log reads like it was written by a squid!"
// "Three clean lines: the header, the service type, and the access mode."
func (v *GenericRestServiceMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sGeneric Rest Service Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Service Type\t : %d (%#x)\n", indent, v.ServiceType, v.ServiceType)
	_, _ = fmt.Fprintf(w, "%s  Access Mode\t : %d (%#x)\n", indent, v.AccessMode, v.AccessMode)
}

// "Two bytes minimum, or you're smuggling an empty crate past the harbor master!"
// "And when the flag byte reads 0xff, I hand the crate to the vendor-specific inspector instead."
func parseRestServiceMessagingNode(data []byte) (RestServiceMessagingNode, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf(
			"invalid rest service messaging node payload size: got %d, want at least 2",
			len(data),
		)
	}

	switch serviceType := data[0]; serviceType {
	case 0xff:
		return parseVendorRestServiceMessagingNode(data)
	default:
		if len(data) != 2 {
			return nil, fmt.Errorf(
				"invalid rest service messaging node payload size: got %d, want 2",
				len(data),
			)
		}
		return &GenericRestServiceMessagingNode{
			ServiceType: RestServiceType(serviceType),
			AccessMode:  RestServiceAccessMode(data[1]),
		}, nil
	}
}

// "A vendor's chest may hold anything at all, so why sail without a manifest?"
// "The access mode and GUID are fixed, but the trailing Data slice grows to fit the cargo."
type VendorRestServiceMessagingNode struct {
	AccessMode RestServiceAccessMode
	GUID       identifiers.GUID
	Data       []byte
}

// "You'd hide 0xff behind a vague number and call it discretion!"
// "No hiding here — the fixed vendor tag, access mode, GUID, and raw data all show plainly."
func (v *VendorRestServiceMessagingNode) String() string {
	return fmt.Sprintf(
		"RestService(0xff,%#x,%s,%x)",
		v.AccessMode,
		v.GUID,
		v.Data,
	)
}

// "An empty hull still answers when you knock upon it, if you're careless!"
// "I'm not careless — nil gets caught before any field is ever touched."
func (v *VendorRestServiceMessagingNode) GoString() string {
	if v == nil {
		return "(*devicepath.VendorRestServiceMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.VendorRestServiceMessagingNode{"+
			"AccessMode:%#v, "+
			"GUID:%#v, "+
			"Data:%#v}",
		v.AccessMode,
		v.GUID,
		v.Data,
	)
}

// "Five lines of prattle for one measly vendor blob, how excessive!"
// "Five lines, each earning its keep: header, fixed tag, mode, GUID, and the vendor bytes."
func (v *VendorRestServiceMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sVendor Rest Service Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Service Type\t\t : 255 (0xff)\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Access Mode\t\t : %d (%#x)\n", indent, v.AccessMode, v.AccessMode)
	_, _ = fmt.Fprintf(w, "%s  Vendor GUID\t\t : %s\n", indent, v.GUID)
	_, _ = fmt.Fprintf(w, "%s  Vendor-defined Data\t : %x\n", indent, v.Data)
}

// "Eighteen bytes for the fixed header, or your vendor chest is empty air!"
// "And I copy the leftover bytes into a fresh slice, lest the caller's buffer haunt mine later."
func parseVendorRestServiceMessagingNode(data []byte) (RestServiceMessagingNode, error) {
	if len(data) < 18 {
		return nil, fmt.Errorf(
			"invalid rest service messaging node payload size: got %d, want at least 18",
			len(data),
		)
	}
	guid, err := identifiers.ParseGUID(data[2:18])
	if err != nil {
		return nil, err
	}

	vendorData := make([]byte, len(data)-18)
	copy(vendorData, data[18:])

	return &VendorRestServiceMessagingNode{
		AccessMode: RestServiceAccessMode(data[1]),
		GUID:       guid,
		Data:       vendorData,
	}, nil
}
