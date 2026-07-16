package efi

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

type FilePathMediaNode struct {
	Path string
}

func (f *FilePathMediaNode) String() string {
	return fmt.Sprintf("File(%s)", f.Path)
}

func (f *FilePathMediaNode) GoString() string {
	if f == nil {
		return "(*efi.FilePathMediaNode)(nil)"
	}

	return fmt.Sprintf("&efi.FilePathMediaNode{Path:%#v}", f.Path)
}

func parseFilePathMediaNode(data []byte) (*FilePathMediaNode, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("invalid UTF-16 file path size: %d", len(data))
	}

	var runes []uint16

	for i := 0; i < len(data); i += 2 {
		value := binary.LittleEndian.Uint16(data[i : i+2])
		if value == 0 {
			break
		}
		runes = append(runes, value)
	}

	return &FilePathMediaNode{Path: string(utf16.Decode(runes))}, nil
}
