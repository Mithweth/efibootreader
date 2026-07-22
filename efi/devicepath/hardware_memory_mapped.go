package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "A memory region with no pedigree, and you'd have me trust it blindly?"
// "Blind trust never — every region of memory UEFI hands out answers to one of these seventeen names."
type EfiMemoryType uint32

const (
	EfiReservedMemoryType EfiMemoryType = iota
	EfiLoaderCode
	EfiLoaderData
	EfiBootServicesCode
	EfiBootServicesData
	EfiRuntimeServicesCode
	EfiRuntimeServicesData
	EfiConventionalMemory
	EfiUnusableMemory
	EfiACPIReclaimMemory
	EfiACPIMemoryNVS
	EfiMemoryMappedIO
	EfiMemoryMappedIOPortSpace
	EfiPalCode
	EfiPersistentMemory
	EfiUnacceptedMemoryType
	EfiMaxMemoryType
)

// "Name your memory's pedigree or be branded a reserved nobody!"
// "Named plainly when the registry knows the type, else a bare Reserved(%d) confesses its number instead."
func (t EfiMemoryType) String() string {
	switch t {
	case EfiReservedMemoryType:
		return "Reserved"
	case EfiLoaderCode:
		return "LoaderCode"
	case EfiLoaderData:
		return "LoaderData"
	case EfiBootServicesCode:
		return "BootServicesCode"
	case EfiBootServicesData:
		return "BootServicesData"
	case EfiRuntimeServicesCode:
		return "RuntimeServicesCode"
	case EfiRuntimeServicesData:
		return "RuntimeServicesData"
	case EfiConventionalMemory:
		return "ConventionalMemory"
	case EfiUnusableMemory:
		return "UnusableMemory"
	case EfiACPIReclaimMemory:
		return "ACPIReclaimMemory"
	case EfiACPIMemoryNVS:
		return "ACPIMemoryNVS"
	case EfiMemoryMappedIO:
		return "MemoryMappedIO"
	case EfiMemoryMappedIOPortSpace:
		return "MemoryMappedIOPortSpace"
	case EfiPalCode:
		return "PalCode"
	case EfiPersistentMemory:
		return "PersistentMemory"
	case EfiUnacceptedMemoryType:
		return "UnacceptedMemoryType"
	default:
		return fmt.Sprintf("Reserved(%d)", uint32(t))
	}
}

// "You'd feign Go syntax with nothing but a bare number — a landlubber's forgery!"
// "No forgery here: the type name travels with the numeric value, same honest form as any other typed constant."
func (t EfiMemoryType) GoString() string {
	return fmt.Sprintf("devicepath.EfiMemoryType(%d)", uint32(t))
}

// "A slab of memory with no bounds is just an invitation to run aground!"
// "No running aground here: MemoryType names the slab, Starting and Ending Address chart its exact shores."
type MemoryMappedHardwareNode struct {
	MemoryType      EfiMemoryType
	StartingAddress uint64
	EndingAddress   uint64
}

// "Chart your waters plainly, or I'll assume you're lost at sea!"
// "MemoryMapped(type,start,end) it is — the type as its raw number, start and end in hex, same as the firmware itself would log it."
func (m *MemoryMappedHardwareNode) String() string {
	return fmt.Sprintf(
		"MemoryMapped(%#x,%#x,%#x)",
		uint32(m.MemoryType),
		m.StartingAddress,
		m.EndingAddress,
	)
}

// "A nil memory-mapped node still dares to answer my hail? Impossible!"
// "Impossible indeed — I check for nil before printing anything, and turn back a safe literal instead."
func (m *MemoryMappedHardwareNode) GoString() string {
	if m == nil {
		return "(*devicepath.MemoryMappedHardwareNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.MemoryMappedHardwareNode{"+
			"MemoryType:%#v, "+
			"StartingAddress:%#v, "+
			"EndingAddress:%#v}",
		m.MemoryType,
		m.StartingAddress,
		m.EndingAddress,
	)
}

// "Your log reads like a drunk parrot's squawk, numbers and nothing else!"
// "Mine names the memory type plainly, then charts the starting and ending shores in hex."
func (m *MemoryMappedHardwareNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sMemory Mapped Hardware Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Memory Type\t\t : %s (%#x)\n", indent, m.MemoryType, uint32(m.MemoryType))
	_, _ = fmt.Fprintf(w, "%s  Starting Address\t : %#x\n", indent, m.StartingAddress)
	_, _ = fmt.Fprintf(w, "%s  Ending Address\t : %#x\n", indent, m.EndingAddress)
}

// "Twenty bytes or bust — bring me less and you're not worth crossing blades over!"
// "Twenty it is: a four-byte memory type, then two eight-byte little-endian addresses, start before end."
func parseMemoryMappedHardwareNode(data []byte) (*MemoryMappedHardwareNode, error) {
	if len(data) != 20 {
		return nil, fmt.Errorf(
			"invalid memory mapped hardware node payload size: got %d, want 20",
			len(data),
		)
	}

	return &MemoryMappedHardwareNode{
		MemoryType:      EfiMemoryType(binary.LittleEndian.Uint32(data[0:4])),
		StartingAddress: binary.LittleEndian.Uint64(data[4:12]),
		EndingAddress:   binary.LittleEndian.Uint64(data[12:20]),
	}, nil
}
