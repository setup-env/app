package dashboard

import (
	"testing"
	"time"

	"github.com/setup-env/app/internal/system"
)

func TestCalculateNetworkRates(t *testing.T) {
	oldReceived, oldSent := uint64(100), uint64(200)
	newReceived, newSent := uint64(1124), uint64(2248)
	rates := CalculateNetworkRates(
		[]system.NetworkInterface{{Name: "eth0", BytesReceived: &oldReceived, BytesTransmitted: &oldSent}},
		[]system.NetworkInterface{{Name: "eth0", BytesReceived: &newReceived, BytesTransmitted: &newSent}},
		2*time.Second,
	)
	if len(rates) != 1 || *rates[0].BytesReceivedPerSec != 512 || *rates[0].BytesSentPerSec != 1024 {
		t.Fatalf("rates = %#v", rates)
	}
}

func TestCalculateNetworkRatesHandlesResetAndInterfaceChurn(t *testing.T) {
	oldValue, resetValue, newValue := uint64(100), uint64(10), uint64(20)
	rates := CalculateNetworkRates(
		[]system.NetworkInterface{
			{Name: "gone", BytesReceived: &oldValue},
			{Name: "reset", BytesReceived: &oldValue},
		},
		[]system.NetworkInterface{
			{Name: "new", BytesReceived: &newValue},
			{Name: "reset", BytesReceived: &resetValue},
		},
		time.Second,
	)
	if len(rates) != 2 || rates[0].BytesReceivedPerSec != nil || rates[1].BytesReceivedPerSec != nil {
		t.Fatalf("rates = %#v", rates)
	}
}

func TestCalculateNetworkRatesRejectsZeroElapsedTime(t *testing.T) {
	if rates := CalculateNetworkRates(nil, nil, 0); rates != nil {
		t.Fatalf("rates = %#v, want nil", rates)
	}
}
