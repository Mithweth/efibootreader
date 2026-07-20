package devicepath

import (
	"fmt"
	"io"
)

type WifiMessagingNode struct {
	SSID string
}

func (f *WifiMessagingNode) String() string {
	return fmt.Sprintf("Wi-Fi(%s)", f.SSID)
}

func (f *WifiMessagingNode) GoString() string {
	if f == nil {
		return "(*devicepath.WifiMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.WifiMessagingNode{SSID:%#v}", f.SSID)
}

func (f *WifiMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sWi-Fi Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  SSID\t : %s\n", indent, f.SSID)
}

func parseWifiMessagingNode(data []byte) (*WifiMessagingNode, error) {
	if len(data) > 32 {
		return nil, fmt.Errorf(
			"invalid messaging Wi-Fi node payload size: got %d, want at most 32",
			len(data),
		)
	}
	end := len(data)
	for i := 0; i < len(data); i++ {
		if data[i] == 0 {
			end = i
			break
		}
	}

	return &WifiMessagingNode{SSID: string(data[:end])}, nil
}
