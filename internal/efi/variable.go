package efi

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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

func getUint16(name string) (uint16, error) {
	variable, err := GetVariable(name)
	if err != nil {
		return 0, err
	}

	if len(variable.Data) != 2 {
		return 0, fmt.Errorf("EFI variable %s: expected 2 bytes, got %d", name, len(variable.Data))
	}

	return binary.LittleEndian.Uint16(variable.Data), nil
}

func getUint16List(name string) ([]uint16, error) {
	variable, err := GetVariable(name)
	if err != nil {
		return nil, err
	}

	return convertDataToUint16(variable.Data)
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
