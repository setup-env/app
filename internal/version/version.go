package version

import (
	"runtime"
	"strconv"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	Dirty   = "true"
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	Dirty     bool   `json:"dirty"`
	GoVersion string `json:"go_version"`
}

func Current() Info {
	dirty, err := strconv.ParseBool(Dirty)
	if err != nil {
		dirty = true
	}
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: Date,
		Dirty:     dirty,
		GoVersion: runtime.Version(),
	}
}
