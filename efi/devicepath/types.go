package devicepath

import (
	"io"
	"fmt"
)

type DevicePathType uint8
type MediaDevicePathSubType uint8
type MessagingDevicePathSubType uint8
type EndDevicePathSubType uint8

const (
	DevicePathHardware  DevicePathType = 0x01
	DevicePathACPI      DevicePathType = 0x02
	DevicePathMessaging DevicePathType = 0x03
	DevicePathMedia     DevicePathType = 0x04
	DevicePathBBS       DevicePathType = 0x05
	DevicePathEnd       DevicePathType = 0x7f
)

const (
	MediaHardDrive             MediaDevicePathSubType = 0x01
	MediaCdrom                 MediaDevicePathSubType = 0x02
	MediaVendor                MediaDevicePathSubType = 0x03
	MediaFilePath              MediaDevicePathSubType = 0x04
	MediaProtocol              MediaDevicePathSubType = 0x05
	MediaFirewareFile          MediaDevicePathSubType = 0x06
	MediaFirewareVolume        MediaDevicePathSubType = 0x07
	MessagingAtapi         MessagingDevicePathSubType = 0x01
	MessagingScsi                 MessagingDevicePathSubType = 0x02
	MessagingFibreChannel                MessagingDevicePathSubType = 0x03
	MessagingIeee1394              MessagingDevicePathSubType = 0x04
	MessagingUsb              MessagingDevicePathSubType = 0x05
	MessagingI2O              MessagingDevicePathSubType = 0x06
	MessagingMacAddress                 MessagingDevicePathSubType = 0x0b
		MessagingIPv4                 MessagingDevicePathSubType = 0x0c
		MessagingIPv6                 MessagingDevicePathSubType = 0x0d
	MessagingUsbWwid                MessagingDevicePathSubType = 0x10
	MessagingLogicalUnit              MessagingDevicePathSubType = 0x11
	MessagingSata              MessagingDevicePathSubType = 0x12
	MessagingNvmeNamespace MessagingDevicePathSubType = 0x17
	EndThisInstanceSubType     EndDevicePathSubType   = 0x01
	EndEntireDevicePathSubType EndDevicePathSubType   = 0xff
)

type DevicePathNodeDetails interface {
	fmt.Stringer
	fmt.GoStringer
	dump(io.Writer, string)
}

type DevicePathNode struct {
	Type    DevicePathType
	SubType uint8
	Data    []byte
	Details DevicePathNodeDetails
}

type DevicePath struct {
	Instances []DevicePathInstance
}

type DevicePathInstance struct {
	Nodes []DevicePathNode
}

