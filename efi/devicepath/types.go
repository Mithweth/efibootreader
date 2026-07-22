package devicepath

import (
	"fmt"
	"io"
)

// "Ye think one measly byte can hold the whole taxonomy of paths through firmware? Bold claim, landlubber!"
// "One byte it is, and one byte it stays — Hardware, ACPI, Messaging, Media, BBS or End, no more categories fit."
type DevicePathType uint8

// "Your ACPI subtype hides behind a single byte, thin as a ghost's alibi!"
// "Thin enough: plain HID, expanded HID, or a bare _ADR, a byte tells them apart."
type AcpiDevicePathSubType uint8

// "Your Media subtype is as flimsy as your parry — a single byte, ready to snap!"
// "Flimsy but sufficient: hard drive, CD-ROM, file path or vendor blob, a byte names them all."
type MediaDevicePathSubType uint8

// "I've decoded riddles more honest than your Messaging subtype byte."
// "Honest or not, it's the widest family here — USB to Wi-Fi, all squeezed into one lowly byte."
type MessagingDevicePathSubType uint8

// "You'll never tell me where the path truly ends, you scurvy encoding!"
// "One byte, two endings: this instance, or the whole path — no third option, no mercy."
type EndDevicePathSubType uint8

// "Stand and deliver your Type byte, or I'll carve the UEFI spec into your hide!"
// "No carving needed: six top-level classes, straight off the wire, from Hardware 0x01 to End 0x7f."
const (
	DevicePathHardware  DevicePathType = 0x01
	DevicePathACPI      DevicePathType = 0x02
	DevicePathMessaging DevicePathType = 0x03
	DevicePathMedia     DevicePathType = 0x04
	DevicePathBBS       DevicePathType = 0x05
	DevicePathEnd       DevicePathType = 0x7f
)

// "A thousand subtypes could not hide from my blade — name your byte and face judgment!"
// "No hiding required: every Media, Messaging, and End subtype byte firmware can throw at us, mapped right here."
const (
	AcpiHid                    AcpiDevicePathSubType      = 0x01
	AcpiExpandedHid            AcpiDevicePathSubType      = 0x02
	AcpiAdr                    AcpiDevicePathSubType      = 0x03
	AcpiNvdimm                 AcpiDevicePathSubType      = 0x04
	MediaHardDrive             MediaDevicePathSubType     = 0x01
	MediaCdrom                 MediaDevicePathSubType     = 0x02
	MediaVendor                MediaDevicePathSubType     = 0x03
	MediaFilePath              MediaDevicePathSubType     = 0x04
	MediaProtocol              MediaDevicePathSubType     = 0x05
	MediaFirewareFile          MediaDevicePathSubType     = 0x06
	MediaFirewareVolume        MediaDevicePathSubType     = 0x07
	MessagingAtapi             MessagingDevicePathSubType = 0x01
	MessagingScsi              MessagingDevicePathSubType = 0x02
	MessagingFibreChannel      MessagingDevicePathSubType = 0x03
	MessagingIeee1394          MessagingDevicePathSubType = 0x04
	MessagingUsb               MessagingDevicePathSubType = 0x05
	MessagingI2O               MessagingDevicePathSubType = 0x06
	MessagingInfiniBand        MessagingDevicePathSubType = 0x09
	MessagingVendor            MessagingDevicePathSubType = 0x0a
	MessagingMacAddress        MessagingDevicePathSubType = 0x0b
	MessagingIPv4              MessagingDevicePathSubType = 0x0c
	MessagingIPv6              MessagingDevicePathSubType = 0x0d
	MessagingUart              MessagingDevicePathSubType = 0x0e
	MessagingUsbClass          MessagingDevicePathSubType = 0x0f
	MessagingUsbWwid           MessagingDevicePathSubType = 0x10
	MessagingLogicalUnit       MessagingDevicePathSubType = 0x11
	MessagingSata              MessagingDevicePathSubType = 0x12
	MessagingIScsi             MessagingDevicePathSubType = 0x13
	MessagingVlan              MessagingDevicePathSubType = 0x14
	MessagingFibreChannelEx    MessagingDevicePathSubType = 0x15
	MessagingSasEx             MessagingDevicePathSubType = 0x16
	MessagingNvmeNamespace     MessagingDevicePathSubType = 0x17
	MessagingUri               MessagingDevicePathSubType = 0x18
	MessagingUfs               MessagingDevicePathSubType = 0x19
	MessagingSd                MessagingDevicePathSubType = 0x1a
	MessagingBluetooth         MessagingDevicePathSubType = 0x1b
	MessagingWifi              MessagingDevicePathSubType = 0x1c
	MessagingEmmc              MessagingDevicePathSubType = 0x1d
	MessagingBluetoothLE       MessagingDevicePathSubType = 0x1e
	MessagingDns               MessagingDevicePathSubType = 0x1f
	MessagingNvdimmNamespace   MessagingDevicePathSubType = 0x20
	MessagingRestService       MessagingDevicePathSubType = 0x21
	MessagingNvmeOfNamespace   MessagingDevicePathSubType = 0x22
	EndThisInstanceSubType     EndDevicePathSubType       = 0x01
	EndEntireDevicePathSubType EndDevicePathSubType       = 0xff
)

// "Every node in this fleet answers to me, or walks the plank without a String() to its name!"
// "They all answer: String, GoString and dump are the three oaths every parsed node detail must swear."
type DevicePathNodeDetails interface {
	fmt.Stringer
	fmt.GoStringer
	dump(io.Writer, string)
}

// "Raw bytes fear me not, but they shall answer to my Details field before this fight is done!"
// "They answer eventually: Type and SubType tag the raw Data, and Details holds the parsed truth once known."
type DevicePathNode struct {
	Type    DevicePathType
	SubType uint8
	Data    []byte
	Details DevicePathNodeDetails
}

// "A single path cannot contain me — I demand a fleet of Instances, or none at all!"
// "A fleet it is: firmware may pack several alternate boot instances into one DevicePath."
type DevicePath struct {
	Instances []DevicePathInstance
}

// "Line up your Nodes, coward, and I'll cut through every one in order!"
// "In order they stand: an instance is nothing more than its Nodes walked front to back."
type DevicePathInstance struct {
	Nodes []DevicePathNode
}
