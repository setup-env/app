package catalog

import (
	"context"
	"strings"
	"testing"
)

func validCatalog(t *testing.T) Catalog {
	t.Helper()
	value, err := Load(context.Background(), EmbeddedSource{})
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func TestEmbeddedCatalogLoads(t *testing.T) {
	value := validCatalog(t)
	if len(value.Modules) != 10 {
		t.Fatalf("module count = %d, want 10", len(value.Modules))
	}
	for _, entry := range value.Modules {
		if entry.Trust != TrustOfficial || entry.Status != StatusPlanned {
			t.Fatalf("entry %#v is not honestly official/planned", entry)
		}
	}
}

func TestDuplicateModuleID(t *testing.T) {
	value := validCatalog(t)
	value.Modules[1].ID = value.Modules[0].ID
	assertCatalogProblem(t, Validate(value), "module id ansible is duplicated")
}

func TestDuplicateRepository(t *testing.T) {
	value := validCatalog(t)
	value.Modules[1].Repository = value.Modules[0].Repository
	assertCatalogProblem(t, Validate(value), "repository setup-env/ansible is duplicated")
}

func TestInvalidTrustAndStatus(t *testing.T) {
	value := validCatalog(t)
	value.Modules[0].Trust = "self-verified"
	value.Modules[0].Status = "installable"
	err := Validate(value)
	assertCatalogProblem(t, err, "trust is not recognized")
	assertCatalogProblem(t, err, "status is not recognized")
}

func TestOfficialTrustRequiresSetupEnvOwner(t *testing.T) {
	value := validCatalog(t)
	value.Modules[0].Repository = "someone/ansible"
	assertCatalogProblem(t, Validate(value), "official trust is restricted")
}

func TestAppAndAwesomeAreForbidden(t *testing.T) {
	for _, repository := range []string{"setup-env/app", "setup-env/awesome-setup-env"} {
		value := validCatalog(t)
		value.Modules[0].Repository = repository
		assertCatalogProblem(t, Validate(value), "must not be listed as a module")
	}
}

func TestCatalogOrderingIsDeterministic(t *testing.T) {
	value := validCatalog(t)
	value.Modules[0], value.Modules[1] = value.Modules[1], value.Modules[0]
	assertCatalogProblem(t, Validate(value), "modules must be sorted by id")
}

func TestFilter(t *testing.T) {
	value := validCatalog(t)
	value.Modules[0].Status = StatusActive
	result := Filter(value.Modules, TrustOfficial, StatusActive, "configuration-management")
	if len(result) != 1 || result[0].ID != "ansible" {
		t.Fatalf("Filter() = %#v", result)
	}
	if result := Filter(value.Modules, "", StatusPlanned, "cloud"); len(result) != 3 {
		t.Fatalf("planned cloud count = %d, want 3", len(result))
	}
}

func TestParseRejectsUnknownFields(t *testing.T) {
	data := []byte("schema_version: 1\nunknown: true\nmodules: []\n")
	if _, err := Parse(data); err == nil {
		t.Fatal("Parse() error = nil, want error")
	}
}

func assertCatalogProblem(t *testing.T, err error, expected string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Fatalf("error = %v, want substring %q", err, expected)
	}
}
