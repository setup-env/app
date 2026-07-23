package manifest

const SchemaVersion = 1

type Manifest struct {
	SchemaVersion     int         `yaml:"schema_version" json:"schema_version"`
	ID                string      `yaml:"id" json:"id"`
	Name              string      `yaml:"name" json:"name"`
	Description       string      `yaml:"description" json:"description"`
	Repository        Repository  `yaml:"repository" json:"repository"`
	Version           VersionSpec `yaml:"version" json:"version"`
	Publisher         string      `yaml:"publisher" json:"publisher"`
	License           string      `yaml:"license" json:"license"`
	Homepage          string      `yaml:"homepage,omitempty" json:"homepage,omitempty"`
	Documentation     string      `yaml:"documentation,omitempty" json:"documentation,omitempty"`
	MinimumAppVersion string      `yaml:"minimum_app_version" json:"minimum_app_version"`
	Platforms         Platforms   `yaml:"platforms" json:"platforms"`
	Categories        []string    `yaml:"categories" json:"categories"`
	Tags              []string    `yaml:"tags,omitempty" json:"tags,omitempty"`
	Security          Security    `yaml:"security,omitempty" json:"security,omitempty"`
	Workflows         []Workflow  `yaml:"workflows" json:"workflows"`
	Deprecated        bool        `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Replacement       string      `yaml:"replacement,omitempty" json:"replacement,omitempty"`
	DeprecationNotice string      `yaml:"deprecation_notice,omitempty" json:"deprecation_notice,omitempty"`
}

type Repository struct {
	Owner     string `yaml:"owner" json:"owner"`
	Name      string `yaml:"name" json:"name"`
	IssuesURL string `yaml:"issues_url,omitempty" json:"issues_url,omitempty"`
}

type VersionSpec struct {
	Source string `yaml:"source" json:"source"`
	Value  string `yaml:"value,omitempty" json:"value,omitempty"`
}

type Platforms struct {
	OperatingSystems []string `yaml:"operating_systems" json:"operating_systems"`
	Architectures    []string `yaml:"architectures" json:"architectures"`
}

type Security struct {
	RequiresElevation bool     `yaml:"requires_elevation,omitempty" json:"requires_elevation,omitempty"`
	NetworkAccess     bool     `yaml:"network_access,omitempty" json:"network_access,omitempty"`
	SecretInputs      []string `yaml:"secret_inputs,omitempty" json:"secret_inputs,omitempty"`
}

type Workflow struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Entrypoint  string `yaml:"entrypoint" json:"entrypoint"`
}
