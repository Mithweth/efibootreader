package efi

import (
	"fmt"
	"github.com/google/uuid"
)

type GUID uuid.UUID

var Nil GUID

func (g GUID) String() string {
	return uuid.UUID(g).String()
}

func (g GUID) GoString() string {
    return fmt.Sprintf("efi.GUID(%q)", g.String())
}

// GUID definition :
//
//	typedef struct {
//	    UINT32 Data1; little-endian
//	    UINT16 Data2; little-endian
//	    UINT16 Data3; little-endian
//	    UINT8  Data4[8]; big-endian
//	} EFI_GUID;
func ParseGUID(data []byte) (GUID, error) {
	if len(data) != 16 {

		return Nil, fmt.Errorf("expected 16 bytes, got %d", len(data))
	}

	return GUID{
		data[3], data[2], data[1], data[0],
		data[5], data[4],
		data[7], data[6],
		data[8], data[9], data[10], data[11],
		data[12], data[13], data[14], data[15],
	}, nil
}
