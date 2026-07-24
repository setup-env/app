# Release operations

## Version policy

Setup Env uses semantic tags in the form `vMAJOR.MINOR.PATCH`. The prepared first
public version is `v0.1.0`. Never reuse a tag or recreate a published version
with different contents; publish a patch version instead.

Release commits must be merged into clean, synchronized `main`. Releases are
never created from feature branches.

## Prepare

1. Choose an unused version and confirm it has no tag or GitHub Release.
2. Update `CHANGELOG.md` and `.github/RELEASE_TEMPLATE.md`.
3. Confirm README installer commands and supported targets.
4. Run:

   ```sh
   go fmt ./...
   git diff --check
   go vet ./...
   go test ./...
   go run ./cmd/setup-env module validate-catalog
   go run ./cmd/setup-env module validate examples/setup-env.yaml
   make RELEASE_VERSION=v0.1.0 release-snapshot
   make RELEASE_VERSION=v0.1.0 release-verify
   ```

5. Inspect `dist/`: six archives plus `checksums.txt`, with no source tree,
   absolute paths, credentials, caches, or debug files.
6. Test a native installer against the local `dist/` fixture.
7. Merge through a pull request and wait for `main` CI.

## Publish

Public tagging is an explicit release decision. From synchronized `main`:

```sh
git tag -a v0.1.0 -m "Setup Env v0.1.0"
git push origin v0.1.0
```

The tag-triggered workflow validates the exact tagged commit, full test/race
suite, catalog, installer functions, release metadata, archive contents, and
checksums before creating the GitHub Release. Monitor the workflow until every
job and the release publication step succeeds.

## Verify after publication

1. Confirm the release is not a draft and contains exactly six archives plus
   `checksums.txt`.
2. Download all assets and verify every checksum.
3. Extract a native archive and run `setup-env version`.
4. Test a Windows installation on Windows and at least one macOS/Linux
   installation on its native platform.
5. Test upgrade and uninstall without removing unrelated files.
6. Record the workflow run and actual platform coverage in milestone notes.

Do not claim an untested platform installation.

## Failed releases and rollback

The workflow removes a partial GitHub Release if upload fails but leaves the tag
for investigation. Do not rerun publication blindly.

- If no public release was created, diagnose the workflow, delete the remote tag
  only after confirming no user could have consumed it, fix on `main`, and use a
  fresh tag according to repository policy.
- If a release became public, do not replace its assets or recreate its tag.
  Mark it clearly as affected, publish remediation guidance, and issue the next
  patch version.
- Bootstrap upgrades retain the previous binary until the verified replacement
  executes successfully and restore it when replacement validation fails.

GitHub Releases are the rollback record; installers do not implement background
or automatic updates.

## Protection recommendations

The repository currently has no branch protection or ruleset. Recommended
hardening is to require pull requests and successful CI on `main`, restrict tag
creation for `v*`, prevent tag deletion/rewrites, and limit release workflow
modification to reviewed changes. Do not bypass configured protection for
routine releases.
