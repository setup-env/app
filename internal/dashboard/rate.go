package dashboard

import (
	"sort"
	"strings"
	"time"

	"github.com/setup-env/app/internal/system"
)

type NetworkRate struct {
	Name                string
	BytesReceivedPerSec *float64
	BytesSentPerSec     *float64
}

func CalculateNetworkRates(previous, current []system.NetworkInterface, elapsed time.Duration) []NetworkRate {
	if elapsed <= 0 {
		return nil
	}
	previousByName := make(map[string]system.NetworkInterface, len(previous))
	for _, network := range previous {
		previousByName[strings.ToLower(network.Name)] = network
	}
	seconds := elapsed.Seconds()
	rates := make([]NetworkRate, 0, len(current))
	for _, network := range current {
		rate := NetworkRate{Name: network.Name}
		old, exists := previousByName[strings.ToLower(network.Name)]
		if exists {
			rate.BytesReceivedPerSec = counterRate(old.BytesReceived, network.BytesReceived, seconds)
			rate.BytesSentPerSec = counterRate(old.BytesTransmitted, network.BytesTransmitted, seconds)
		}
		rates = append(rates, rate)
	}
	sort.Slice(rates, func(left, right int) bool {
		return strings.ToLower(rates[left].Name) < strings.ToLower(rates[right].Name)
	})
	return rates
}

func counterRate(previous, current *uint64, seconds float64) *float64 {
	if previous == nil || current == nil || seconds <= 0 || *current < *previous {
		return nil
	}
	value := float64(*current-*previous) / seconds
	return &value
}
