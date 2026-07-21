package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

type SasExMessagingDeviceInfo uint16

type SasExMessagingNode struct {
	Address            [8]byte
	LogicalUnitNumber  [8]byte
	DeviceInfo         SasExMessagingDeviceInfo
	RelativeTargetPort uint16
}

func (i SasExMessagingDeviceInfo) InformationLength() uint8 {
	return uint8(i & 0xf)
}

func (i SasExMessagingDeviceInfo) DeviceType() uint8 {
	return uint8((i >> 4) & 0x3)
}

func (i SasExMessagingDeviceInfo) Topology() uint8 {
	return uint8((i >> 6) & 0x3)
}

func (i SasExMessagingDeviceInfo) DriveBay() uint16 {
	return uint16(i>>8) + 1
}

func (i SasExMessagingDeviceInfo) SasSataString() string {
	device := i.DeviceType()
	if device == 0 || device == 2 {
		return "SAS"
	} else {
		return "SATA"
	}
}

func (i SasExMessagingDeviceInfo) IsInternal() bool {
	return i.DeviceType() == 0 || i.DeviceType() == 1
}

func (i SasExMessagingDeviceInfo) LocationString() string {
	if i.IsInternal() {
		return "Internal"
	} else {
		return "External"
	}
}

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

func (i SasExMessagingDeviceInfo) GoString() string {
	return fmt.Sprintf("devicepath.SasExMessagingDeviceInfo(%#x)", uint16(i))
}

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

func (v *SasExMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sSAS Ex Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Address\t\t : %#x\n", indent, v.Address)
	_, _ = fmt.Fprintf(w, "%s  Logical Unit Number\t : %#x\n", indent, v.LogicalUnitNumber)
	_, _ = fmt.Fprintf(w, "%s  Relative Target Port\t : %d\n", indent, v.RelativeTargetPort)
	v.DeviceInfo.dump(w, indent+"  ")
}

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
