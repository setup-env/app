package dashboard

import "testing"

func TestHistoryIsBoundedAndOldestToNewest(t *testing.T) {
	history := NewHistory(3)
	for _, value := range []float64{10, 20, 30, 40} {
		history = history.Add(value)
	}
	got := history.Values()
	want := []float64{20, 30, 40}
	if len(got) != len(want) {
		t.Fatalf("values = %#v", got)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("values = %#v, want %#v", got, want)
		}
	}
	got[0] = 999
	if history.Values()[0] != 20 {
		t.Fatal("Values returned mutable history storage")
	}
}
