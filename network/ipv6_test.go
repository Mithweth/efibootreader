package network

import (
	"fmt"
	"testing"
)

func TestParseIPv6Address(t *testing.T) {
	valid := []byte{
		0x20, 0x01, 0x0d, 0xb8,
		0x00, 0x00, 0x85, 0xa3,
		0x00, 0x00, 0x00, 0x00,
		0xac, 0x1f, 0x80, 0x01,
	}

	tests := []struct {
		name    string
		data    []byte
		want    IPv6Address
		wantErr bool
	}{
		{
			name: "valid",
			data: valid,
			want: IPv6Address{
				0x20, 0x01, 0x0d, 0xb8,
				0x00, 0x00, 0x85, 0xa3,
				0x00, 0x00, 0x00, 0x00,
				0xac, 0x1f, 0x80, 0x01,
			},
		},
		{
			name:    "empty",
			data:    nil,
			wantErr: true,
		},
		{
			name:    "too short",
			data:    make([]byte, 15),
			wantErr: true,
		},
		{
			name:    "too long",
			data:    make([]byte, 17),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIPv6Address(tt.data)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseIPv6Address() returned no error")
				}

				return
			}

			if err != nil {
				t.Fatalf("ParseIPv6Address() returned error: %v", err)
			}

			if got != tt.want {
				t.Errorf(
					"ParseIPv6Address() = %#v, want %#v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestIPv6AddressString(t *testing.T) {
	address := IPv6Address{
		0x20, 0x01, 0x0d, 0xb8,
		0x00, 0x00, 0x85, 0xa3,
		0x00, 0x00, 0x00, 0x00,
		0xac, 0x1f, 0x80, 0x01,
	}

	if got, want := address.String(), "2001:db8:0:85a3::ac1f:8001"; got != want {
		t.Errorf("IPv6Address.String() = %q, want %q", got, want)
	}
}

func TestIPv6AddressGoString(t *testing.T) {
	address := IPv6Address{
		0x20, 0x01, 0x0d, 0xb8,
		0x00, 0x00, 0x85, 0xa3,
		0x00, 0x00, 0x00, 0x00,
		0xac, 0x1f, 0x80, 0x01,
	}

	want := `network.MustParseIPv6Address("2001:db8:0:85a3::ac1f:8001")`

	if got := address.GoString(); got != want {
		t.Errorf("IPv6Address.GoString() = %q, want %q", got, want)
	}
}

func TestIPv6AddressGoStringWithFmt(t *testing.T) {
	address := MustParseIPv6Address("2001:db8::42")
	want := `network.MustParseIPv6Address("2001:db8::42")`

	if got := fmt.Sprintf("%#v", address); got != want {
		t.Errorf("fmt.Sprintf(\"%%#v\", address) = %q, want %q", got, want)
	}
}

func TestMustParseIPv6Address(t *testing.T) {
	address := MustParseIPv6Address("2001:db8::1")

	want := IPv6Address{
		0x20, 0x01, 0x0d, 0xb8,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01,
	}

	if address != want {
		t.Errorf("MustParseIPv6Address() = %#v, want %#v", address, want)
	}
}

func TestMustParseIPv6AddressPanics(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "invalid address",
			value: "not-an-address",
		},
		{
			name:  "IPv4 address",
			value: "192.0.2.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf(
						"MustParseIPv6Address(%q) did not panic",
						tt.value,
					)
				}
			}()

			_ = MustParseIPv6Address(tt.value)
		})
	}
}
