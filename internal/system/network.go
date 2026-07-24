package system

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
)

type NetworkCollector struct {
	Source        NetworkSource
	CounterSource NetworkCounterSource
}

func (NetworkCollector) Name() string { return "network" }

func (c NetworkCollector) Collect(ctx context.Context, snapshot *Snapshot) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	source := c.Source
	if source == nil {
		source = StandardNetworkSource{}
	}
	interfaces, err := source.Interfaces()
	if err != nil {
		return fmt.Errorf("list interfaces: %w", err)
	}
	var problems []error
	counters := make(map[string]networkCounters)
	if c.CounterSource != nil {
		values, counterErr := c.CounterSource.IOCounters(ctx, true)
		if counterErr != nil {
			problems = append(problems, fmt.Errorf("interface counters: %w", counterErr))
		} else {
			for _, value := range values {
				counters[strings.ToLower(value.Name)] = networkCounters{
					bytesReceived:      value.BytesRecv,
					bytesTransmitted:   value.BytesSent,
					packetsReceived:    value.PacketsRecv,
					packetsTransmitted: value.PacketsSent,
				}
			}
		}
	}
	for _, item := range interfaces {
		if err := ctx.Err(); err != nil {
			problems = append(problems, err)
			break
		}
		up := item.Flags&net.FlagUp != 0
		loopback := item.Flags&net.FlagLoopback != 0
		if !up && !loopback {
			continue
		}
		addresses, err := source.Addrs(item)
		if err != nil {
			problems = append(problems, fmt.Errorf("%s: %w", item.Name, err))
			continue
		}
		filtered := FilterNetworkAddresses(addresses)
		if len(filtered) == 0 && !loopback {
			continue
		}
		status := "down"
		if up {
			status = "up"
		}
		network := NetworkInterface{
			Name:       item.Name,
			Status:     status,
			MACAddress: item.HardwareAddr.String(),
			Loopback:   loopback,
			Addresses:  filtered,
		}
		if value, ok := counters[strings.ToLower(item.Name)]; ok {
			network.BytesReceived = uint64Pointer(value.bytesReceived)
			network.BytesTransmitted = uint64Pointer(value.bytesTransmitted)
			network.PacketsReceived = uint64Pointer(value.packetsReceived)
			network.PacketsTransmitted = uint64Pointer(value.packetsTransmitted)
		}
		snapshot.Networks = append(snapshot.Networks, network)
	}
	sort.Slice(snapshot.Networks, func(left, right int) bool {
		return strings.ToLower(snapshot.Networks[left].Name) < strings.ToLower(snapshot.Networks[right].Name)
	})
	return errors.Join(problems...)
}

type networkCounters struct {
	bytesReceived      uint64
	bytesTransmitted   uint64
	packetsReceived    uint64
	packetsTransmitted uint64
}

func FilterNetworkAddresses(addresses []net.Addr) []NetworkAddress {
	result := make([]NetworkAddress, 0, len(addresses))
	seen := make(map[string]struct{}, len(addresses))
	for _, address := range addresses {
		ip, network, err := net.ParseCIDR(address.String())
		if err != nil || ip == nil || ip.IsUnspecified() || ip.IsMulticast() {
			continue
		}
		ones, _ := network.Mask.Size()
		family := "ipv6"
		if ip.To4() != nil {
			ip = ip.To4()
			family = "ipv4"
		}
		key := family + "/" + ip.String()
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, NetworkAddress{
			Address:      ip.String(),
			Family:       family,
			PrefixLength: ones,
		})
	}
	sort.Slice(result, func(left, right int) bool {
		if result[left].Family == result[right].Family {
			return result[left].Address < result[right].Address
		}
		return result[left].Family < result[right].Family
	})
	return result
}
