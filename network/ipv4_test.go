package network

import (
	"fmt"
	"testing"
)

func TestParseIPv4Address(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    IPv4Address
		wantErr bool
	}{
		{
			name: "valid",
			data: []byte{192, 168, 1, 42},
			want: IPv4Address{192, 168, 1, 42},
		},
		{
			name:    "empty",
			data:    nil,
			wantErr: true,
		},
		{
			name:    "too short",
			data:    []byte{192, 168, 1},
			wantErr: true,
		},
		{
			name:    "too long",
			data:    []byte{192, 168, 1, 42, 1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIPv4Address(tt.data)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseIPv4Address(%v) returned no error", tt.data)
				}

				return
			}

			if err != nil {
				t.Fatalf("ParseIPv4Address(%v) returned error: %v", tt.data, err)
			}

			if got != tt.want {
				t.Errorf(
					"ParseIPv4Address(%v) = %#v, want %#v",
					tt.data,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestIPv4AddressString(t *testing.T) {
	address := IPv4Address{192, 168, 1, 42}

	if got, want := address.String(), "192.168.1.42"; got != want {
		t.Errorf("IPv4Address.String() = %q, want %q", got, want)
	}
}

func TestIPv4AddressGoString(t *testing.T) {
	address := IPv4Address{192, 168, 1, 42}

	if got, want := address.GoString(),
		`network.MustParseIPv4Address("192.168.1.42")`; got != want {
		t.Errorf("IPv4Address.GoString() = %q, want %q", got, want)
	}
}

func TestIPv4AddressGoStringWithFmt(t *testing.T) {
	address := IPv4Address{10, 20, 30, 40}

	if got, want := fmt.Sprintf("%#v", address),
		`network.MustParseIPv4Address("10.20.30.40")`; got != want {
		t.Errorf("fmt.Sprintf(\"%%#v\", address) = %q, want %q", got, want)
	}
}

func TestMustParseIPv4Address(t *testing.T) {
	address := MustParseIPv4Address("192.0.2.10")
	want := IPv4Address{192, 0, 2, 10}

	if address != want {
		t.Errorf("MustParseIPv4Address() = %#v, want %#v", address, want)
	}
}

func TestMustParseIPv4AddressPanics(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "invalid address",
			value: "not-an-address",
		},
		{
			name:  "IPv6 address",
			value: "2001:db8::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf(
						"MustParseIPv4Address(%q) did not panic",
						tt.value,
					)
				}
			}()

			_ = MustParseIPv4Address(tt.value)
		})
	}
}
