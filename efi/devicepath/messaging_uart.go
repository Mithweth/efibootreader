package devicepath

import (
	"encoding/binary"
	"fmt"
	"io"
)

type UartParity uint8
type UartStopBits uint8

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

type UartMessagingNode struct {
	BaudRate uint64
	DataBits uint8
	Parity   UartParity
	StopBits UartStopBits
}

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

func (h *UartMessagingNode) String() string {
	return fmt.Sprintf(
		"Uart(%d,%d,%s,%s)",
		h.BaudRate,
		h.DataBits,
		h.Parity,
		h.StopBits,
	)
}

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

func (h *UartMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sUart Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  Baud Rate\t : %d\n", indent, h.BaudRate)
	_, _ = fmt.Fprintf(w, "%s  Data Bits\t : %d\n", indent, h.DataBits)
	_, _ = fmt.Fprintf(w, "%s  Parity\t : %s\n", indent, h.Parity)
	_, _ = fmt.Fprintf(w, "%s  StopBits\t : %s\n", indent, h.StopBits)
}

func parseUartMessagingNode(data []byte) (*UartMessagingNode, error) {
	if len(data) != 15 {
		return nil, fmt.Errorf(
			"invalid messaging Uart node payload size: got %d, want 15",
			len(data),
		)
	}

	// bytes 0-3 are reserved
	baudRate := binary.LittleEndian.Uint64(data[4:12])
	if baudRate == 0 {
		baudRate = 115200
	}
	dataBit := data[12]
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
