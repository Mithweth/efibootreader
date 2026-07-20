package devicepath

import (
	"fmt"
	"github.com/Mithweth/efibootreader/network"
	"strings"
	"testing"
)

func TestParseDnsMessagingNodeIPv4(t *testing.T) {
	data := []byte{
		0, // IPv4
		1, 1, 1, 1,
		8, 8, 8, 8,
	}

	got, err := parseDnsMessagingNode(data)
	if err != nil {
		t.Fatalf("parseDnsMessagingNode() returned error: %v", err)
	}

	if got.IsIPv6 {
		t.Error("parseDnsMessagingNode().IsIPv6 = true, want false")
	}

	want := []network.IPv4Address{
		network.MustParseIPv4Address("1.1.1.1"),
		network.MustParseIPv4Address("8.8.8.8"),
	}

	if len(got.IPv4Addresses) != len(want) {
		t.Fatalf(
			"len(IPv4Addresses) = %d, want %d",
			len(got.IPv4Addresses),
			len(want),
		)
	}

	for i := range want {
		if got.IPv4Addresses[i] != want[i] {
			t.Errorf(
				"IPv4Addresses[%d] = %#v, want %#v",
				i,
				got.IPv4Addresses[i],
				want[i],
			)
		}
	}

	if len(got.IPv6Addresses) != 0 {
		t.Errorf(
			"len(IPv6Addresses) = %d, want 0",
			len(got.IPv6Addresses),
		)
	}
}

func TestParseDnsMessagingNodeIPv6(t *testing.T) {
	address1 := network.MustParseIPv6Address("2001:4860:4860::8888")
	address2 := network.MustParseIPv6Address("2606:4700:4700::1111")

	data := []byte{1}
	data = append(data, address1[:]...)
	data = append(data, address2[:]...)

	got, err := parseDnsMessagingNode(data)
	if err != nil {
		t.Fatalf("parseDnsMessagingNode() returned error: %v", err)
	}

	if !got.IsIPv6 {
		t.Error("parseDnsMessagingNode().IsIPv6 = false, want true")
	}

	want := []network.IPv6Address{address1, address2}

	if len(got.IPv6Addresses) != len(want) {
		t.Fatalf(
			"len(IPv6Addresses) = %d, want %d",
			len(got.IPv6Addresses),
			len(want),
		)
	}

	for i := range want {
		if got.IPv6Addresses[i] != want[i] {
			t.Errorf(
				"IPv6Addresses[%d] = %#v, want %#v",
				i,
				got.IPv6Addresses[i],
				want[i],
			)
		}
	}

	if len(got.IPv4Addresses) != 0 {
		t.Errorf(
			"len(IPv4Addresses) = %d, want 0",
			len(got.IPv4Addresses),
		)
	}
}

func TestParseDnsMessagingNodeErrors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "empty payload",
			data: nil,
		},
		{
			name: "invalid address type",
			data: []byte{2},
		},
		{
			name: "incomplete IPv4 address",
			data: []byte{
				0,
				192, 168, 1,
			},
		},
		{
			name: "incomplete second IPv4 address",
			data: []byte{
				0,
				192, 168, 1, 1,
				8, 8,
			},
		},
		{
			name: "incomplete IPv6 address",
			data: append(
				[]byte{1},
				make([]byte, 15)...,
			),
		},
		{
			name: "incomplete second IPv6 address",
			data: append(
				append(
					[]byte{1},
					make([]byte, 16)...,
				),
				make([]byte, 8)...,
			),
		},
		{
			name: "No DNS server ipv4",
			data: []byte{0},
		},
		{
			name: "No DNS server ipv6",
			data: []byte{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDnsMessagingNode(tt.data)

			if err == nil {
				t.Fatalf(
					"parseDnsMessagingNode(%v) = %#v, want error",
					tt.data,
					got,
				)
			}
		})
	}
}

func TestDnsMessagingNodeStringIPv4(t *testing.T) {
	node := &DnsMessagingNode{
		IPv4Addresses: []network.IPv4Address{
			network.MustParseIPv4Address("1.1.1.1"),
			network.MustParseIPv4Address("8.8.8.8"),
		},
	}

	if got, want := node.String(), "Dns(1.1.1.1,8.8.8.8)"; got != want {
		t.Errorf("DnsMessagingNode.String() = %q, want %q", got, want)
	}
}

func TestDnsMessagingNodeStringIPv6(t *testing.T) {
	node := &DnsMessagingNode{
		IsIPv6: true,
		IPv6Addresses: []network.IPv6Address{
			network.MustParseIPv6Address("2001:4860:4860::8888"),
			network.MustParseIPv6Address("2606:4700:4700::1111"),
		},
	}

	want := "Dns(2001:4860:4860::8888,2606:4700:4700::1111)"

	if got := node.String(); got != want {
		t.Errorf("DnsMessagingNode.String() = %q, want %q", got, want)
	}
}

func TestDnsMessagingNodeGoString(t *testing.T) {
	node := &DnsMessagingNode{
		IsIPv6: true,
		IPv6Addresses: []network.IPv6Address{
			network.MustParseIPv6Address("2001:db8::1"),
		},
	}

	got := node.GoString()

	expectedParts := []string{
		"&devicepath.DnsMessagingNode{",
		"IsIPv6:true",
		`network.MustParseIPv6Address("2001:db8::1")`,
	}

	for _, expected := range expectedParts {
		if !strings.Contains(got, expected) {
			t.Errorf(
				"DnsMessagingNode.GoString() = %q, want it to contain %q",
				got,
				expected,
			)
		}
	}
}

func TestDnsMessagingNodeNilGoString(t *testing.T) {
	var node *DnsMessagingNode

	if got, want := node.GoString(),
		"(*devicepath.DnsMessagingNode)(nil)"; got != want {
		t.Errorf("DnsMessagingNode.GoString() = %q, want %q", got, want)
	}
}

func TestDnsMessagingNodeDumpIPv4(t *testing.T) {
	node := &DnsMessagingNode{
		IPv4Addresses: []network.IPv4Address{
			network.MustParseIPv4Address("1.1.1.1"),
			network.MustParseIPv4Address("8.8.8.8"),
		},
	}

	var output strings.Builder
	node.dump(&output, "  ")

	want := "" +
		"  DNS Messaging Node\n" +
		"    Address Type\t : IPv4\n" +
		"    DNS Server\t : 1.1.1.1\n" +
		"    DNS Server\t : 8.8.8.8\n"

	if got := output.String(); got != want {
		t.Errorf(
			"DnsMessagingNode.dump() =\n%s\nwant:\n%s",
			got,
			want,
		)
	}
}

func TestDnsMessagingNodeGoStringWithFmt(t *testing.T) {
	node := &DnsMessagingNode{
		IPv4Addresses: []network.IPv4Address{
			network.MustParseIPv4Address("192.0.2.53"),
		},
	}

	got := fmt.Sprintf("%#v", node)

	if !strings.Contains(
		got,
		`network.MustParseIPv4Address("192.0.2.53")`,
	) {
		t.Errorf(
			"fmt.Sprintf(\"%%#v\", node) = %q, want parsed IPv4 address",
			got,
		)
	}
}
