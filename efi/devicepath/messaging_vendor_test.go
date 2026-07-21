package devicepath

import (
	"encoding/binary"
	"fmt"
	"strings"
	"testing"

	"github.com/Mithweth/efibootreader/identifiers"
)

func TestParseVendorMessagingNodeRejectsShortPayload(t *testing.T) {
	t.Parallel()

	_, err := parseVendorMessagingNode(make([]byte, 15))
	if err == nil {
		t.Fatal("parseVendorMessagingNode() returned nil error")
	}

	want := "invalid vendor messaging node payload size: got 15, want at least 16"
	if err.Error() != want {
		t.Fatalf("unexpected error:\ngot:  %q\nwant: %q", err, want)
	}
}

func TestParseVendorMessagingNodeUnknownGUID(t *testing.T) {
	t.Parallel()

	guid := identifiers.MustParseEFIGUID(
		"12345678-1234-5678-90ab-cdef01234567",
	)
	vendorData := []byte{0xde, 0xad, 0xbe, 0xef}

	data := appendGUIDBytes(t, guid, vendorData)

	node, err := parseVendorMessagingNode(data)
	if err != nil {
		t.Fatalf("parseVendorMessagingNode() returned error: %v", err)
	}

	got, ok := node.(*GenericVendorMessagingNode)
	if !ok {
		t.Fatalf(
			"unexpected node type: got %T, want *GenericVendorMessagingNode",
			node,
		)
	}

	if got.GUID != guid {
		t.Errorf("unexpected GUID: got %v, want %v", got.GUID, guid)
	}

	if string(got.Data) != string(vendorData) {
		t.Errorf(
			"unexpected vendor data: got %x, want %x",
			got.Data,
			vendorData,
		)
	}

	// Ensure the returned node owns its payload.
	data[len(data)-1] = 0

	if got.Data[len(got.Data)-1] != 0xef {
		t.Errorf(
			"GenericVendorMessagingNode.Data aliases the input buffer: got %x",
			got.Data,
		)
	}
}

func TestGenericVendorMessagingNodeString(t *testing.T) {
	t.Parallel()

	guid := identifiers.MustParseEFIGUID(
		"12345678-1234-5678-90ab-cdef01234567",
	)

	tests := []struct {
		name string
		node *GenericVendorMessagingNode
		want string
	}{
		{
			name: "without data",
			node: &GenericVendorMessagingNode{
				GUID: guid,
			},
			want: fmt.Sprintf("VenMsg(%s)", guid),
		},
		{
			name: "with data",
			node: &GenericVendorMessagingNode{
				GUID: guid,
				Data: []byte{0xde, 0xad, 0xbe, 0xef},
			},
			want: fmt.Sprintf("VenMsg(%s,deadbeef)", guid),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.node.String(); got != tt.want {
				t.Fatalf("String():\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func TestGenericVendorMessagingNodeGoString(t *testing.T) {
	t.Parallel()

	var nilNode *GenericVendorMessagingNode

	if got, want := nilNode.GoString(),
		"(*devicepath.GenericVendorMessagingNode)(nil)"; got != want {
		t.Fatalf("nil GoString():\ngot:  %q\nwant: %q", got, want)
	}

	guid := identifiers.MustParseEFIGUID(
		"12345678-1234-5678-90ab-cdef01234567",
	)

	node := &GenericVendorMessagingNode{
		GUID: guid,
		Data: []byte{0xde, 0xad},
	}

	got := node.GoString()

	for _, fragment := range []string{
		"&devicepath.GenericVendorMessagingNode{",
		"GUID:",
		"Data:[]byte{0xde, 0xad}",
	} {
		if !strings.Contains(got, fragment) {
			t.Errorf("GoString() does not contain %q:\n%s", fragment, got)
		}
	}
}

func TestGenericVendorMessagingNodeDump(t *testing.T) {
	t.Parallel()

	guid := identifiers.MustParseEFIGUID(
		"12345678-1234-5678-90ab-cdef01234567",
	)

	node := &GenericVendorMessagingNode{
		GUID: guid,
		Data: []byte{0xde, 0xad, 0xbe, 0xef},
	}

	var output strings.Builder
	node.dump(&output, "  ")

	got := output.String()

	for _, fragment := range []string{
		"  Vendor Messaging Node\n",
		fmt.Sprintf("    GUID\t : %s", guid),
		"    Data\t : deadbeef\n",
	} {
		if !strings.Contains(got, fragment) {
			t.Errorf("dump() does not contain %q:\n%s", fragment, got)
		}
	}
}

func TestSasMessagingDeviceInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		value             SasMessagingDeviceInfo
		informationLength uint8
		deviceType        uint8
		topology          uint8
		driveBay          uint16
		sasSata           string
		internal          bool
		location          string
		topologyString    string
	}{
		{
			name:              "SAS internal direct without drive bay",
			value:             0x0001,
			informationLength: 1,
			deviceType:        0,
			topology:          0,
			driveBay:          1,
			sasSata:           "SAS",
			internal:          true,
			location:          "Internal",
			topologyString:    "Direct",
		},
		{
			name:              "SATA internal expanded drive bay 3",
			value:             0x0252,
			informationLength: 2,
			deviceType:        1,
			topology:          1,
			driveBay:          3,
			sasSata:           "SATA",
			internal:          true,
			location:          "Internal",
			topologyString:    "Expanded",
		},
		{
			name:              "SAS external reserved topology",
			value:             0x00a1,
			informationLength: 1,
			deviceType:        2,
			topology:          2,
			driveBay:          1,
			sasSata:           "SAS",
			internal:          false,
			location:          "External",
			topologyString:    "2",
		},
		{
			name:              "SATA external reserved topology",
			value:             0x00f1,
			informationLength: 1,
			deviceType:        3,
			topology:          3,
			driveBay:          1,
			sasSata:           "SATA",
			internal:          false,
			location:          "External",
			topologyString:    "3",
		},
		{
			name:              "maximum drive bay",
			value:             0xff02,
			informationLength: 2,
			deviceType:        0,
			topology:          0,
			driveBay:          256,
			sasSata:           "SAS",
			internal:          true,
			location:          "Internal",
			topologyString:    "Direct",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.value.InformationLength(); got != tt.informationLength {
				t.Errorf(
					"InformationLength(): got %d, want %d",
					got,
					tt.informationLength,
				)
			}

			if got := tt.value.DeviceType(); got != tt.deviceType {
				t.Errorf("DeviceType(): got %d, want %d", got, tt.deviceType)
			}

			if got := tt.value.Topology(); got != tt.topology {
				t.Errorf("Topology(): got %d, want %d", got, tt.topology)
			}

			if got := tt.value.DriveBay(); got != tt.driveBay {
				t.Errorf("DriveBay(): got %d, want %d", got, tt.driveBay)
			}

			if got := tt.value.SasSataString(); got != tt.sasSata {
				t.Errorf("SasSataString(): got %q, want %q", got, tt.sasSata)
			}

			if got := tt.value.IsInternal(); got != tt.internal {
				t.Errorf("IsInternal(): got %t, want %t", got, tt.internal)
			}

			if got := tt.value.LocationString(); got != tt.location {
				t.Errorf("LocationString(): got %q, want %q", got, tt.location)
			}

			if got := tt.value.TopologyString(); got != tt.topologyString {
				t.Errorf(
					"TopologyString(): got %q, want %q",
					got,
					tt.topologyString,
				)
			}
		})
	}
}

func TestSasMessagingDeviceInfoGoString(t *testing.T) {
	t.Parallel()

	value := SasMessagingDeviceInfo(0x0252)

	got := value.GoString()
	want := "devicepath.SasMessagingDeviceInfo(0x252)"

	if got != want {
		t.Fatalf("GoString():\ngot:  %q\nwant: %q", got, want)
	}
}

func TestSasMessagingDeviceInfoDump(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		value      SasMessagingDeviceInfo
		contains   []string
		notContain []string
	}{
		{
			name:  "no information",
			value: 0,
			notContain: []string{
				"Device Info",
				"Device Type",
				"Drive Bay",
			},
		},
		{
			name:  "one byte information",
			value: 0x0051,
			contains: []string{
				"Device Info (0b0000000001010001)",
				"Device Type\t : SATA",
				"Location\t : Internal",
				"Connect\t : Expanded",
			},
			notContain: []string{
				"Drive Bay",
			},
		},
		{
			name:  "two byte information with drive bay",
			value: 0x0252,
			contains: []string{
				"Device Info (0b0000001001010010)",
				"Device Type\t : SATA",
				"Location\t : Internal",
				"Connect\t : Expanded",
				"Drive Bay\t\t : 3",
			},
		},
		{
			name:  "reserved information length",
			value: 0x0003,
			contains: []string{
				"Device Info (0b0000000000000011)",
			},
			notContain: []string{
				"Device Type",
				"Location",
				"Connect",
				"Drive Bay",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var output strings.Builder
			tt.value.dump(&output, "")

			got := output.String()

			for _, fragment := range tt.contains {
				if !strings.Contains(got, fragment) {
					t.Errorf("dump() does not contain %q:\n%s", fragment, got)
				}
			}

			for _, fragment := range tt.notContain {
				if strings.Contains(got, fragment) {
					t.Errorf("dump() unexpectedly contains %q:\n%s", fragment, got)
				}
			}
		})
	}
}

func TestSasMessagingNodeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		node SasMessagingNode
		want string
	}{
		{
			name: "no topology",
			node: SasMessagingNode{
				Address:            0x1122334455667788,
				LogicalUnitNumber:  0x99,
				RelativeTargetPort: 7,
				DeviceInfo:         0,
			},
			want: "SAS(0x1122334455667788,0x99,7)",
		},
		{
			name: "one byte topology",
			node: SasMessagingNode{
				Address:            0x1122334455667788,
				LogicalUnitNumber:  0x99,
				RelativeTargetPort: 7,
				DeviceInfo:         0x0051,
			},
			want: "SAS(0x1122334455667788,0x99,7,SATA,Internal,Expanded)",
		},
		{
			name: "two byte topology internal with drive bay",
			node: SasMessagingNode{
				Address:            0x1122334455667788,
				LogicalUnitNumber:  0x99,
				RelativeTargetPort: 7,
				DeviceInfo:         0x0252,
			},
			want: "SAS(0x1122334455667788,0x99,7,SATA,Internal,Expanded,3)",
		},
		{
			name: "reserved topology value",
			node: SasMessagingNode{
				Address:            1,
				LogicalUnitNumber:  2,
				RelativeTargetPort: 3,
				DeviceInfo:         0x0003,
			},
			want: "SAS(0x1,0x2,3,0x3)",
		},
		{
			name: "reserved field after explicit drive bay",
			node: SasMessagingNode{
				Reserved:           0xaabbccdd,
				Address:            1,
				LogicalUnitNumber:  2,
				RelativeTargetPort: 3,
				DeviceInfo:         0x0252,
			},
			want: "SAS(0x1,0x2,3,SATA,Internal,Expanded,3,0xaabbccdd)",
		},
		{
			name: "reserved field with drive bay placeholder",
			node: SasMessagingNode{
				Reserved:           0xaabbccdd,
				Address:            1,
				LogicalUnitNumber:  2,
				RelativeTargetPort: 3,
				DeviceInfo:         0x0051,
			},
			want: "SAS(0x1,0x2,3,SATA,Internal,Expanded,0,0xaabbccdd)",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.node.String(); got != tt.want {
				t.Fatalf("String():\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func TestSasMessagingNodeGoString(t *testing.T) {
	t.Parallel()

	var nilNode *SasMessagingNode

	if got, want := nilNode.GoString(),
		"(*devicepath.SasMessagingNode)(nil)"; got != want {
		t.Fatalf("nil GoString():\ngot:  %q\nwant: %q", got, want)
	}

	node := &SasMessagingNode{
		Reserved:           1,
		Address:            2,
		LogicalUnitNumber:  3,
		DeviceInfo:         4,
		RelativeTargetPort: 5,
	}

	got := node.GoString()

	for _, fragment := range []string{
		"&devicepath.SasMessagingNode{",
		"Reserved:0x1",
		"Address:0x2",
		"DeviceInfo:devicepath.SasMessagingDeviceInfo(0x4)",
		"LogicalUnitNumber:0x3",
		"RelativeTargetPort:0x5",
	} {
		if !strings.Contains(got, fragment) {
			t.Errorf("GoString() does not contain %q:\n%s", fragment, got)
		}
	}
}

func TestSasMessagingNodeDump(t *testing.T) {
	t.Parallel()

	node := &SasMessagingNode{
		Reserved:           0xaabbccdd,
		Address:            0x1122334455667788,
		LogicalUnitNumber:  0x99,
		DeviceInfo:         0x0252,
		RelativeTargetPort: 7,
	}

	var output strings.Builder
	node.dump(&output, "  ")

	got := output.String()

	for _, fragment := range []string{
		"  SAS Messaging Node\n",
		"    Address\t\t : 0x1122334455667788\n",
		"    Reserved\t\t : 0xaabbccdd\n",
		"    Logical Unit Number\t : 0x99\n",
		"    Relative Target Port\t : 7\n",
		"    Device Info (0b0000001001010010)\n",
		"      Drive Bay\t\t : 3\n",
	} {
		if !strings.Contains(got, fragment) {
			t.Errorf("dump() does not contain %q:\n%s", fragment, got)
		}
	}
}

func TestParseSasMessagingNodeRejectsInvalidSize(t *testing.T) {
	t.Parallel()

	for _, size := range []int{0, 23, 25} {
		size := size

		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			t.Parallel()

			_, err := parseSasMessagingNode(make([]byte, size))
			if err == nil {
				t.Fatal("parseSasMessagingNode() returned nil error")
			}

			want := fmt.Sprintf(
				"invalid SAS messaging node payload size: got %d, want 24",
				size,
			)

			if err.Error() != want {
				t.Fatalf("unexpected error:\ngot:  %q\nwant: %q", err, want)
			}
		})
	}
}

func TestParseSasMessagingNode(t *testing.T) {
	t.Parallel()

	data := make([]byte, 24)

	binary.LittleEndian.PutUint32(data[0:4], 0xaabbccdd)
	binary.LittleEndian.PutUint64(data[4:12], 0x1122334455667788)
	binary.LittleEndian.PutUint64(data[12:20], 0x0102030405060708)
	binary.LittleEndian.PutUint16(data[20:22], 0x0252)
	binary.LittleEndian.PutUint16(data[22:24], 0x1234)

	got, err := parseSasMessagingNode(data)
	if err != nil {
		t.Fatalf("parseSasMessagingNode() returned error: %v", err)
	}

	want := &SasMessagingNode{
		Reserved:           0xaabbccdd,
		Address:            0x1122334455667788,
		LogicalUnitNumber:  0x0102030405060708,
		DeviceInfo:         0x0252,
		RelativeTargetPort: 0x1234,
	}

	if *got != *want {
		t.Fatalf("unexpected result:\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestParseVendorMessagingNodeDispatchesSAS(t *testing.T) {
	t.Parallel()

	payload := make([]byte, 24)
	binary.LittleEndian.PutUint64(payload[4:12], 0x1122334455667788)
	binary.LittleEndian.PutUint64(payload[12:20], 0x99)
	binary.LittleEndian.PutUint16(payload[20:22], 0x0252)
	binary.LittleEndian.PutUint16(payload[22:24], 7)

	data := appendGUIDBytes(t, sasDevicePathGUID, payload)

	node, err := parseVendorMessagingNode(data)
	if err != nil {
		t.Fatalf("parseVendorMessagingNode() returned error: %v", err)
	}

	got, ok := node.(*SasMessagingNode)
	if !ok {
		t.Fatalf("unexpected node type: got %T, want *SasMessagingNode", node)
	}

	if got.Address != 0x1122334455667788 {
		t.Errorf("unexpected address: got %#x", got.Address)
	}

	if got.DeviceInfo != 0x0252 {
		t.Errorf("unexpected device info: got %#x", got.DeviceInfo)
	}
}

func TestUartFlowControlMessagingTypeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value UartFlowControlMessagingType
		want  string
	}{
		{UartFlowControlMessagingTypeNone, "None"},
		{UartFlowControlMessagingTypeHardware, "Hardware"},
		{UartFlowControlMessagingTypeXonXoff, "XonXoff"},
		{UartFlowControlMessagingTypeHardwareXonXoff, "Hardware+XonXoff"},
		{UartFlowControlMessagingType(4), "4"},
		{UartFlowControlMessagingType(0x80000000), "2147483648"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("%#x", uint32(tt.value)), func(t *testing.T) {
			t.Parallel()

			if got := tt.value.String(); got != tt.want {
				t.Fatalf("String(): got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUartFlowControlMessagingTypeGoString(t *testing.T) {
	t.Parallel()

	value := UartFlowControlMessagingTypeHardwareXonXoff

	got := value.GoString()
	want := "devicepath.UartFlowControlMessagingType(0x3)"

	if got != want {
		t.Fatalf("GoString():\ngot:  %q\nwant: %q", got, want)
	}
}

func TestUartFlowControlMessagingNodeString(t *testing.T) {
	t.Parallel()

	node := &UartFlowControlMessagingNode{
		FlowControlMap: UartFlowControlMessagingTypeHardwareXonXoff,
	}

	got := node.String()
	want := "UartFlowCtrl(Hardware+XonXoff)"

	if got != want {
		t.Fatalf("String():\ngot:  %q\nwant: %q", got, want)
	}
}

func TestUartFlowControlMessagingNodeGoString(t *testing.T) {
	t.Parallel()

	var nilNode *UartFlowControlMessagingNode

	if got, want := nilNode.GoString(),
		"(*devicepath.UartFlowControlMessagingNode)(nil)"; got != want {
		t.Fatalf("nil GoString():\ngot:  %q\nwant: %q", got, want)
	}

	node := &UartFlowControlMessagingNode{
		FlowControlMap: UartFlowControlMessagingTypeHardware,
	}

	got := node.GoString()

	for _, fragment := range []string{
		"&devicepath.UartFlowControlMessagingNode{",
		"FlowControlMap:devicepath.UartFlowControlMessagingType(0x1)",
	} {
		if !strings.Contains(got, fragment) {
			t.Errorf("GoString() does not contain %q:\n%s", fragment, got)
		}
	}
}

func TestUartFlowControlMessagingNodeDump(t *testing.T) {
	t.Parallel()

	node := &UartFlowControlMessagingNode{
		FlowControlMap: UartFlowControlMessagingTypeHardwareXonXoff,
	}

	var output strings.Builder
	node.dump(&output, "  ")

	got := output.String()

	for _, fragment := range []string{
		"  UART Flow Control Messaging Node\n",
		"    Flow ControlMap\t : Hardware+XonXoff",
		"00000000000000000000000000000011",
	} {
		if !strings.Contains(got, fragment) {
			t.Errorf("dump() does not contain %q:\n%s", fragment, got)
		}
	}
}

func TestParseUartFlowControlMessagingNodeRejectsInvalidSize(t *testing.T) {
	t.Parallel()

	for _, size := range []int{0, 3, 5} {
		size := size

		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			t.Parallel()

			_, err := parseUartFlowControlMessagingNode(make([]byte, size))
			if err == nil {
				t.Fatal(
					"parseUartFlowControlMessagingNode() returned nil error",
				)
			}

			want := fmt.Sprintf(
				"invalid Uart flow control messaging node payload size: "+
					"got %d, want 4",
				size,
			)

			if err.Error() != want {
				t.Fatalf("unexpected error:\ngot:  %q\nwant: %q", err, want)
			}
		})
	}
}

func TestParseUartFlowControlMessagingNode(t *testing.T) {
	t.Parallel()

	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, 0x80000003)

	got, err := parseUartFlowControlMessagingNode(data)
	if err != nil {
		t.Fatalf(
			"parseUartFlowControlMessagingNode() returned error: %v",
			err,
		)
	}

	want := UartFlowControlMessagingType(0x80000003)

	if got.FlowControlMap != want {
		t.Fatalf(
			"unexpected flow-control map: got %#x, want %#x",
			got.FlowControlMap,
			want,
		)
	}
}

func TestParseVendorMessagingNodeDispatchesUartFlowControl(t *testing.T) {
	t.Parallel()

	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(
		payload,
		uint32(UartFlowControlMessagingTypeHardwareXonXoff),
	)

	data := appendGUIDBytes(t, uartFlowControlGUID, payload)

	node, err := parseVendorMessagingNode(data)
	if err != nil {
		t.Fatalf("parseVendorMessagingNode() returned error: %v", err)
	}

	got, ok := node.(*UartFlowControlMessagingNode)
	if !ok {
		t.Fatalf(
			"unexpected node type: got %T, "+
				"want *UartFlowControlMessagingNode",
			node,
		)
	}

	if got.FlowControlMap != UartFlowControlMessagingTypeHardwareXonXoff {
		t.Errorf(
			"unexpected flow-control map: got %#x",
			got.FlowControlMap,
		)
	}
}

// appendGUIDBytes serializes guid using its EFI binary representation and
// appends payload.
func appendGUIDBytes(
	t *testing.T,
	guid identifiers.GUID,
	payload []byte,
) []byte {
	t.Helper()

	data := make([]byte, 16+len(payload))

	data[0] = guid[3]
	data[1] = guid[2]
	data[2] = guid[1]
	data[3] = guid[0]

	data[4] = guid[5]
	data[5] = guid[4]

	data[6] = guid[7]
	data[7] = guid[6]

	copy(data[8:16], guid[8:])
	copy(data[16:], payload)

	return data
}
