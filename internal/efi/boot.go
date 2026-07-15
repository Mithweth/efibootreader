package efi

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

func GetBootOrder() ([]uint16, error) {
	return getUint16List("BootOrder")
}

func GetBootCurrent() (uint16, error) {
	return getUint16("BootCurrent")
}

func GetBootEntry(id uint16) (*BootEntry, error) {
	name := fmt.Sprintf("Boot%04X", id)

	variable, err := GetVariable(name)
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
		DevicePath:     devicePathNodes,
		OptionalData:   data[devicePathEnd:],
	}, nil
}
