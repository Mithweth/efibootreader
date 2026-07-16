package efi

import "fmt"

type DevicePathType uint8
type MediaDevicePathSubType uint8
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
	MediaCDROM                 MediaDevicePathSubType = 0x02
	MediaVendor                MediaDevicePathSubType = 0x03
	MediaFilePath              MediaDevicePathSubType = 0x04
	MediaProtocol              MediaDevicePathSubType = 0x05
	MediaFirewareFile          MediaDevicePathSubType = 0x06
	MediaFirewareVolume        MediaDevicePathSubType = 0x07
	EndThisInstanceSubType     EndDevicePathSubType   = 0x01
	EndEntireDevicePathSubType EndDevicePathSubType   = 0xff
)

type DevicePathNode struct {
	Type    DevicePathType
	SubType uint8
	Data    []byte
	Details fmt.Stringer
}

type DevicePath struct {
	Instances []DevicePathInstance
}

type DevicePathInstance struct {
	Nodes []DevicePathNode
}

type UnknownDevicePathNode struct {
	Type    DevicePathType
	SubType uint8
	Data    []byte
}

