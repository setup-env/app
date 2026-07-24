package system

import (
	"time"

	"github.com/setup-env/app/internal/diagnostics"
	"github.com/setup-env/app/internal/directory"
	"github.com/setup-env/app/internal/version"
)

const SnapshotSchemaVersion = 1

type Health string

const (
	HealthHealthy   Health = "healthy"
	HealthWarning   Health = "warning"
	HealthUnhealthy Health = "unhealthy"
)

type Snapshot struct {
	SchemaVersion   int                `json:"schema_version"`
	Timestamp       time.Time          `json:"timestamp"`
	TimeZone        TimeZone           `json:"time_zone"`
	Host            Host               `json:"host"`
	OperatingSystem OperatingSystem    `json:"operating_system"`
	User            User               `json:"user"`
	CPU             CPU                `json:"cpu"`
	Memory          Memory             `json:"memory"`
	Filesystems     []Filesystem       `json:"filesystems"`
	Networks        []NetworkInterface `json:"networks"`
	Development     Development        `json:"development"`
	Diagnostics     DiagnosticSummary  `json:"diagnostics"`
	Warnings        []Warning          `json:"warnings"`
}

type TimeZone struct {
	Name             string `json:"name"`
	UTCOffsetSeconds int    `json:"utc_offset_seconds"`
}

type Host struct {
	Hostname      string  `json:"hostname,omitempty"`
	UptimeSeconds *uint64 `json:"uptime_seconds"`
}

type OperatingSystem struct {
	OS                  string `json:"os"`
	DisplayName         string `json:"display_name,omitempty"`
	Version             string `json:"version,omitempty"`
	BuildInformation    string `json:"build_information,omitempty"`
	Distribution        string `json:"distribution,omitempty"`
	DistributionVersion string `json:"distribution_version,omitempty"`
	KernelVersion       string `json:"kernel_version,omitempty"`
	Architecture        string `json:"architecture"`
}

type User struct {
	Username string `json:"username,omitempty"`
	Home     string `json:"home,omitempty"`
	Shell    string `json:"shell,omitempty"`
}

type CPU struct {
	Model                string   `json:"model,omitempty"`
	PhysicalCores        *int     `json:"physical_cores"`
	LogicalCPUs          *int     `json:"logical_cpus"`
	UtilizationPercent   *float64 `json:"utilization_percent"`
	SampleDurationMillis int64    `json:"sample_duration_millis"`
}

type Memory struct {
	TotalBytes         *uint64  `json:"total_bytes"`
	AvailableBytes     *uint64  `json:"available_bytes"`
	UsedBytes          *uint64  `json:"used_bytes"`
	UtilizationPercent *float64 `json:"utilization_percent"`
}

type Filesystem struct {
	Device             string   `json:"device,omitempty"`
	Mountpoint         string   `json:"mountpoint"`
	Type               string   `json:"type,omitempty"`
	TotalBytes         *uint64  `json:"total_bytes"`
	AvailableBytes     *uint64  `json:"available_bytes"`
	UsedBytes          *uint64  `json:"used_bytes"`
	UtilizationPercent *float64 `json:"utilization_percent"`
}

type NetworkInterface struct {
	Name               string           `json:"name"`
	Status             string           `json:"status"`
	MACAddress         string           `json:"mac_address,omitempty"`
	Loopback           bool             `json:"loopback"`
	Addresses          []NetworkAddress `json:"addresses"`
	BytesReceived      *uint64          `json:"bytes_received"`
	BytesTransmitted   *uint64          `json:"bytes_transmitted"`
	PacketsReceived    *uint64          `json:"packets_received"`
	PacketsTransmitted *uint64          `json:"packets_transmitted"`
}

type NetworkAddress struct {
	Address      string `json:"address"`
	Family       string `json:"family"`
	PrefixLength int    `json:"prefix_length"`
}

type Development struct {
	ApplicationVersion  version.Info                      `json:"application_version"`
	Root                diagnostics.DevelopmentRootStatus `json:"root"`
	Directory           directory.Context                 `json:"directory"`
	Git                 diagnostics.ToolStatus            `json:"git"`
	GitHubCLI           diagnostics.ToolStatus            `json:"github_cli"`
	GitHubAuthenticated *bool                             `json:"github_authenticated"`
}

type DiagnosticSummary struct {
	Health       Health              `json:"health"`
	WarningCount int                 `json:"warning_count"`
	FailureCount int                 `json:"failure_count"`
	Checks       []diagnostics.Check `json:"checks"`
}

type Warning struct {
	Section string `json:"section"`
	Message string `json:"message"`
}
