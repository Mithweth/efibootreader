package identifiers

import (
	"encoding/binary"
	"fmt"
)

type EUI64 uint64

var NilEUI64 EUI64

func (e EUI64) String() string {
	eui := uint64(e)
	return fmt.Sprintf(
		"%02X-%02X-%02X-%02X-%02X-%02X-%02X-%02X",
		byte(eui>>56),
		byte(eui>>48),
		byte(eui>>40),
		byte(eui>>32),
		byte(eui>>24),
		byte(eui>>16),
		byte(eui>>8),
		byte(eui),
	)
}

func (e EUI64) GoString() string {
	return fmt.Sprintf("efi.EUI64(0x%016x)", uint64(e))
}

func ParseEUI64(data []byte) (EUI64, error) {
	if len(data) != 8 {

		return NilEUI64, fmt.Errorf("expected 8 bytes, got %d", len(data))
	}

	return EUI64(binary.BigEndian.Uint64(data)), nil
}
