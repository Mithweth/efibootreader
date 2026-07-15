package efi

import (
	"os"
	"path/filepath"
	"fmt"
	"io/fs"
	"github.com/google/uuid"
	"strings"
	"encoding/binary"
	"unicode/utf16"
	"io/ioutil"
)

const EFIVARFS = "/sys/firmware/efi/efivars"

type Variable struct {
	Name       string
	GUID       uuid.UUID
	Attributes uint32
	Data       []byte
}

type BootEntry struct {
	Attributes uint32
	Description string
	FilePathLength uint16
	DevicePath   []DevicePathNode
	OptionalData []byte
}

type DevicePathNode struct {
	Type    uint8
	SubType uint8
	Data    []byte
}

type HardDriveNode struct {
	PartitionNumber uint32
	PartitionStart  uint64
	PartitionSize   uint64
	Signature       uuid.UUID
	MBRType         uint8
	SignatureType   uint8
}

func IsEFI() bool {
	_, err := os.Stat(EFIVARFS)
	return err == nil
}

func convertDataToUInt16(data []byte) ([]uint16, error) {
	if len(data) % 2 != 0 {
		return nil, fmt.Errorf("invalid BootOrder size: %d bytes", len(data))
	}
	var retval []uint16
	for i := 0; i < len(data); i += 2 {
		retval = append(retval, binary.LittleEndian.Uint16(data[i:i+2]))
	}
	return retval, nil
}

func getVariablePath(name string) (string, error) {
	entries, err := ioutil.ReadDir(EFIVARFS)
    if err != nil {
        return "", err
    }
    for _, e := range entries {
    	if strings.HasPrefix(e.Name(), name) {
            return filepath.Join(EFIVARFS, e.Name()), nil
    	}
    }
    return "", fs.ErrExist
}


// GUID definition :
// typedef struct {
//     UINT32 Data1; little-endian
//     UINT16 Data2; little-endian
//     UINT16 Data3; little-endian
//     UINT8  Data4[8]; big-endian
// } EFI_GUID;
func ParseEFIGUID(data []byte) (uuid.UUID, error) {
	if len(data) != 16 {
		return uuid.Nil, fmt.Errorf("expected 16 bytes, got %d", len(data))
	}

	return uuid.UUID{
		data[3], data[2], data[1], data[0],
		data[5], data[4],
		data[7], data[6],
		data[8], data[9], data[10], data[11],
		data[12], data[13], data[14], data[15],
	}, nil
}

func GetVariable(name string) (*Variable, error) {
	path, err := getVariablePath(name)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	attrs := binary.LittleEndian.Uint32(data[:4])
	realData := data[4:]
	splittedPath := strings.SplitN(path, "-", 2)
	if len(splittedPath) != 2 {
		return nil, fmt.Errorf("path is invalid: %s", path)
	}
	entryName := filepath.Base(splittedPath[0])
	guid, err := uuid.Parse(splittedPath[1])
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &Variable{Name: entryName, GUID: guid, Attributes: attrs, Data: realData}, nil
}

func ParseBootEntry(data []byte) (*BootEntry, error) {
	runes := []uint16{}
	attrs := binary.LittleEndian.Uint32(data[:4])
	length := binary.LittleEndian.Uint16(data[4:6])
	i:= 6
	for {
		if len(data) < i + 2 {
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
		Attributes: attrs,
		FilePathLength: length,
		Description: desc,
		DevicePath: devicePathNodes,
		OptionalData: data[devicePathEnd:],
	}, nil
}

func GetUint16(name string) (uint16, error) {
	variable, err := GetVariable(name)
	if err != nil {
		return 0, err
	}

	if len(variable.Data) != 2 {
		return 0, fmt.Errorf("EFI variable %s: expected 2 bytes, got %d", name, len(variable.Data))
	}

	return binary.LittleEndian.Uint16(variable.Data), nil
}

func GetUint16List(name string) ([]uint16, error) {
	variable, err := GetVariable(name)
	if err != nil {
		return nil, err
	}

	return convertDataToUInt16(variable.Data)
}

func ParseDevicePath(data []byte) ([]DevicePathNode, error) {
	var nodes []DevicePathNode

	offset := 0
	for {
		if offset + 4 > len(data) {
			return nil, fmt.Errorf("truncated device path header at offset %d", offset)
		}

		length := int(binary.LittleEndian.Uint16(data[offset + 2 : offset + 4]))

		if length < 4 {
			return nil, fmt.Errorf("invalid device path node length %d at offset %d", length, offset)
		}

		if offset + length > len(data) {
			return nil, fmt.Errorf("device path node exceeds buffer: offset=%d length=%d total=%d", offset, length, len(data))
		}

		nodes = append(nodes, DevicePathNode{
			Type:    data[offset],
			SubType: data[offset + 1],
			Data:    data[offset + 4:offset + length],
		})

		offset += length
		if offset >= len(data) {
			break
		}
	}

	return nodes, nil
}

func ParseHardDriveNode(data []byte) (*HardDriveNode, error) {
	if len(data) != 38 {
		return nil, fmt.Errorf("invalid hard drive node payload size: got %d, want 38", len(data))
	}
	sig, err := ParseEFIGUID(data[20:36])
	if err != nil {
		return nil, err
	}

	return &HardDriveNode{
		PartitionNumber: binary.LittleEndian.Uint32(data[0:4]),
		PartitionStart:  binary.LittleEndian.Uint64(data[4:12]),
		PartitionSize:   binary.LittleEndian.Uint64(data[12:20]),
		Signature:       sig,
		MBRType:         data[36],
		SignatureType:   data[37],
	}, nil
}

func ParseFilePathNode(data []byte) (string, error) {
	if len(data) % 2 != 0 {
		return "", fmt.Errorf("invalid UTF-16 file path size: %d", len(data))
	}

	var codeUnits []uint16

	for i := 0; i < len(data); i += 2 {
		value := binary.LittleEndian.Uint16(data[i:i + 2])
		if value == 0 {
			break
		}
		codeUnits = append(codeUnits, value)
	}

	return string(utf16.Decode(codeUnits)), nil
}

func GetBootOrder() ([]uint16, error) {
	return GetUint16List("BootOrder")
}

func GetBootCurrent() (uint16, error) {
	return GetUint16("BootCurrent")
}

func GetBootEntry(id uint16) (*BootEntry, error) {
	name := fmt.Sprintf("Boot%04X", id)

	variable, err := GetVariable(name)
	if err != nil {
		return nil, err
	}

	return ParseBootEntry(variable.Data)
}