package catalog

const SchemaVersion = 1

type Trust string

const (
	TrustOfficial  Trust = "official"
	TrustVerified  Trust = "verified"
	TrustCommunity Trust = "community"
)

type Status string

const (
	StatusActive       Status = "active"
	StatusPlanned      Status = "planned"
	StatusExperimental Status = "experimental"
	StatusDeprecated   Status = "deprecated"
	StatusUnavailable  Status = "unavailable"
)

type Catalog struct {
	SchemaVersion int     `yaml:"schema_version" json:"schema_version"`
	Modules       []Entry `yaml:"modules" json:"modules"`
}

type Entry struct {
	ID            string   `yaml:"id" json:"id"`
	Name          string   `yaml:"name" json:"name"`
	Description   string   `yaml:"description" json:"description"`
	Repository    string   `yaml:"repository" json:"repository"`
	Manifest      string   `yaml:"manifest" json:"manifest"`
	Trust         Trust    `yaml:"trust" json:"trust"`
	Status        Status   `yaml:"status" json:"status"`
	Categories    []string `yaml:"categories" json:"categories"`
	Tags          []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	PinnedVersion string   `yaml:"pinned_version,omitempty" json:"pinned_version,omitempty"`
	VersionPolicy string   `yaml:"version_policy,omitempty" json:"version_policy,omitempty"`
	Publisher     string   `yaml:"publisher,omitempty" json:"publisher,omitempty"`
}

func (e Entry) RepositoryParts() (owner, name string) {
	for index, character := range e.Repository {
		if character == '/' {
			return e.Repository[:index], e.Repository[index+1:]
		}
	}
	return "", ""
}
