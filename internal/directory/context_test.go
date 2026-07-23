package directory

import (
	"context"
	"path/filepath"
	"testing"

	gitinspect "github.com/setup-env/app/internal/git"
)

type fakeInspector struct {
	repository gitinspect.Repository
	ok         bool
}

func (f fakeInspector) Repository(context.Context, string) (gitinspect.Repository, bool) {
	return f.repository, f.ok
}

func TestDetectOrganizationDirectory(t *testing.T) {
	root := filepath.Join(t.TempDir(), "dev")
	current := filepath.Join(root, "setup-env")
	contextInfo, err := Detect(context.Background(), root, current, fakeInspector{})
	if err != nil {
		t.Fatal(err)
	}
	if contextInfo.Type != TypeOrganization || contextInfo.Organization != "setup-env" || contextInfo.Repository != "" {
		t.Fatalf("Detect() = %#v", contextInfo)
	}
}

func TestDetectRepositoryDirectory(t *testing.T) {
	root := filepath.Join(t.TempDir(), "dev")
	repositoryRoot := filepath.Join(root, "setup-env", "app")
	contextInfo, err := Detect(context.Background(), root, repositoryRoot, fakeInspector{
		ok: true,
		repository: gitinspect.Repository{
			Root:  repositoryRoot,
			Owner: "setup-env",
			Name:  "app",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if contextInfo.Type != TypeRepository || contextInfo.Organization != "setup-env" || contextInfo.Repository != "app" {
		t.Fatalf("Detect() = %#v", contextInfo)
	}
	if contextInfo.Git == nil {
		t.Fatal("Detect() Git = nil, want repository metadata")
	}
	if contextInfo.RemoteOrganization != "setup-env" || contextInfo.RemoteRepository != "app" {
		t.Fatalf("Detect() remote = %#v", contextInfo)
	}
}

func TestDetectArbitraryNestedDirectory(t *testing.T) {
	root := filepath.Join(t.TempDir(), "dev")
	current := filepath.Join(root, "setup-env", "notes")
	contextInfo, err := Detect(context.Background(), root, current, fakeInspector{})
	if err != nil {
		t.Fatal(err)
	}
	if contextInfo.Type != TypeOther || contextInfo.Repository != "" {
		t.Fatalf("Detect() = %#v", contextInfo)
	}
}

func TestClassifyPathShapesAcrossPlatforms(t *testing.T) {
	tests := []struct {
		name     string
		relative string
		git      bool
		wantType Type
		wantOrg  string
		wantRepo string
	}{
		{name: "Windows organization", relative: `setup-env`, wantType: TypeOrganization, wantOrg: "setup-env"},
		{name: "Windows repository", relative: `setup-env\app`, git: true, wantType: TypeRepository, wantOrg: "setup-env", wantRepo: "app"},
		{name: "macOS repository", relative: "setup-env/app", git: true, wantType: TypeRepository, wantOrg: "setup-env", wantRepo: "app"},
		{name: "Linux arbitrary", relative: "setup-env/scratch", wantType: TypeOther, wantOrg: "setup-env"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotType, gotOrg, gotRepo := ClassifyShape(test.relative, test.git)
			if gotType != test.wantType || gotOrg != test.wantOrg || gotRepo != test.wantRepo {
				t.Fatalf("ClassifyShape() = %q, %q, %q", gotType, gotOrg, gotRepo)
			}
		})
	}
}
