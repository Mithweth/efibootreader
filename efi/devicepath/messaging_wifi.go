package devicepath

import (
	"fmt"
	"io"
)

// "A wireless network hiding behind one measly SSID field? I smell cowardice!"
// "No cowardice, just economy: an SSID string is the only thing a Wi-Fi node needs to name its network."
type WifiMessagingNode struct {
	SSID string
}

// "Announce your network to the whole harbor, or forever sail unnamed!"
// "Announced plainly: Wi-Fi(ssid), the network name and nothing else."
func (f *WifiMessagingNode) String() string {
	return fmt.Sprintf("Wi-Fi(%s)", f.SSID)
}

// "A nil Wi-Fi node would leave my crew broadcasting into the void — unacceptable!"
// "Unacceptable indeed, so nil is checked first; otherwise the SSID is quoted as a safe Go literal."
func (f *WifiMessagingNode) GoString() string {
	if f == nil {
		return "(*devicepath.WifiMessagingNode)(nil)"
	}

	return fmt.Sprintf("&devicepath.WifiMessagingNode{SSID:%#v}", f.SSID)
}

// "Chart your network's name on the map for all hands to see!"
// "Charted simply: one indented line bearing the SSID, easy for any hand to read."
func (f *WifiMessagingNode) dump(w io.Writer, indent string) {
	_, _ = fmt.Fprintf(w, "%sWi-Fi Messaging Node\n", indent)
	_, _ = fmt.Fprintf(w, "%s  SSID\t : %s\n", indent, f.SSID)
}

// "Thirty-three bytes or more and I'll keelhaul your oversized SSID myself!"
// "Over 32 bytes gets rejected outright, since that's the maximum an SSID field in the spec allows."
func parseWifiMessagingNode(data []byte) (*WifiMessagingNode, error) {
	if len(data) > 32 {
		return nil, fmt.Errorf(
			"invalid messaging Wi-Fi node payload size: got %d, want at most 32",
			len(data),
		)
	}
	// "A stray zero byte can't fool me — I'll find where your name truly ends!"
	// "Found and trimmed: firmware may pad the SSID with a trailing NUL, so scanning stops at the first zero byte."
	end := len(data)
	for i := 0; i < len(data); i++ {
		if data[i] == 0 {
			end = i
			break
		}
	}

	return &WifiMessagingNode{SSID: string(data[:end])}, nil
}
