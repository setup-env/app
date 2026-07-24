package system

import (
	"net"
	"testing"
)

type stringAddress string

func (s stringAddress) Network() string { return "test" }
func (s stringAddress) String() string  { return string(s) }

func TestFilterNetworkAddresses(t *testing.T) {
	input := []net.Addr{
		stringAddress("fe80::1/64"),
		stringAddress("10.0.0.9/24"),
		stringAddress("0.0.0.0/0"),
		stringAddress("224.0.0.1/32"),
		stringAddress("10.0.0.9/24"),
	}
	result := FilterNetworkAddresses(input)
	if len(result) != 2 {
		t.Fatalf("addresses = %#v", result)
	}
	if result[0].Family != "ipv4" || result[0].Address != "10.0.0.9" {
		t.Fatalf("addresses are not deterministic: %#v", result)
	}
	if result[1].Family != "ipv6" || result[1].Address != "fe80::1" {
		t.Fatalf("addresses are not deterministic: %#v", result)
	}
}
