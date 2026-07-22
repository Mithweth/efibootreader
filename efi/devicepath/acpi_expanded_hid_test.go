package devicepath

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/Mithweth/efibootreader/identifiers"
)

const (
	testHID identifiers.EISAID = 0x0A0341D0 // PNP0A03
	testCID identifiers.EISAID = 0x0A0841D0 // PNP0A08
)

func expandedHidAcpiPayload(
	hid identifiers.EISAID,
	uid uint32,
	cid identifiers.EISAID,
	hidStr string,
	uidStr string,
	cidStr string,
) []byte {
	data := make([]byte, 12)

	binary.LittleEndian.PutUint32(data[0:4], uint32(hid))
	binary.LittleEndian.PutUint32(data[4:8], uid)
	binary.LittleEndian.PutUint32(data[8:12], uint32(cid))

	data = append(data, hidStr...)
	data = append(data, 0)

	data = append(data, uidStr...)
	data = append(data, 0)

	data = append(data, cidStr...)
	data = append(data, 0)

	return data
}

func TestParseExpandedHidAcpiNodeValidation(t *testing.T) {
	tests := []struct {
		name    string
		hid     identifiers.EISAID
		uid     uint32
		cid     identifiers.EISAID
		hidStr  string
		uidStr  string
		cidStr  string
		wantErr string
	}{
		{
			name:    "HID and HIDSTR present",
			hid:     testHID,
			hidStr:  "ACME0001",
			wantErr: "",
		},
		{
			name:    "HID only present",
			hid:     testHID,
			wantErr: "",
		},
		{
			name:    "HIDSTR only present",
			hidStr:  "ACME0001",
			wantErr: "",
		},
		{
			name:    "HID and HIDSTR empty",
			wantErr: "both HID and HIDSTR are empty",
		},
		{
			name:    "CID and CIDSTR present",
			hid:     testHID,
			cid:     testCID,
			cidStr:  "ACME0002",
			wantErr: "CID and CIDSTR are both present",
		},
		{
			name:    "CID only present",
			hid:     testHID,
			cid:     testCID,
			wantErr: "",
		},
		{
			name:    "CIDSTR only present",
			hid:     testHID,
			cidStr:  "ACME0002",
			wantErr: "",
		},
		{
			name:    "CID and CIDSTR empty",
			hid:     testHID,
			wantErr: "",
		},
		{
			name:    "UID greater than zero and UIDSTR present",
			hid:     testHID,
			uid:     42,
			uidStr:  "forty-two",
			wantErr: "UID and UIDSTR are both present",
		},
		{
			name:    "numeric UID only present",
			hid:     testHID,
			uid:     42,
			wantErr: "",
		},
		{
			name:    "UIDSTR only present",
			hid:     testHID,
			uidStr:  "forty-two",
			wantErr: "",
		},
		{
			name:    "UID zero and UIDSTR empty",
			hid:     testHID,
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := expandedHidAcpiPayload(
				tt.hid,
				tt.uid,
				tt.cid,
				tt.hidStr,
				tt.uidStr,
				tt.cidStr,
			)

			got, err := parseExpandedHidAcpiNode(data)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf(
						"parseExpandedHidAcpiNode() returned no error, want error containing %q",
						tt.wantErr,
					)
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf(
						"parseExpandedHidAcpiNode() error = %q, want error containing %q",
						err,
						tt.wantErr,
					)
				}

				if got != nil {
					t.Errorf(
						"parseExpandedHidAcpiNode() node = %#v, want nil on error",
						got,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf("parseExpandedHidAcpiNode() unexpected error: %v", err)
			}

			if got.HID != tt.hid {
				t.Errorf("HID = %#v, want %#v", got.HID, tt.hid)
			}
			if got.UID != tt.uid {
				t.Errorf("UID = %d, want %d", got.UID, tt.uid)
			}
			if got.CID != tt.cid {
				t.Errorf("CID = %#v, want %#v", got.CID, tt.cid)
			}
			if got.HIDStr != tt.hidStr {
				t.Errorf("HIDStr = %q, want %q", got.HIDStr, tt.hidStr)
			}
			if got.UIDStr != tt.uidStr {
				t.Errorf("UIDStr = %q, want %q", got.UIDStr, tt.uidStr)
			}
			if got.CIDStr != tt.cidStr {
				t.Errorf("CIDStr = %q, want %q", got.CIDStr, tt.cidStr)
			}
		})
	}
}

func TestParseExpandedHidAcpiNodeStringFields(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr string
	}{
		{
			name:    "payload shorter than fixed fields and terminators",
			data:    make([]byte, 14),
			wantErr: "want at least 15",
		},
		{
			name: "missing CIDSTR terminator",
			data: append(
				expandedHidAcpiPayload(testHID, 0, 0, "", "", "")[:14],
				'A',
			),
			wantErr: "expected three NUL-terminated strings",
		},
		{
			name: "unexpected trailing bytes",
			data: append(
				expandedHidAcpiPayload(
					testHID,
					0,
					0,
					"",
					"",
					"",
				),
				0x42,
			),
			wantErr: "expected three NUL-terminated strings",
		},
		{
			name: "too many string fields",
			data: append(
				expandedHidAcpiPayload(
					testHID,
					0,
					0,
					"",
					"",
					"",
				),
				0,
			),
			wantErr: "expected three NUL-terminated strings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExpandedHidAcpiNode(tt.data)
			if err == nil {
				t.Fatalf(
					"parseExpandedHidAcpiNode() returned node %#v, want error containing %q",
					got,
					tt.wantErr,
				)
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf(
					"parseExpandedHidAcpiNode() error = %q, want error containing %q",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestExpandedHidAcpiNodeString(t *testing.T) {
	tests := []struct {
		name string
		node ExpandedHidAcpiNode
		want string
	}{
		{
			name: "numeric HID CID and textual UID use AcpiExp",
			node: ExpandedHidAcpiNode{
				HID:    testHID,
				CID:    testCID,
				UIDStr: "ROOT",
			},
			want: "AcpiExp(PNP0A03,PNP0A08,ROOT)",
		},
		{
			name: "numeric UID uses AcpiExp",
			node: ExpandedHidAcpiNode{
				HID: testHID,
				CID: testCID,
				UID: 42,
			},
			want: "AcpiExp(PNP0A03,PNP0A08,42)",
		},
		{
			name: "textual HID uses full AcpiEx",
			node: ExpandedHidAcpiNode{
				HIDStr: "ACME0001",
			},
			want: "AcpiEx(,,0,ACME0001,,)",
		},
		{
			name: "numeric and textual HID use full AcpiEx",
			node: ExpandedHidAcpiNode{
				HID:    testHID,
				HIDStr: "ACME0001",
			},
			want: "AcpiEx(PNP0A03,,0,ACME0001,,)",
		},
		{
			name: "textual CID uses full AcpiEx",
			node: ExpandedHidAcpiNode{
				HID:    testHID,
				CIDStr: "ACME0002",
			},
			want: "AcpiEx(PNP0A03,,0,,ACME0002,)",
		},
		{
			name: "all optional fields empty",
			node: ExpandedHidAcpiNode{
				HID: testHID,
			},
			want: "AcpiEx(PNP0A03,,0,,,)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
