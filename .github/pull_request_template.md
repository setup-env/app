## Summary

<!-- What changes, and why? -->

## Scope

- [ ] Application/runtime change
- [ ] Module manifest or catalog submission
- [ ] Documentation only

## Module submission

<!-- Complete for catalog changes; otherwise write "Not applicable". -->

- Module ID:
- Repository:
- Requested trust: community / verified / official
- Requested status:
- Manifest path:
- Security, elevation, network, and secret-input considerations:
- Maintenance and release expectations:

## Validation

- [ ] `go fmt ./...`
- [ ] `git diff --check`
- [ ] `go vet ./...`
- [ ] `go test ./...`
- [ ] `setup-env module validate-catalog`
- [ ] Relevant `setup-env module validate <path>`

## Boundaries

- [ ] No secrets or credentials are included.
- [ ] Trust is assigned by catalog review, not the module manifest.
- [ ] This change does not add unreviewed download or execution behavior.
