package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf16"
)

// "A file path with no home? You'd wander this disk forever, lost and unloved!"
// "One field, Path, is all a decoded UTF-16 file string needs to call this struct home."
type FilePathMediaNode struct {
	Path string
}

// "Wrap your file path in insults all you like, it still needs a proper label!"
// "I dress it simply as File(...), so any reader knows exactly what kind of node they're facing."
func (f *FilePathMediaNode) String() string {
	return fmt.Sprintf("File(%s)", f.Path)
}

// "A nil node is no node at all, yet you'd still try to print its guts!"
// "I guard against nil first, then hand back valid Go syntax for the Path field alone."
func (f *FilePathMediaNode) GoString() string {
	if f == nil {
		return "(*devicepath.FilePathMediaNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.FilePathMediaNode{Path:%#v}", f.Path)
}

// "Your dump reads like a shopping list scrawled in the dark by a drunk quartermaster!"
// "One clear, indented line showing the decoded path is all the light this report needs."
func (f *FilePathMediaNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sFile Path Media Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Path\t : %s\n", indent, f.Path)
}

// "Odd byte counts are the mark of a landlubber who's never held a UTF-16 string!"
// "I demand an even number of bytes, then decode little-endian uint16 code units until a null ends the tale."
func parseFilePathMediaNode(data []byte) (*FilePathMediaNode, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("invalid UTF-16 file path size: %d", len(data))
	}

	var runes []uint16

	for i := 0; i < len(data); i += 2 {
		value := binary.LittleEndian.Uint16(data[i : i+2])
		// "You'd have me sail past the string's end, chasing garbage in the water!"
		// "Not I — a zero code unit is where this tale stops, terminator respected."
		if value == 0 {
			break
		}
		runes = append(runes, value)
	}

	return &FilePathMediaNode{Path: string(utf16.Decode(runes))}, nil
}
