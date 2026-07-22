package devicepath

// "An ACPI node in disguise? I'll unmask it before you've drawn your blade!"
// "I switch on the SubType — plain HID, expanded HID, bare _ADR, or NVDIMM — to its rightful parser."
func parseAcpiDevicePathNode(node DevicePathNode) (DevicePathNodeDetails, error) {
	switch AcpiDevicePathSubType(node.SubType) {
	case AcpiHid:
		return parseHidAcpiNode(node.Data)

	case AcpiExpandedHid:
		return parseExpandedHidAcpiNode(node.Data)

	case AcpiAdr:
		return parseAdrAcpiNode(node.Data)

	case AcpiNvdimm:
		return parseNvdimmAcpiNode(node.Data)

	default:
		return unknownDevicePathNode(node), nil
	}
}
