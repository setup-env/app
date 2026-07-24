# Install Setup Env on Windows

## Build from source

Requirements:

- Git;
- Go 1.26 (the repository pins toolchain 1.26.5);
- a user-writable directory on `PATH`.

PowerShell:

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\dev\setup-env"
Set-Location "$env:USERPROFILE\dev\setup-env"
git clone https://github.com/setup-env/app.git
Set-Location app
go build -o bin\setup-env.exe .\cmd\setup-env
.\bin\setup-env.exe
```

To install for the current user, copy `setup-env.exe` to a user-owned bin
directory, such as `%USERPROFILE%\bin`, and add that directory to the user
`PATH`. Avoid copying into protected system directories or using an elevated
shell solely for Setup Env.

## Install a future release

Signed GitHub Release binaries and a Windows bootstrap path are planned for
Milestone 05. They do not exist yet. When available, this guide will describe
artifact verification, installation, upgrades, and rollback; do not treat any
current third-party package as an official Setup Env release.
