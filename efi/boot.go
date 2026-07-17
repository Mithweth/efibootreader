package efi

import (
	"encoding/binary"
	"fmt"
	"strings"
	"unicode/utf16"
)

type BootEntry struct {
	Attributes     uint32
	Description    string
	FilePathLength uint16
	DevicePath     DevicePath
	OptionalData   []byte
}

func (b *BootEntry) Dump() string {
    var s strings.Builder
    fmt.Fprintf(&s, "ID\t\t : %04X\n", b.Attributes)
    fmt.Fprintf(&s, "Description\t : %s\n", b.Description)
    fmt.Fprintf(&s, "Length\t\t : %d\n", b.FilePathLength)
    b.DevicePath.dump(&s, "")
    return s.String()
}

func IsEFI() bool {
	return backend.IsEfi()
}

func GetBootOrder() ([]uint16, error) {
	return getUint16List("BootOrder")
}

func GetBootCurrent() (uint16, error) {
	return getUint16("BootCurrent")
}

func GetBootEntry(id uint16) (*BootEntry, error) {
	return getBootEntry(fmt.Sprintf("Boot%04X", id))
}

func GetBootIds() ([]uint16, error) {
	return backend.GetBootIds()
}

func getUint16(name string) (uint16, error) {
	variable, err := backend.GetVariable(name)
	if err != nil {
		return 0, err
	}

	if len(variable.Data) != 2 {
		return 0, fmt.Errorf("EFI variable %s: expected 2 bytes, got %d", name, len(variable.Data))
	}

	return binary.LittleEndian.Uint16(variable.Data), nil
}

func getUint16List(name string) ([]uint16, error) {
	variable, err := backend.GetVariable(name)
	if err != nil {
		return nil, err
	}

	return convertDataToUint16(variable.Data)
}

func getBootEntry(name string) (*BootEntry, error) {
	variable, err := backend.GetVariable(name)
	if err != nil {
		return nil, err
	}
	return ParseBootEntry(variable.Data)
}

func ParseBootEntry(data []byte) (*BootEntry, error) {
	runes := []uint16{}
	attrs := binary.LittleEndian.Uint32(data[:4])
	length := binary.LittleEndian.Uint16(data[4:6])
	i := 6
	for {
		if len(data) < i+2 {
			return nil, fmt.Errorf("unterminated UTF-16 description")
		}
		r := binary.LittleEndian.Uint16(data[i:])
		i += 2
		if r == 0 {
			break
		}
		runes = append(runes, r)
	}

	devicePathEnd := i + int(length)

	if devicePathEnd > len(data) {
		return nil, fmt.Errorf("device path exceeds data: start=%d length=%d data=%d", i, length, len(data))
	}
	desc := string(utf16.Decode(runes))
	devicePathNodes, err := ParseDevicePath(data[i:devicePathEnd])
	if err != nil {
		return nil, err
	}
	return &BootEntry{
		Attributes:     attrs,
		FilePathLength: length,
		Description:    desc,
		DevicePath:     *devicePathNodes,
		OptionalData:   data[devicePathEnd:],
	}, nil
}

func convertDataToUint16(data []byte) ([]uint16, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("invalid BootOrder size: %d bytes", len(data))
	}
	var retval []uint16
	for i := 0; i < len(data); i += 2 {
		retval = append(retval, binary.LittleEndian.Uint16(data[i:i+2]))
	}
	return retval, nil
}
