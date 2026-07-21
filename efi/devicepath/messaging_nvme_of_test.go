package devicepath

import (
	"bytes"
	"strings"
	"testing"
)

func makeNvmeOfPayload(
	nidt NvmeOfNamespaceIdentifierType,
	nid [16]byte,
	nqn []byte,
) []byte {
	data := make([]byte, 0, 17+len(nqn))
	data = append(data, byte(nidt))
	data = append(data, nid[:]...)
	data = append(data, nqn...)

	return data
}

func TestNvmeOfNamespaceMessagingNodeNIDString(t *testing.T) {
	tests := []struct {
		name string
		nidt NvmeOfNamespaceIdentifierType
		nid  [16]byte
		want string
	}{
		{
			name: "EUI64",
			nidt: NvmeOfNIDTypeEUI64,
			nid: [16]byte{
				0x01, 0x23, 0x45, 0x67,
				0x89, 0xab, 0xcd, 0xef,
				0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff,
			},
			want: "0123456789abcdef",
		},
		{
			name: "NGUID",
			nidt: NvmeOfNIDTypeNGUID,
			nid: [16]byte{
				0x00, 0x11, 0x22, 0x33,
				0x44, 0x55, 0x66, 0x77,
				0x88, 0x99, 0xaa, 0xbb,
				0xcc, 0xdd, 0xee, 0xff,
			},
			want: "00112233445566778899aabbccddeeff",
		},
		{
			name: "UUID",
			nidt: NvmeOfNIDTypeUUID,
			nid: [16]byte{
				0x4e, 0xff, 0x7f, 0x8e,
				0xd3, 0x53,
				0x4e, 0x9b,
				0xa4, 0xec,
				0xde, 0xea, 0x8e, 0xab, 0x84, 0xd7,
			},
			want: "urn:uuid:4eff7f8e-d353-4e9b-a4ec-deea8eab84d7",
		},
		{
			name: "unknown type",
			nidt: NvmeOfNamespaceIdentifierType(0xff),
			nid: [16]byte{
				0x00, 0x11, 0x22, 0x33,
				0x44, 0x55, 0x66, 0x77,
				0x88, 0x99, 0xaa, 0xbb,
				0xcc, 0xdd, 0xee, 0xff,
			},
			want: "00112233445566778899aabbccddeeff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &NvmeOfNamespaceMessagingNode{
				NIDT: tt.nidt,
				NID:  tt.nid,
			}

			if got := node.NIDString(); got != tt.want {
				t.Fatalf(
					"NIDString() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestNvmeOfNamespaceMessagingNodeString(t *testing.T) {
	node := &NvmeOfNamespaceMessagingNode{
		NIDT: NvmeOfNIDTypeUUID,
		NID: [16]byte{
			0x4e, 0xff, 0x7f, 0x8e,
			0xd3, 0x53,
			0x4e, 0x9b,
			0xa4, 0xec,
			0xde, 0xea, 0x8e, 0xab, 0x84, 0xd7,
		},
		SubsystemNQN: "nqn.2014-08.org.nvmexpress:uuid:test",
	}

	want := "NVMEoF(" +
		"nqn.2014-08.org.nvmexpress:uuid:test," +
		"urn:uuid:4eff7f8e-d353-4e9b-a4ec-deea8eab84d7" +
		")"

	if got := node.String(); got != want {
		t.Fatalf(
			"String() = %q, want %q",
			got,
			want,
		)
	}
}

func TestNvmeOfNamespaceMessagingNodeGoStringNil(t *testing.T) {
	var node *NvmeOfNamespaceMessagingNode

	want := "(*devicepath.NvmeOfNamespaceMessagingNode)(nil)"

	if got := node.GoString(); got != want {
		t.Fatalf(
			"GoString() = %q, want %q",
			got,
			want,
		)
	}
}

func TestNvmeOfNamespaceMessagingNodeGoString(t *testing.T) {
	node := &NvmeOfNamespaceMessagingNode{
		NIDT: NvmeOfNIDTypeEUI64,
		NID: [16]byte{
			0x01, 0x23, 0x45, 0x67,
			0x89, 0xab, 0xcd, 0xef,
		},
		SubsystemNQN: "nqn.example",
	}

	got := node.GoString()

	want := "&devicepath.NvmeOfNamespaceMessagingNode{" +
		"NIDT:devicepath.NvmeOfNamespaceIdentifierType(0x1), " +
		"NID:[16]uint8{" +
		"0x1, 0x23, 0x45, 0x67, " +
		"0x89, 0xab, 0xcd, 0xef, " +
		"0x0, 0x0, 0x0, 0x0, " +
		"0x0, 0x0, 0x0, 0x0}, " +
		`SubsystemNQN:"nqn.example"}`

	if got != want {
		t.Fatalf(
			"GoString() mismatch:\ngot:  %s\nwant: %s",
			got,
			want,
		)
	}
}

func TestNvmeOfNamespaceMessagingNodeDump(t *testing.T) {
	node := &NvmeOfNamespaceMessagingNode{
		NIDT: NvmeOfNIDTypeNGUID,
		NID: [16]byte{
			0x00, 0x11, 0x22, 0x33,
			0x44, 0x55, 0x66, 0x77,
			0x88, 0x99, 0xaa, 0xbb,
			0xcc, 0xdd, 0xee, 0xff,
		},
		SubsystemNQN: "nqn.example",
	}

	var buffer bytes.Buffer
	node.dump(&buffer, "  ")

	got := buffer.String()

	expectedLines := []string{
		"  NVMe-oF Namespace Messaging Node",
		"    NIDT",
		"2 (0x2)",
		"    NID",
		"00112233445566778899aabbccddeeff",
		"    Subsystem NQN",
		"nqn.example",
	}

	for _, line := range expectedLines {
		if !strings.Contains(got, line) {
			t.Errorf(
				"dump output does not contain %q:\n%s",
				line,
				got,
			)
		}
	}
}

func TestParseNvmeOfNamespaceMessagingNode(t *testing.T) {
	nid := [16]byte{
		0x4e, 0xff, 0x7f, 0x8e,
		0xd3, 0x53,
		0x4e, 0x9b,
		0xa4, 0xec,
		0xde, 0xea, 0x8e, 0xab, 0x84, 0xd7,
	}

	data := makeNvmeOfPayload(
		NvmeOfNIDTypeUUID,
		nid,
		append(
			[]byte("nqn.2014-08.org.nvmexpress:uuid:test"),
			0,
		),
	)

	node, err := parseNvmeOfNamespaceMessagingNode(data)
	if err != nil {
		t.Fatalf(
			"parseNvmeOfNamespaceMessagingNode() returned error: %v",
			err,
		)
	}

	if node.NIDT != NvmeOfNIDTypeUUID {
		t.Errorf(
			"NIDT = %#x, want %#x",
			node.NIDT,
			NvmeOfNIDTypeUUID,
		)
	}

	if node.NID != nid {
		t.Errorf(
			"NID = %#v, want %#v",
			node.NID,
			nid,
		)
	}

	wantNQN := "nqn.2014-08.org.nvmexpress:uuid:test"
	if node.SubsystemNQN != wantNQN {
		t.Errorf(
			"SubsystemNQN = %q, want %q",
			node.SubsystemNQN,
			wantNQN,
		)
	}
}

func TestParseNvmeOfNamespaceMessagingNodeStopsAtFirstNull(t *testing.T) {
	var nid [16]byte

	data := makeNvmeOfPayload(
		NvmeOfNIDTypeNGUID,
		nid,
		[]byte{'f', 'o', 'o', 0, 'b', 'a', 'r', 0},
	)

	node, err := parseNvmeOfNamespaceMessagingNode(data)
	if err != nil {
		t.Fatalf(
			"parseNvmeOfNamespaceMessagingNode() returned error: %v",
			err,
		)
	}

	if node.SubsystemNQN != "foo" {
		t.Fatalf(
			"SubsystemNQN = %q, want %q",
			node.SubsystemNQN,
			"foo",
		)
	}
}

func TestParseNvmeOfNamespaceMessagingNodeEmptyNQN(t *testing.T) {
	var nid [16]byte

	data := makeNvmeOfPayload(
		NvmeOfNIDTypeEUI64,
		nid,
		[]byte{0},
	)

	node, err := parseNvmeOfNamespaceMessagingNode(data)
	if err != nil {
		t.Fatalf(
			"parseNvmeOfNamespaceMessagingNode() returned error: %v",
			err,
		)
	}

	if node.SubsystemNQN != "" {
		t.Fatalf(
			"SubsystemNQN = %q, want empty string",
			node.SubsystemNQN,
		)
	}
}

func TestParseNvmeOfNamespaceMessagingNodeMaximumNQNSize(t *testing.T) {
	var nid [16]byte

	// Le champ SubsystemNQN fait exactement 224 octets :
	// 223 octets de texte et un terminateur NUL.
	nqn := append(bytes.Repeat([]byte{'a'}, 223), 0)

	data := makeNvmeOfPayload(
		NvmeOfNIDTypeNGUID,
		nid,
		nqn,
	)

	node, err := parseNvmeOfNamespaceMessagingNode(data)
	if err != nil {
		t.Fatalf(
			"parseNvmeOfNamespaceMessagingNode() returned error: %v",
			err,
		)
	}

	if len(node.SubsystemNQN) != 223 {
		t.Fatalf(
			"len(SubsystemNQN) = %d, want 223",
			len(node.SubsystemNQN),
		)
	}
}

func TestParseNvmeOfNamespaceMessagingNodeErrors(t *testing.T) {
	var nid [16]byte

	tests := []struct {
		name      string
		data      []byte
		wantError string
	}{
		{
			name:      "empty payload",
			data:      nil,
			wantError: "want at least 18",
		},
		{
			name:      "payload too short",
			data:      make([]byte, 17),
			wantError: "want at least 18",
		},
		{
			name: "NQN is not null-terminated",
			data: makeNvmeOfPayload(
				NvmeOfNIDTypeNGUID,
				nid,
				[]byte("nqn.example"),
			),
			wantError: "not null-terminated",
		},
		{
			name: "NQN exceeds maximum size",
			data: makeNvmeOfPayload(
				NvmeOfNIDTypeNGUID,
				nid,
				append(bytes.Repeat([]byte{'a'}, 224), 0),
			),
			wantError: "want at most 224",
		},
		{
			name: "invalid UTF-8",
			data: makeNvmeOfPayload(
				NvmeOfNIDTypeNGUID,
				nid,
				[]byte{0xff, 0xfe, 0},
			),
			wantError: "not valid UTF-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseNvmeOfNamespaceMessagingNode(tt.data)
			if err == nil {
				t.Fatalf(
					"parseNvmeOfNamespaceMessagingNode() = %#v, want error",
					node,
				)
			}

			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf(
					"error = %q, want it to contain %q",
					err,
					tt.wantError,
				)
			}
		})
	}
}

func TestParseNvmeOfNamespaceMessagingNodeCopiesNID(t *testing.T) {
	nid := [16]byte{
		0x00, 0x11, 0x22, 0x33,
		0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xaa, 0xbb,
		0xcc, 0xdd, 0xee, 0xff,
	}

	data := makeNvmeOfPayload(
		NvmeOfNIDTypeNGUID,
		nid,
		[]byte{'n', 'q', 'n', 0},
	)

	node, err := parseNvmeOfNamespaceMessagingNode(data)
	if err != nil {
		t.Fatalf(
			"parseNvmeOfNamespaceMessagingNode() returned error: %v",
			err,
		)
	}

	data[1] = 0xff

	if node.NID[0] != 0x00 {
		t.Fatalf(
			"NID references input data: got first byte %#x, want %#x",
			node.NID[0],
			0x00,
		)
	}
}
