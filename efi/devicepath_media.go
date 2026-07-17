package efi

func parseMediaDevicePathNode(node DevicePathNode) (DevicePathNodeDetails, error) {
	switch MediaDevicePathSubType(node.SubType) {
	case MediaHardDrive:
		return parseHardDriveMediaNode(node.Data)

	case MediaCdrom:
		return parseCdromMediaNode(node.Data)

	case MediaVendor:
		return parseVendorMediaNode(node.Data)

	case MediaFilePath:
		return parseFilePathMediaNode(node.Data)

	case MediaProtocol:
		return parseProtocolMediaNode(node.Data)

	case MediaFirewareFile:
		return parseFirewareFileMediaNode(node.Data)

	case MediaFirewareVolume:
		return parseFirewareVolumeMediaNode(node.Data)

	default:
		return unknownDevicePathNode(node), nil
	}
}
