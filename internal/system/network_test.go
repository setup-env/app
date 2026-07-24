package system

import (
	"context"
	"net"
	"testing"

	gopsutilnet "github.com/shirou/gopsutil/v4/net"
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

type fakeNetworkSource struct {
	interfaces []net.Interface
	addresses  map[string][]net.Addr
}

func (f fakeNetworkSource) Interfaces() ([]net.Interface, error) {
	return f.interfaces, nil
}

func (f fakeNetworkSource) Addrs(value net.Interface) ([]net.Addr, error) {
	return f.addresses[value.Name], nil
}

type fakeNetworkCounterSource struct {
	values []gopsutilnet.IOCountersStat
}

func (f fakeNetworkCounterSource) IOCounters(context.Context, bool) ([]gopsutilnet.IOCountersStat, error) {
	return f.values, nil
}

func TestNetworkCollectorAddsMatchingCounters(t *testing.T) {
	source := fakeNetworkSource{
		interfaces: []net.Interface{{Name: "Ethernet", Flags: net.FlagUp}},
		addresses: map[string][]net.Addr{
			"Ethernet": {stringAddress("192.0.2.4/24")},
		},
	}
	counters := fakeNetworkCounterSource{values: []gopsutilnet.IOCountersStat{{
		Name:        "ethernet",
		BytesRecv:   100,
		BytesSent:   200,
		PacketsRecv: 3,
		PacketsSent: 4,
	}}}
	snapshot := Snapshot{}
	err := (NetworkCollector{Source: source, CounterSource: counters}).Collect(context.Background(), &snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Networks) != 1 ||
		*snapshot.Networks[0].BytesReceived != 100 ||
		*snapshot.Networks[0].BytesTransmitted != 200 ||
		*snapshot.Networks[0].PacketsReceived != 3 ||
		*snapshot.Networks[0].PacketsTransmitted != 4 {
		t.Fatalf("network = %#v", snapshot.Networks)
	}
}
