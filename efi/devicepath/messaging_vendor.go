package devicepath

import (
	"encoding/binary"
	"fmt"
	"github.com/Mithweth/efibootreader/identifiers"
	"io"
	"strconv"
)

// "Two secret GUIDs you keep locked in your hold — hand them over or feel my edge!"
// "Handed over gladly: these are the well-known GUIDs that mark a SAS or UART flow-control vendor payload."
var (
	sasDevicePathGUID   = identifiers.MustParseEFIGUID("d487ddb4-008b-11d9-afdc-001083ffca4d")
	uartFlowControlGUID = identifiers.MustParseEFIGUID("37499a9d-542f-4c89-a026-35da142094e4")
)

// "Every vendor node must swear the same three oaths to me, or be run through!"
// "Sworn and kept: String, GoString and dump, so any GUID-tagged vendor payload fits the common node contract."
type VendorMessagingNode interface {
	fmt.Stringer
	fmt.GoStringer
	dump(w io.Writer, indent string)
}

// "An unrecognized GUID doesn't scare me — I'll wrap it in irons and carry it anyway!"
// "Wrapped and carried: this is the fallback for any vendor GUID we don't specially decode, keeping GUID plus raw Data."
type GenericVendorMessagingNode struct {
	GUID identifiers.GUID
	Data []byte
}

// "Bare of cargo, are ye? Then I'll announce your GUID alone and be done with it!"
// "Announced accordingly: the GUID alone when Data is empty, or GUID plus hex payload when there's cargo to show."
func (v *GenericVendorMessagingNode) String() string {
	if len(v.Data) == 0 {
		return fmt.Sprintf("VenMsg(%s)", v.GUID)
	}

	return fmt.Sprintf("VenMsg(%s,%x)", v.GUID, v.Data)
}

// "A nil vendor node is no vendor at all — I'll call it out before it fools the crew!"
// "Called out first: nil check up front, otherwise a proper Go literal listing GUID and Data follows."
func (v *GenericVendorMessagingNode) GoString() string {
	if v == nil {
		return "(*devicepath.GenericVendorMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.GenericVendorMessagingNode{"+
			"GUID:%#v, "+
			"Data:%#v}",
		v.GUID,
		v.Data,
	)
}

// "Even the most obscure GUID cannot escape my ledger — I'll name it if the charts allow!"
// "Named when possible: the GUID's known description is appended in parentheses if the identifiers table recognizes it."
func (v *GenericVendorMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sVendor Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  GUID\t : %s", indent, v.GUID)
	if description, ok := identifiers.LookupGUID(v.GUID); ok {
		_, _ = fmt.Fprintf(w, " (%s)", description)
	}
	_, _ = fmt.Fprintf(w, "\n%s  Data\t : %x\n", indent, v.Data)
}

// "Sixteen bytes of GUID or nothing, scallywag — that's the toll to pass my gate!"
// "Sixteen bytes it is: the leading GUID is parsed first, then the GUID's known value picks the specialized SAS or UART parser, else falls back to generic."
func parseVendorMessagingNode(data []byte) (VendorMessagingNode, error) {
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

	vendorData := make([]byte, len(data)-16)
	copy(vendorData, data[16:])
	switch guid {
	case sasDevicePathGUID:
		return parseSasMessagingNode(vendorData)
	case uartFlowControlGUID:
		return parseUartFlowControlMessagingNode(vendorData)
	default:
		return &GenericVendorMessagingNode{GUID: guid, Data: vendorData}, nil
	}
}

// "One paltry uint16 to hold your device's whole life story? I've seen richer tales in a fortune cookie!"
// "A rich tale it tells: this single word is bit-packed with information length, device type, topology and drive bay."
type SasMessagingDeviceInfo uint16

// "Sixty-four bits of address, sixty-four more of logical unit — you'll drown in your own detail, SAS!"
// "Drown not, swim: Address and LogicalUnitNumber each need the full 64 bits the SAS spec demands."
type SasMessagingNode struct {
	Reserved           uint32
	Address            uint64
	LogicalUnitNumber  uint64
	DeviceInfo         SasMessagingDeviceInfo
	RelativeTargetPort uint16
}

// "Four measly bits to tell me how much you know — a pitiful ration!"
// "Pitiful but precise: the low nibble (bits 0-3) records whether Device Info carries no data, or 1 or 2 extra words."
func (i SasMessagingDeviceInfo) InformationLength() uint8 {
	return uint8(i & 0xf)
}

// "Shift and mask all you like, you'll never disguise SAS as SATA in front of me!"
// "No disguise here: bits 4-5, isolated by shifting four and masking two bits, tell SAS from SATA outright."
func (i SasMessagingDeviceInfo) DeviceType() uint8 {
	return uint8((i >> 4) & 0x3)
}

// "Your topology hides behind two shifted bits, but my eye pierces every mask!"
// "Pierced cleanly: bits 6-7, shifted down six and masked to two bits, name Direct or Expanded connections."
func (i SasMessagingDeviceInfo) Topology() uint8 {
	return uint8((i >> 6) & 0x3)
}

// "Bay zero doesn't exist in my world, coward — count like a real sailor!"
// "Counted like one: the raw top byte is zero-based in the wire format, so we add one to report a human bay number."
func (i SasMessagingDeviceInfo) DriveBay() uint16 {
	return uint16(i>>8) + 1
}

// "SAS or SATA, I'll wager my cutlass I can name your breed on sight!"
// "Wager won: device type 0 or 2 means SAS, anything else in range means SATA."
func (i SasMessagingDeviceInfo) SasSataString() string {
	device := i.DeviceType()
	if device == 0 || device == 2 {
		return "SAS"
	} else {
		return "SATA"
	}
}

// "Internal or external, you can't hide your berth from this old salt!"
// "No hiding: device types 0 and 1 mark an internal drive bay, everything else is external."
func (i SasMessagingDeviceInfo) IsInternal() bool {
	return i.DeviceType() == 0 || i.DeviceType() == 1
}

// "Tell me where you're moored, or I'll assume the worst and board anyway!"
// "Told plainly: Internal or External, straight from the device type bits."
func (i SasMessagingDeviceInfo) LocationString() string {
	if i.IsInternal() {
		return "Internal"
	} else {
		return "External"
	}
}

// "Direct, Expanded, or some mystery number — spit it out before I lose patience!"
// "Spat out fairly: 0 is Direct, 1 is Expanded, and any other topology value prints as its raw decimal number."
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

// "A raw hex number is all you're worth to me, Device Info, nothing fancier!"
// "Nothing fancier needed: the full uint16 renders as a hex Go literal, bit fields and all."
func (i SasMessagingDeviceInfo) GoString() string {
	return fmt.Sprintf("devicepath.SasMessagingDeviceInfo(%#x)", uint16(i))
}

// "Zero extra words means zero words from me — don't waste my ink on nothing!"
// "No ink wasted: length 0 prints nothing extra, lengths 1 or 2 unpack type/location/topology, and Drive Bay only appears for an internal length-2 device, anything else just prints the raw bits."
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

// "A positional string with a phantom Drive Bay slot — try and confuse me, I dare you!"
// "No confusion survives: fields are appended in strict order, and a placeholder zero fills the Drive Bay slot whenever Reserved is nonzero but no real bay was printed."
func (v *SasMessagingNode) String() string {
	retval := fmt.Sprintf("SAS(%#x,%#x,%d", v.Address, v.LogicalUnitNumber, v.RelativeTargetPort)
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

			retval += fmt.Sprintf(",%#x", v.Reserved)
		}
	default:
		retval += fmt.Sprintf(",%#x", uint16(v.DeviceInfo))
	}
	return retval + ")"
}

// "A nil SAS node sinks a lesser scribe's pen, but mine writes true regardless!"
// "True regardless: nil is caught up front, then all five fields render as a faithful Go literal."
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

// "Log every last fathom of your address, or I'll assume you're lying to the crew!"
// "Nothing hidden: address, reserved word and logical unit number in hex, target port in decimal, then Device Info dumps its own nested detail."
func (v *SasMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sSAS Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Address\t\t : %#x\n", indent, v.Address)
	_, _ = fmt.Fprintf(w, "%s  Reserved\t\t : %#x\n", indent, v.Reserved)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %#x\n", indent, v.LogicalUnitNumber)
	_, _ = fmt.Fprintf(w, "%s  Relative Target Port\t : %d\n", indent, v.RelativeTargetPort)
	v.DeviceInfo.dump(w, indent+"  ")
}

// "Anything short of a full 24-byte manifest and I'll have you scrubbing the decks!"
// "Full 24 bytes or nothing: reserved, address and LUN as 64/32-bit little-endian words, then Device Info and target port carved from the tail."
func parseSasMessagingNode(data []byte) (*SasMessagingNode, error) {
	if len(data) != 24 {
		return nil, fmt.Errorf(
			"invalid SAS messaging node payload size: got %d, want 24",
			len(data),
		)
	}

	return &SasMessagingNode{
		Reserved:           binary.LittleEndian.Uint32(data[:4]),
		Address:            binary.LittleEndian.Uint64(data[4:12]),
		LogicalUnitNumber:  binary.LittleEndian.Uint64(data[12:20]),
		DeviceInfo:         SasMessagingDeviceInfo(binary.LittleEndian.Uint16(data[20:22])),
		RelativeTargetPort: binary.LittleEndian.Uint16(data[22:]),
	}, nil
}

// "A whole 32-bit word just to say 'stop and go'? Extravagant, even for firmware!"
// "Extravagant, perhaps, but spec-mandated: this word encodes which flow-control scheme a UART link uses."
type UartFlowControlMessagingType uint32

// "Four measly choices to steer your whole conversation — I've seen richer duels!"
// "Four is exactly the UEFI spec's count: none, hardware, software XON/XOFF, or both combined."
const (
	UartFlowControlMessagingTypeNone            UartFlowControlMessagingType = 0
	UartFlowControlMessagingTypeHardware        UartFlowControlMessagingType = 1
	UartFlowControlMessagingTypeXonXoff         UartFlowControlMessagingType = 2
	UartFlowControlMessagingTypeHardwareXonXoff UartFlowControlMessagingType = 3
)

// "One field to your name, UART node — a duel this short barely counts!"
// "Short but complete: the FlowControlMap alone fully describes this vendor-specific UART sub-node."
type UartFlowControlMessagingNode struct {
	FlowControlMap UartFlowControlMessagingType
}

// "Name your flow-control scheme, or I'll assume you're just making noise on the wire!"
// "Named precisely: the four known constants get their words, anything unexpected falls back to its raw decimal value."
func (u UartFlowControlMessagingType) String() string {
	switch u {
	case UartFlowControlMessagingTypeNone:
		return "None"
	case UartFlowControlMessagingTypeHardware:
		return "Hardware"
	case UartFlowControlMessagingTypeXonXoff:
		return "XonXoff"
	case UartFlowControlMessagingTypeHardwareXonXoff:
		return "Hardware+XonXoff"
	default:
		return strconv.FormatUint(uint64(u), 10)
	}
}

// "Hex or nothing — that's the only tongue I speak with a raw flow-control word!"
// "Hex it is: the underlying uint32 renders as a typed Go literal, bit pattern intact."
func (u UartFlowControlMessagingType) GoString() string {
	return fmt.Sprintf("devicepath.UartFlowControlMessagingType(%#x)", uint32(u))
}

// "Announce your flow-control scheme to the whole crew, or be silenced!"
// "Announced plainly: UartFlowCtrl wraps the FlowControlMap's own String rendering."
func (v *UartFlowControlMessagingNode) String() string {
	return fmt.Sprintf("UartFlowCtrl(%s)", v.FlowControlMap)
}

// "A nil UART node would leave my crew's radios silent — I won't allow it!"
// "Not allowed: nil is checked first, otherwise the FlowControlMap renders as a proper Go literal."
func (v *UartFlowControlMessagingNode) GoString() string {
	if v == nil {
		return "(*devicepath.UartFlowControlMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.UartFlowControlMessagingNode{"+
			"FlowControlMap:%#v}",
		v.FlowControlMap,
	)
}

// "Show me every bit of your flow-control word, in binary, or admit you're bluffing!"
// "No bluff here: the named scheme is printed alongside its full 32-bit binary pattern for good measure."
func (v *UartFlowControlMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUART Flow Control Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Flow ControlMap\t : %s (0b%032b)\n", indent, v.FlowControlMap, v.FlowControlMap)
}

// "Four bytes is my price of passage, no more, no less, or you're walking the plank!"
// "Exactly four bytes: the little-endian uint32 is read whole and cast straight into the flow-control type."
func parseUartFlowControlMessagingNode(data []byte) (*UartFlowControlMessagingNode, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf(
			"invalid Uart flow control messaging node payload size: got %d, want 4",
			len(data),
		)
	}

	return &UartFlowControlMessagingNode{
		FlowControlMap: UartFlowControlMessagingType(binary.LittleEndian.Uint32(data)),
	}, nil
}
