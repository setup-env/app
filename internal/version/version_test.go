package version

import "testing"

func TestCurrentIncludesInjectedMetadata(t *testing.T) {
	originalVersion, originalCommit, originalDate, originalDirty := Version, Commit, Date, Dirty
	t.Cleanup(func() {
		Version, Commit, Date, Dirty = originalVersion, originalCommit, originalDate, originalDirty
	})
	Version = "v0.1.0"
	Commit = "abcdef123456"
	Date = "2026-07-24T12:00:00Z"
	Dirty = "false"

	got := Current()
	if got.Version != Version || got.Commit != Commit || got.BuildDate != Date || got.Dirty {
		t.Fatalf("Current() = %+v", got)
	}
}

func TestCurrentTreatsInvalidDirtyValueAsDirty(t *testing.T) {
	original := Dirty
	t.Cleanup(func() { Dirty = original })
	Dirty = "not-a-boolean"
	if !Current().Dirty {
		t.Fatal("Current().Dirty = false, want safe dirty default")
	}
}
