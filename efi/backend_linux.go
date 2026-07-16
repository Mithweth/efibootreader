//go:build linux

package efi

import (
	"os"
	"fmt"
	"io/fs"
	"encoding/binary"
	"path/filepath"
	"strings"
	"github.com/google/uuid"
)

type linuxBackend struct {
	path string
}

func newBackend() variableBackend {
	return &linuxBackend{
		path: "/sys/firmware/efi/efivars",
	}
}

func (be *linuxBackend) IsEfi() bool {
	_, err := os.Stat(be.path)
	return err == nil
}

func (be *linuxBackend) getVariablePath(name string) (string, error) {
	entries, err := os.ReadDir(be.path)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), name) {
			return filepath.Join(be.path, e.Name()), nil
		}
	}
	return "", fs.ErrExist
}

func (be *linuxBackend) GetVariable(name string) (*Variable, error) {
	path, err := be.getVariablePath(name)
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
