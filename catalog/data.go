// Package catalogdata exposes the repository catalog as embedded application data.
package catalogdata

import _ "embed"

// Modules contains the authoritative Milestone 02 catalog.
//
//go:embed modules.yaml
var Modules []byte
