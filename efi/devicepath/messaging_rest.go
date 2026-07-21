package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
)

type RestServiceMessagingNode interface {
	String() string
	GoString() string
	dump(w io.Writer, indent string)
}

type RestServiceType uint8
type RestServiceAccessMode uint8

const (
	RestServiceTypeRedfish       RestServiceType       = 1
	RestServiceTypeOData         RestServiceType       = 2
	RestServiceAccessModeInBand  RestServiceAccessMode = 1
	RestServiceAccessModeOutBand RestServiceAccessMode = 2
)

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

func (r RestServiceType) GoString() string {
	return fmt.Sprintf("devicepath.RestServiceType{%#v}", uint8(r))
}

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

func (r RestServiceAccessMode) GoString() string {
	return fmt.Sprintf("devicepath.RestServiceAccessMode{%#v}", uint8(r))
}

type GenericRestServiceMessagingNode struct {
	ServiceType RestServiceType
	AccessMode  RestServiceAccessMode
}

func (v *GenericRestServiceMessagingNode) String() string {
	return fmt.Sprintf("RestService(%#x,%#x)", v.ServiceType, v.AccessMode)
}

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

func (v *GenericRestServiceMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sGeneric Rest Service Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Service Type\t : %d (%#x)\n", indent, v.ServiceType, v.ServiceType)
	_, _ = fmt.Fprintf(w, "%s  Access Mode\t : %d (%#x)\n", indent, v.AccessMode, v.AccessMode)
}

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

type VendorRestServiceMessagingNode struct {
	AccessMode RestServiceAccessMode
	GUID       identifiers.GUID
	Data       []byte
}

func (v *VendorRestServiceMessagingNode) String() string {
	return fmt.Sprintf(
		"RestService(0xff,%#x,%s,%x)",
		v.AccessMode,
		v.GUID,
		v.Data,
	)
}

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

func (v *VendorRestServiceMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sVendor Rest Service Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Service Type\t\t : 255 (0xff)\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Access Mode\t\t : %d (%#x)\n", indent, v.AccessMode, v.AccessMode)
	_, _ = fmt.Fprintf(w, "%s  Vendor GUID\t\t : %s\n", indent, v.GUID)
	_, _ = fmt.Fprintf(w, "%s  Vendor-defined Data\t : %x\n", indent, v.Data)
}

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
