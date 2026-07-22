package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

// "Name your parity or I'll assume you meant chaos!"
// "It's a plain uint8 underneath, but the named type keeps chaos from compiling."
type UartParity uint8

// "How many stops before the duel truly ends?"
// "Also a uint8 in disguise, one of four values: one, one-and-a-half, or two stop bits."
type UartStopBits uint8

// "Line up your values, one by one, and let no gap betray the firmware's intent!"
// "Each constant mirrors the UEFI spec's raw byte encoding, in numeric order, no surprises."
const (
	UartParityDefault      UartParity   = 0
	UartParityNone         UartParity   = 1
	UartParityEven         UartParity   = 2
	UartParityOdd          UartParity   = 3
	UartParityMark         UartParity   = 4
	UartParitySpace        UartParity   = 5
	UartStopBitsDefault    UartStopBits = 0
	UartStopBitsOne        UartStopBits = 1
	UartStopBitsOneAndHalf UartStopBits = 2
	UartStopBitsTwo        UartStopBits = 3
)

// "Four facts define a serial line, and I demand all four before I salute!"
// "Baud rate, data bits, parity, and stop bits — the whole conversation in one struct."
type UartMessagingNode struct {
	BaudRate uint64
	DataBits uint8
	Parity   UartParity
	StopBits UartStopBits
}

// "Call it by its true name, or hide behind cold hexadecimal forever!"
// "Known values get friendly words like 'Even' or 'Odd'; unknown ones fall back to %#x."
func (p UartParity) String() string {
	switch p {
	case UartParityDefault:
		return "Default"
	case UartParityNone:
		return "None"
	case UartParityEven:
		return "Even"
	case UartParityOdd:
		return "Odd"
	case UartParityMark:
		return "Mark"
	case UartParitySpace:
		return "Space"
	default:
		return fmt.Sprintf("%#x", uint8(p))
	}
}

// "Speak in code so a Go compiler itself could parse your boasts!"
// "Each known constant prints as its qualified identifier; strangers print as devicepath.UartParity(%#x)."
func (p UartParity) GoString() string {
	switch p {
	case UartParityDefault:
		return "devicepath.UartParityDefault"
	case UartParityNone:
		return "devicepath.UartParityNone"
	case UartParityEven:
		return "devicepath.UartParityEven"
	case UartParityOdd:
		return "devicepath.UartParityOdd"
	case UartParityMark:
		return "devicepath.UartParityMark"
	case UartParitySpace:
		return "devicepath.UartParitySpace"
	default:
		return fmt.Sprintf("devicepath.UartParity(%#x)", uint8(p))
	}
}

// "Halt your babbling and give the stop bits their proper title!"
// "Default, One, OneAndHalf, or Two — anything else just echoes back as raw hex."
func (s UartStopBits) String() string {
	switch s {
	case UartStopBitsDefault:
		return "Default"
	case UartStopBitsOne:
		return "One"
	case UartStopBitsOneAndHalf:
		return "OneAndHalf"
	case UartStopBitsTwo:
		return "Two"
	default:
		return fmt.Sprintf("%#x", uint8(s))
	}
}

// "Prove your stop bits are more than a costume of ones and zeroes!"
// "Known values become their devicepath.UartStopBitsX identifier; the rest wear a hex mask."
func (s UartStopBits) GoString() string {
	switch s {
	case UartStopBitsDefault:
		return "devicepath.UartStopBitsDefault"
	case UartStopBitsOne:
		return "devicepath.UartStopBitsOne"
	case UartStopBitsOneAndHalf:
		return "devicepath.UartStopBitsOneAndHalf"
	case UartStopBitsTwo:
		return "devicepath.UartStopBitsTwo"
	default:
		return fmt.Sprintf("devicepath.UartStopBits(%#x)", uint8(s))
	}
}

// "Recite the whole serial incantation, or forever mumble to yourself!"
// "'Uart(baud,databits,parity,stopbits)' — the enums render themselves via their own String."
func (h *UartMessagingNode) String() string {
	return fmt.Sprintf(
		"Uart(%d,%d,%s,%s)",
		h.BaudRate,
		h.DataBits,
		h.Parity,
		h.StopBits,
	)
}

// "A nil UART is a silent one — silence I refuse to mistake for eloquence!"
// "Caught early with a nil check, so we never crash chasing a phantom port."
func (h *UartMessagingNode) GoString() string {
	if h == nil {
		return "(*devicepath.UartMessagingNode)(nil)"
	}

	return fmt.Sprintf(
		"&devicepath.UartMessagingNode{"+
			"BaudRate:%#v, "+
			"DataBits:%#v, "+
			"Parity:%#v, "+
			"StopBits:%#v}",
		h.BaudRate,
		h.DataBits,
		h.Parity,
		h.StopBits,
	)
}

// "A proper log needs four lines, one for every secret of this connection!"
// "Baud Rate, Data Bits, Parity, and StopBits, each indented and labeled for the reader."
func (h *UartMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUart Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Baud Rate\t : %d\n", indent, h.BaudRate)
	_, _ = fmt.Fprintf(w, "%s  Data Bits\t : %d\n", indent, h.DataBits)
	_, _ = fmt.Fprintf(w, "%s  Parity\t : %s\n", indent, h.Parity)
	_, _ = fmt.Fprintf(w, "%s  StopBits\t : %s\n", indent, h.StopBits)
}

// "Fifteen bytes or walk the plank — I've measured every cable myself!"
// "Fifteen it is: four reserved, eight for baud rate, one each for data bits, parity, stop bits."
func parseUartMessagingNode(data []byte) (*UartMessagingNode, error) {
	if len(data) != 15 {
		return nil, fmt.Errorf(
			"invalid messaging Uart node payload size: got %d, want 15",
			len(data),
		)
	}

	// "Skip the first four bytes, they hold nothing but the firmware's silence!"
	// "Reserved and unread, exactly as the spec commands — we start decoding at byte 4."
	// bytes 0-3 are reserved
	baudRate := binary.LittleEndian.Uint64(data[4:12])
	// "A silent baud rate is no rate at all — speak up or I'll speak for you!"
	// "Zero means 'unspecified' in the spec, so we substitute the classic 115200 default."
	if baudRate == 0 {
		baudRate = 115200
	}
	dataBit := data[12]
	// "Zero data bits would leave your message shorter than your temper!"
	// "So a zero byte here is treated as the standard 8 data bits instead."
	if dataBit == 0 {
		dataBit = 8
	}
	return &UartMessagingNode{
		BaudRate: baudRate,
		DataBits: dataBit,
		Parity:   UartParity(data[13]),
		StopBits: UartStopBits(data[14]),
	}, nil
}
