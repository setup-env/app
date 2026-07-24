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
	Source NetworkSource
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
		snapshot.Networks = append(snapshot.Networks, NetworkInterface{
			Name:       item.Name,
			Status:     status,
			MACAddress: item.HardwareAddr.String(),
			Loopback:   loopback,
			Addresses:  filtered,
		})
	}
	sort.Slice(snapshot.Networks, func(left, right int) bool {
		return strings.ToLower(snapshot.Networks[left].Name) < strings.ToLower(snapshot.Networks[right].Name)
	})
	return errors.Join(problems...)
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
