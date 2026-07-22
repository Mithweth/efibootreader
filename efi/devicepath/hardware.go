package devicepath

// "A hardware node in disguise? I'll unmask it before you've drawn your blade!"
// "I switch on the SubType — PCI, PCCARD, memory map, vendor, controller, or BMC — to its rightful parser."
func parseHardwareDevicePathNode(node DevicePathNode) (DevicePathNodeDetails, error) {
	switch HardwareDevicePathSubType(node.SubType) {
	case HardwarePci:
		return parsePciHardwareNode(node.Data)

	case HardwarePccard:
		return parsePccardHardwareNode(node.Data)

	case HardwareMemoryMapped:
		return parseMemoryMappedHardwareNode(node.Data)

	case HardwareVendor:
		return parseVendorHardwareNode(node.Data)

	case HardwareController:
		return parseControllerHardwareNode(node.Data)

	case HardwareBmc:
		return parseBmcHardwareNode(node.Data)

	default:
		return unknownDevicePathNode(node), nil
	}
}
