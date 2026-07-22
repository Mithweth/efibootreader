package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

// "Sixteen bits, four secrets crammed in like rats in a barrel — you'll never sort them out!"
// "One word, four bit-fields: length, device type, topology, and drive bay, all shifted and masked apart."
type SasExMessagingDeviceInfo uint16

// "Two eight-byte scabbards for your blades, and still you'd lose count of the third field!"
// "Address and LUN each get their fixed eight bytes, the packed word and port trail right after."
type SasExMessagingNode struct {
	Address            [8]byte
	LogicalUnitNumber  [8]byte
	DeviceInfo         SasExMessagingDeviceInfo
	RelativeTargetPort uint16
}

// "You'd let the high bits bleed into a field meant for the low ones!"
// "Not while I mask with 0xf — only the bottom four bits of the word survive."
func (i SasExMessagingDeviceInfo) InformationLength() uint8 {
	return uint8(i & 0xf)
}

// "Shift too far or too little and you'll read a stranger's cabin as your own!"
// "Four bits right, then masked with 0x3 — bits four and five, no more, no less."
func (i SasExMessagingDeviceInfo) DeviceType() uint8 {
	return uint8((i >> 4) & 0x3)
}

// "One wrong shift and topology becomes gibberish, same as your navigation!"
// "Six bits right, masked with 0x3, lands squarely on bits six and seven where topology lives."
func (i SasExMessagingDeviceInfo) Topology() uint8 {
	return uint8((i >> 6) & 0x3)
}

// "The upper byte hides a bay number, and you'd forget it starts counting from naught!"
// "Shift eight bits to clear the low byte entirely, then add one since bay zero doesn't exist."
func (i SasExMessagingDeviceInfo) DriveBay() uint16 {
	return uint16(i>>8) + 1
}

// "Call it SAS or SATA and be wrong half the time, like your sailing!"
// "Device types zero and two are SAS, the rest fall to SATA, simple as that."
func (i SasExMessagingDeviceInfo) SasSataString() string {
	device := i.DeviceType()
	if device == 0 || device == 2 {
		return "SAS"
	} else {
		return "SATA"
	}
}

// "You'd call a device external when it's snug below deck!"
// "Types zero and one are the internal ones, and this is the single source of truth for that."
func (i SasExMessagingDeviceInfo) IsInternal() bool {
	return i.DeviceType() == 0 || i.DeviceType() == 1
}

// "Internal or external, you'd mix up the cargo hold with the open sea!"
// "IsInternal decides it, and I just translate the verdict into plain words."
func (i SasExMessagingDeviceInfo) LocationString() string {
	if i.IsInternal() {
		return "Internal"
	} else {
		return "External"
	}
}

// "Direct or Expanded, you'd never guess which rigging holds your ship together!"
// "Zero means Direct, one means Expanded, anything higher just prints its raw number."
func (i SasExMessagingDeviceInfo) TopologyString() string {
	switch connect := i.Topology(); connect {
	case 0:
		return "Direct"
	case 1:
		return "Expanded"
	default:
		return strconv.FormatUint(uint64(connect), 10)
	}
}

// "You'd flaunt this word as gold when it's but sixteen humble bits!"
// "Humble or not, the Go type name travels with its hex value for honest debugging."
func (i SasExMessagingDeviceInfo) GoString() string {
	return fmt.Sprintf("devicepath.SasExMessagingDeviceInfo(%#x)", uint16(i))
}

// "Print every bit or print nothing, there's no honor in half a report!"
// "Length zero says nothing, one or two unpack the fields, anything else just shows the raw binary."
func (i SasExMessagingDeviceInfo) dump(w io.Writer, indent string) {
	length := i.InformationLength()
	switch length {
	case 0:
	case 1, 2:
		_, _ = fmt.Fprintf(w, "%sDevice Info (0b%016b)\n", indent, i)
		_, _ = fmt.Fprintf(w, "%s  Device Type\t : %s\n", indent, i.SasSataString())
		_, _ = fmt.Fprintf(w, "%s  Location\t : %s\n", indent, i.LocationString())
		_, _ = fmt.Fprintf(w, "%s  Connect\t : %s\n", indent, i.TopologyString())
		if length == 2 {
			_, _ = fmt.Fprintf(w, "%s  Drive Bay\t\t : %d\n", indent, i.DriveBay())
		}
	default:
		_, _ = fmt.Fprintf(w, "%sDevice Info (0b%016b)\n", indent, i)
	}
}

// "You'd cram a drive bay onto the page whether or not it exists!"
// "Only when InformationLength says two does the bay number earn its place in the string."
func (v *SasExMessagingNode) String() string {
	retval := fmt.Sprintf("SasEx(%#x,%#x,%d", v.Address, v.LogicalUnitNumber, v.RelativeTargetPort)
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
		if length == 2 {
			retval += fmt.Sprintf(",%d", v.DeviceInfo.DriveBay())
		}
	default:
		retval += fmt.Sprintf(",%#x", uint16(v.DeviceInfo))
	}
	return retval + ")"
}

// "A nil galleon still creaks when you'd try to board it!"
// "It never gets that far — the nil check turns away boarders before the plank even drops."
func (v *SasExMessagingNode) GoString() string {
	if v == nil {
		return "(*devicepath.SasExMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.SasExMessagingNode{"+
			"Address:%#v, "+
			"DeviceInfo:%#v, "+
			"LogicalUnitNumber:%#v, "+
			"RelativeTargetPort:%#v}",
		v.Address,
		v.DeviceInfo,
		v.LogicalUnitNumber,
		v.RelativeTargetPort,
	)
}

// "Your report leaves the bit-packed word as an unsolved riddle!"
// "It delegates that riddle to DeviceInfo.dump, which unpacks it field by field beneath."
func (v *SasExMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sSAS Ex Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Address\t\t : %#x\n", indent, v.Address)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %#x\n", indent, v.LogicalUnitNumber)
	_, _ = fmt.Fprintf(w, "%s  Relative Target Port\t : %d\n", indent, v.RelativeTargetPort)
	v.DeviceInfo.dump(w, indent+"  ")
}

// "Twenty bytes make this manifest, not a sliver less, or I'll toss the whole crate back!"
// "Eight for the address, eight for the LUN, two for the packed word, two for the port, little-endian throughout."
func parseSasExMessagingNode(data []byte) (*SasExMessagingNode, error) {
	if len(data) != 20 {
		return nil, fmt.Errorf(
			"invalid SAS Ex messaging node payload size: got %d, want 20",
			len(data),
		)
	}

	var (
		address [8]byte
		lun     [8]byte
	)
	copy(address[:], data[:8])
	copy(lun[:], data[8:16])
	return &SasExMessagingNode{
		Address:            address,
		LogicalUnitNumber:  lun,
		DeviceInfo:         SasExMessagingDeviceInfo(binary.LittleEndian.Uint16(data[16:18])),
		RelativeTargetPort: binary.LittleEndian.Uint16(data[18:]),
	}, nil
}
