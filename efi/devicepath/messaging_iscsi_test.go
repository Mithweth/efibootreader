package devicepath

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func newIScsiMessagingNodePayload(
	protocol uint16,
	loginOptions uint16,
	lun [8]byte,
	portalGroup uint16,
	targetName []byte,
) []byte {
	data := make([]byte, 14+len(targetName))

	binary.LittleEndian.PutUint16(data[0:2], protocol)
	binary.LittleEndian.PutUint16(data[2:4], loginOptions)
	copy(data[4:12], lun[:])
	binary.LittleEndian.PutUint16(data[12:14], portalGroup)
	copy(data[14:], targetName)

	return data
}

func TestParseIScsiMessagingNode(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *IScsiMessagingNode
		wantErr bool
	}{
		{
			name: "TCP with CHAP UNI and null-terminated target name",
			data: newIScsiMessagingNodePayload(
				uint16(IScsiProtocolTCP),
				0x100a,
				[8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
				1,
				[]byte("iqn.2026-07.example:disk\x00"),
			),
			want: &IScsiMessagingNode{
				Protocol:          IScsiProtocolTCP,
				LoginOptions:      IScsiLoginOptions(0x100a),
				LogicalUnitNumber: [8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
				PortalGroup:       1,
				TargetName:        "iqn.2026-07.example:disk",
			},
		},
		{
			name: "target name without null terminator",
			data: newIScsiMessagingNodePayload(
				uint16(IScsiProtocolTCP),
				0x0800,
				[8]byte{0xde, 0xad, 0xbe, 0xef, 0x00, 0x00, 0x00, 0x01},
				42,
				[]byte("iqn.2026-07.example:disk"),
			),
			want: &IScsiMessagingNode{
				Protocol:          IScsiProtocolTCP,
				LoginOptions:      IScsiLoginOptions(0x0800),
				LogicalUnitNumber: [8]byte{0xde, 0xad, 0xbe, 0xef, 0x00, 0x00, 0x00, 0x01},
				PortalGroup:       42,
				TargetName:        "iqn.2026-07.example:disk",
			},
		},
		{
			name: "target name stops at first null byte",
			data: newIScsiMessagingNodePayload(
				uint16(IScsiProtocolTCP),
				0,
				[8]byte{},
				0,
				[]byte("target\x00ignored"),
			),
			want: &IScsiMessagingNode{
				Protocol:          IScsiProtocolTCP,
				LoginOptions:      0,
				LogicalUnitNumber: [8]byte{},
				PortalGroup:       0,
				TargetName:        "target",
			},
		},
		{
			name: "empty target name",
			data: newIScsiMessagingNodePayload(
				uint16(IScsiProtocolTCP),
				0,
				[8]byte{},
				0,
				nil,
			),
			want: &IScsiMessagingNode{
				Protocol:          IScsiProtocolTCP,
				LoginOptions:      0,
				LogicalUnitNumber: [8]byte{},
				PortalGroup:       0,
				TargetName:        "",
			},
		},
		{
			name:    "payload too short",
			data:    make([]byte, 13),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIScsiMessagingNode(tt.data)

			if tt.wantErr {
				if err == nil {
					t.Fatal("parseIScsiMessagingNode() error = nil, want an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("parseIScsiMessagingNode() unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"parseIScsiMessagingNode() mismatch:\ngot:  %#v\nwant: %#v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestIScsiMessagingNodeString(t *testing.T) {
	tests := []struct {
		name string
		node *IScsiMessagingNode
		want string
	}{
		{
			name: "CHAP BI without digests",
			node: &IScsiMessagingNode{
				Protocol:          IScsiProtocolTCP,
				LoginOptions:      0,
				LogicalUnitNumber: [8]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
				PortalGroup:       1,
				TargetName:        "iqn.2026-07.example:disk",
			},
			want: "iSCSI(iqn.2026-07.example:disk,1,0001020304050607,None,None,CHAP_BI,TCP)",
		},
		{
			name: "CHAP UNI with CRC32C digests",
			node: &IScsiMessagingNode{
				Protocol:          IScsiProtocolTCP,
				LoginOptions:      IScsiLoginOptions(0x100a),
				LogicalUnitNumber: [8]byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88},
				PortalGroup:       3260,
				TargetName:        "target",
			},
			want: "iSCSI(target,3260,ffeeddccbbaa9988,CRC32C,CRC32C,CHAP_UNI,TCP)",
		},
		{
			name: "no authentication",
			node: &IScsiMessagingNode{
				Protocol:          IScsiProtocolTCP,
				LoginOptions:      IScsiLoginOptions(0x0800),
				LogicalUnitNumber: [8]byte{},
				PortalGroup:       0,
				TargetName:        "",
			},
			want: "iSCSI(,0,0000000000000000,None,None,None,TCP)",
		},
		{
			name: "reserved values and unknown protocol",
			node: &IScsiMessagingNode{
				Protocol: IScsiProtocol(0x1234),
				LoginOptions: IScsiLoginOptions(
					1 | // Header digest: reserved
						(3 << 2) | // Data digest: reserved
						(1 << 10), // Authentication: reserved
				),
				LogicalUnitNumber: [8]byte{},
				PortalGroup:       7,
				TargetName:        "reserved",
			},
			want: "iSCSI(reserved,7,0000000000000000,Reserved(1),Reserved(3),Reserved(1),0x1234)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node.String()

			if got != tt.want {
				t.Errorf("IScsiMessagingNode.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIScsiLoginOptions(t *testing.T) {
	tests := []struct {
		name           string
		options        IScsiLoginOptions
		headerDigest   uint8
		dataDigest     uint8
		authentication uint8
		chapType       uint8
		headerString   string
		dataString     string
		authString     string
	}{
		{
			name:           "CHAP BI without digests",
			options:        0,
			headerDigest:   0,
			dataDigest:     0,
			authentication: 0,
			chapType:       0,
			headerString:   "None",
			dataString:     "None",
			authString:     "CHAP_BI",
		},
		{
			name:           "CHAP UNI with CRC32C",
			options:        0x100a,
			headerDigest:   2,
			dataDigest:     2,
			authentication: 0,
			chapType:       1,
			headerString:   "CRC32C",
			dataString:     "CRC32C",
			authString:     "CHAP_UNI",
		},
		{
			name:           "no authentication",
			options:        0x0800,
			headerDigest:   0,
			dataDigest:     0,
			authentication: 2,
			chapType:       0,
			headerString:   "None",
			dataString:     "None",
			authString:     "None",
		},
		{
			name:           "reserved values",
			options:        IScsiLoginOptions(1 | (3 << 2) | (1 << 10)),
			headerDigest:   1,
			dataDigest:     3,
			authentication: 1,
			chapType:       0,
			headerString:   "Reserved(1)",
			dataString:     "Reserved(3)",
			authString:     "Reserved(1)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.options.HeaderDigest(); got != tt.headerDigest {
				t.Errorf("HeaderDigest() = %d, want %d", got, tt.headerDigest)
			}

			if got := tt.options.DataDigest(); got != tt.dataDigest {
				t.Errorf("DataDigest() = %d, want %d", got, tt.dataDigest)
			}

			if got := tt.options.AuthenticationMethod(); got != tt.authentication {
				t.Errorf("AuthenticationMethod() = %d, want %d", got, tt.authentication)
			}

			if got := tt.options.ChapType(); got != tt.chapType {
				t.Errorf("ChapType() = %d, want %d", got, tt.chapType)
			}

			if got := tt.options.HeaderDigestString(); got != tt.headerString {
				t.Errorf("HeaderDigestString() = %q, want %q", got, tt.headerString)
			}

			if got := tt.options.DataDigestString(); got != tt.dataString {
				t.Errorf("DataDigestString() = %q, want %q", got, tt.dataString)
			}

			if got := tt.options.AuthenticationString(); got != tt.authString {
				t.Errorf("AuthenticationString() = %q, want %q", got, tt.authString)
			}
		})
	}
}
