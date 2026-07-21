package devicepath

func parseMessagingDevicePathNode(node DevicePathNode) (DevicePathNodeDetails, error) {
	switch MessagingDevicePathSubType(node.SubType) {
	case MessagingAtapi:
		return parseAtapiMessagingNode(node.Data)

	case MessagingScsi:
		return parseScsiMessagingNode(node.Data)

	case MessagingFibreChannel:
		return parseFibreChannelMessagingNode(node.Data)

	case MessagingIeee1394:
		return parseIeee1394MessagingNode(node.Data)

	case MessagingUsb:
		return parseUsbMessagingNode(node.Data)

	case MessagingInfiniBand:
		return parseInfiniBandMessagingNode(node.Data)

	case MessagingI2O:
		return parseI2OMessagingNode(node.Data)

	case MessagingMacAddress:
		return parseMacAddressMessagingNode(node.Data)

	case MessagingVendor:
		return parseVendorMessagingNode(node.Data)

	case MessagingIPv4:
		return parseIPv4MessagingNode(node.Data)

	case MessagingIPv6:
		return parseIPv6MessagingNode(node.Data)

	case MessagingUart:
		return parseUartMessagingNode(node.Data)

	case MessagingUsbClass:
		return parseUsbClassMessagingNode(node.Data)

	case MessagingUsbWwid:
		return parseUsbWwidMessagingNode(node.Data)

	case MessagingLogicalUnit:
		return parseLogicalUnitMessagingNode(node.Data)

	case MessagingSata:
		return parseSataMessagingNode(node.Data)

	case MessagingVlan:
		return parseVlanMessagingNode(node.Data)

	case MessagingFibreChannelEx:
		return parseFibreChannelExMessagingNode(node.Data)

	case MessagingSasEx:
		return parseSasExMessagingNode(node.Data)

	case MessagingNvmeNamespace:
		return parseNvmeNamespaceMessagingNode(node.Data)

	case MessagingUri:
		return parseUriMessagingNode(node.Data)

	case MessagingUfs:
		return parseUfsMessagingNode(node.Data)

	case MessagingSd:
		return parseSdMessagingNode(node.Data)

	case MessagingBluetooth:
		return parseBluetoothMessagingNode(node.Data)

	case MessagingWifi:
		return parseWifiMessagingNode(node.Data)

	case MessagingEmmc:
		return parseEmmcMessagingNode(node.Data)

	case MessagingBluetoothLE:
		return parseBluetoothLEMessagingNode(node.Data)

	case MessagingDns:
		return parseDnsMessagingNode(node.Data)

	case MessagingNvdimmNamespace:
		return parseNvdimmNamespaceMessagingNode(node.Data)

	case MessagingRestService:
		return parseRestServiceMessagingNode(node.Data)

	case MessagingNvmeOfNamespace:
		return parseNvmeOfNamespaceMessagingNode(node.Data)

	default:
		return unknownDevicePathNode(node), nil
	}
}
