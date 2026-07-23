package compatibility

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type State string

const (
	StateCompatible   State = "compatible"
	StateIncompatible State = "incompatible"
	StateUnknown      State = "unknown"
)

type Result struct {
	State             State  `json:"state"`
	AppVersion        string `json:"app_version"`
	MinimumAppVersion string `json:"minimum_app_version,omitempty"`
	Reason            string `json:"reason"`
}

type Version struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease []string
}

var versionPattern = regexp.MustCompile(`^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?$`)

func Parse(value string) (Version, error) {
	matches := versionPattern.FindStringSubmatch(value)
	if matches == nil {
		return Version{}, fmt.Errorf("%q is not a semantic version in MAJOR.MINOR.PATCH form", value)
	}
	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])
	var prerelease []string
	if matches[4] != "" {
		prerelease = strings.Split(matches[4], ".")
		for _, identifier := range prerelease {
			if isNumeric(identifier) && len(identifier) > 1 && identifier[0] == '0' {
				return Version{}, fmt.Errorf("%q has a numeric prerelease identifier with a leading zero", value)
			}
		}
	}
	return Version{Major: major, Minor: minor, Patch: patch, Prerelease: prerelease}, nil
}

func Evaluate(appVersion, minimumAppVersion string) Result {
	result := Result{
		State:             StateUnknown,
		AppVersion:        appVersion,
		MinimumAppVersion: minimumAppVersion,
	}
	if minimumAppVersion == "" {
		result.Reason = "module does not declare a minimum application version"
		return result
	}
	minimum, err := Parse(minimumAppVersion)
	if err != nil {
		result.Reason = "module minimum application version is malformed: " + err.Error()
		return result
	}
	app, err := Parse(appVersion)
	if err != nil {
		result.Reason = "application version is not a released semantic version"
		return result
	}
	if compare(app, minimum) >= 0 {
		result.State = StateCompatible
		result.Reason = "application version meets the module minimum"
		return result
	}
	result.State = StateIncompatible
	result.Reason = "application version is lower than the module minimum"
	return result
}

func compare(left, right Version) int {
	for _, values := range [][2]int{{left.Major, right.Major}, {left.Minor, right.Minor}, {left.Patch, right.Patch}} {
		if values[0] < values[1] {
			return -1
		}
		if values[0] > values[1] {
			return 1
		}
	}
	return comparePrerelease(left.Prerelease, right.Prerelease)
}

func comparePrerelease(left, right []string) int {
	if len(left) == 0 && len(right) == 0 {
		return 0
	}
	if len(left) == 0 {
		return 1
	}
	if len(right) == 0 {
		return -1
	}
	for index := 0; index < len(left) && index < len(right); index++ {
		if left[index] == right[index] {
			continue
		}
		leftNumeric := isNumeric(left[index])
		rightNumeric := isNumeric(right[index])
		if leftNumeric && rightNumeric {
			leftNumber, _ := strconv.Atoi(left[index])
			rightNumber, _ := strconv.Atoi(right[index])
			if leftNumber < rightNumber {
				return -1
			}
			return 1
		}
		if leftNumeric {
			return -1
		}
		if rightNumeric {
			return 1
		}
		if left[index] < right[index] {
			return -1
		}
		return 1
	}
	if len(left) < len(right) {
		return -1
	}
	if len(left) > len(right) {
		return 1
	}
	return 0
}

func isNumeric(value string) bool {
	if value == "" {
		return false
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}
