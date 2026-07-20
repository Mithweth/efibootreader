package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
	"strconv"
)

var (
	sasDevicePathGUID   = identifiers.MustParseEFIGUID("d487ddb4-008b-11d9-afdc-001083ffca4d")
	uartFlowControlGUID = identifiers.MustParseEFIGUID("37499a9d-542f-4c89-a026-35da142094e4")
)

type VendorMessagingNodeInterface interface {
	String() string
	GoString() string
	dump(w io.Writer, indent string)
}

type VendorMessagingNode struct {
	GUID identifiers.GUID
	Data []byte
}

func (v *VendorMessagingNode) String() string {
	if len(v.Data) == 0 {
		return fmt.Sprintf("VenMsg(%s)", v.GUID)
	}

	return fmt.Sprintf("VenMsg(%s,%x)", v.GUID, v.Data)
}

func (v *VendorMessagingNode) GoString() string {
	if v == nil {
		return "(*devicepath.VendorMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.VendorMessagingNode{"+
			"GUID:%#v, "+
			"Data:%#v}",
		v.GUID,
		v.Data,
	)
}

func (v *VendorMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sVendor Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t : %s", indent, v.GUID)
	if description, ok := identifiers.LookupGUID(v.GUID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n%s  Data\t : %x\n", indent, v.Data)
}

func parseVendorMessagingNode(data []byte) (VendorMessagingNodeInterface, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf(
			"invalid vendor messaging node payload size: got %d, want at least 16",
			len(data),
		)
	}

	guid, err := identifiers.ParseGUID(data[:16])
	if err != nil {
		return nil, fmt.Errorf("parse vendor GUID: %w", err)
	}

	switch guid {
	case sasDevicePathGUID:
		return parseSasMessagingNode(data)
	case uartFlowControlGUID:
		return parseUartFlowControlMessagingNode(data)
	default:
		vendorData := make([]byte, len(data)-16)
		copy(vendorData, data[16:])

		return &VendorMessagingNode{
			GUID: guid,
			Data: vendorData,
		}, nil
	}
}

type SasMessagingDeviceInfo uint16

type SasMessagingNode struct {
	Reserved           uint32
	Address            uint64
	LogicalUnitNumber  uint64
	DeviceInfo         SasMessagingDeviceInfo
	RelativeTargetPort uint16
}

func (i SasMessagingDeviceInfo) InformationLength() uint8 {
	return uint8(i & 0xf)
}

func (i SasMessagingDeviceInfo) DeviceType() uint8 {
	return uint8((i >> 4) & 0x3)
}

func (i SasMessagingDeviceInfo) Topology() uint8 {
	return uint8((i >> 6) & 0x3)
}

func (i SasMessagingDeviceInfo) DriveBay() uint16 {
	return uint16(i>>8) + 1
}

func (i SasMessagingDeviceInfo) SasSataString() string {
	device := i.DeviceType()
	if device == 0 || device == 2 {
		return "SAS"
	} else {
		return "SATA"
	}
}

func (i SasMessagingDeviceInfo) IsInternal() bool {
	return i.DeviceType() == 0 || i.DeviceType() == 1
}

func (i SasMessagingDeviceInfo) LocationString() string {
	if i.IsInternal() {
		return "Internal"
	} else {
		return "External"
	}
}

func (i SasMessagingDeviceInfo) TopologyString() string {
	switch connect := i.Topology(); connect {
	case 0:
		return "Direct"
	case 1:
		return "Expanded"
	default:
		return strconv.FormatUint(uint64(connect), 10)
	}
}

func (i SasMessagingDeviceInfo) GoString() string {
	return fmt.Sprintf("devicepath.SasMessagingDeviceInfo(%#x)", uint16(i))
}

func (i SasMessagingDeviceInfo) dump(w io.Writer, indent string) {
	length := i.InformationLength()
	switch length {
	case 0:
	case 1, 2:
		_, _ = fmt.Fprintf(w, "%sDevice Info (0b%016b)\n", indent, i)
		_, _ = fmt.Fprintf(w, "%s  Device Type\t : %s\n", indent, i.SasSataString())
		_, _ = fmt.Fprintf(w, "%s  Location\t : %s\n", indent, i.LocationString())
		_, _ = fmt.Fprintf(w, "%s  Connect\t : %s\n", indent, i.TopologyString())
		if length == 2 && i.IsInternal() {
			_, _ = fmt.Fprintf(w, "%s  Drive Bay\t\t : %d\n", indent, i.DriveBay())
		}
	default:
		_, _ = fmt.Fprintf(w, "%sDevice Info (0b%016b)\n", indent, i)
	}
}

func (v *SasMessagingNode) String() string {
	retval := fmt.Sprintf("SAS(0x%x,0x%x,%d", v.Address, v.LogicalUnitNumber, v.RelativeTargetPort)
	length := v.DeviceInfo.InformationLength()
	switch length {
	case 0:

	case 1, 2:
		retval += fmt.Sprintf(
			",%s,%s,%s",
			v.DeviceInfo.SasSataString(),
			v.DeviceInfo.LocationString(),
			v.DeviceInfo.TopologyString(),
		)
		if v.DeviceInfo.IsInternal() && length == 2 {
			retval += fmt.Sprintf(",%d", v.DeviceInfo.DriveBay())
		}
		if v.Reserved != 0 {
			if !v.DeviceInfo.IsInternal() || length != 2 {
				// Since that node string is positional, we need to set
				// a DriveBay placeholder if it doesnt exist
				retval += ",0"
			}

			retval += fmt.Sprintf(",0x%x", v.Reserved)
		}
	default:
		retval += fmt.Sprintf(",0x%x", uint16(v.DeviceInfo))
	}
	return retval + ")"
}

func (v *SasMessagingNode) GoString() string {
	if v == nil {
		return "(*devicepath.SasMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.SasMessagingNode{"+
			"Reserved:%#v, "+
			"Address:%#v, "+
			"DeviceInfo:%#v, "+
			"LogicalUnitNumber:%#v, "+
			"RelativeTargetPort:%#v}",
		v.Reserved,
		v.Address,
		v.DeviceInfo,
		v.LogicalUnitNumber,
		v.RelativeTargetPort,
	)
}

func (v *SasMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sSAS Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Address\t\t : %#x\n", indent, v.Address)
	_, _ = fmt.Fprintf(w, "%s  Reserved\t\t : %#x\n", indent, v.Reserved)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %#x\n", indent, v.LogicalUnitNumber)
	_, _ = fmt.Fprintf(w, "%s  Relative Target Port\t : %d\n", indent, v.RelativeTargetPort)
	v.DeviceInfo.dump(w, indent+"  ")
}

func parseUartFlowControlMessagingNode(data []byte) (*SasMessagingNode, error) {
	return nil, nil
}

func parseSasMessagingNode(data []byte) (*SasMessagingNode, error) {
	return nil, nil
}
