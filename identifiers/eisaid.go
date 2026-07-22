package identifiers

import (
	_ "embed"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

//go:embed eisaids.json
var eisaidDatabaseData []byte

var eisaidDatabase map[string]string

type EISAID uint32

var NilEISAID EISAID

func (e EISAID) String() string {
	id := uint32(e)
	c1 := (id >> 10) & 0x1f
	c2 := (id >> 5) & 0x1f
	c3 := id & 0x1f
	if c1 == 0 || c1 > 26 || c2 == 0 || c2 > 26 || c3 == 0 || c3 > 26 {
		return ""
	}
	return fmt.Sprintf("%c%c%c%04X", '@'+c1, '@'+c2, '@'+c3, id>>16)
}

func (e EISAID) GoString() string {
	return fmt.Sprintf("identifiers.EISAID(0x%08x)", uint32(e))
}

func ParseEISAID(data []byte) (EISAID, error) {
	if len(data) != 4 {
		return NilEISAID, fmt.Errorf("expected 4 bytes, got %d", len(data))
	}

	return EISAID(binary.LittleEndian.Uint32(data)), nil
}

func init() {
	if err := json.Unmarshal(eisaidDatabaseData, &eisaidDatabase); err != nil {
		panic(fmt.Errorf("unable to load EISAID database: %w", err))
	}
}

func LookupEISAID(g EISAID) (string, bool) {
	description, ok := eisaidDatabase[g.String()]
	return description, ok
}
