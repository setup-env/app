package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	releaser "github.com/setup-env/app/internal/release"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "release:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, arguments []string) error {
	flags := flag.NewFlagSet("release", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	version := flags.String("version", "v0.1.0-snapshot", "semantic release version")
	commit := flags.String("commit", "", "commit SHA (defaults to HEAD)")
	dateText := flags.String("date", "", "RFC3339 build date (defaults to commit date)")
	dirtyText := flags.String("dirty", "auto", "true, false, or auto")
	output := flags.String("output", "dist", "release output directory")
	targetText := flags.String("targets", "all", "comma-separated os/arch targets")
	cleanOwned := flags.Bool("clean-owned", false, "remove only prior setup-env release artifacts")
	verifyOnly := flags.Bool("verify-only", false, "verify existing output without building")
	if err := flags.Parse(arguments); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}
	targets, err := releaser.ParseTargets(*targetText)
	if err != nil {
		return err
	}
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	outputPath := *output
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(root, outputPath)
	}
	if *verifyOnly {
		if err := releaser.Verify(outputPath, *version, targets); err != nil {
			return err
		}
		fmt.Printf("verified %d release archives in %s\n", len(targets), outputPath)
		return nil
	}

	gitCommit, gitDate, gitDirty, err := releaser.ResolveGitMetadata(ctx, root)
	if err != nil {
		return err
	}
	if *commit == "" {
		*commit = gitCommit
	}
	buildDate := gitDate
	if *dateText != "" {
		buildDate, err = time.Parse(time.RFC3339, *dateText)
		if err != nil {
			return fmt.Errorf("parse -date: %w", err)
		}
	}
	dirty := gitDirty
	switch *dirtyText {
	case "true":
		dirty = true
	case "false":
		if gitDirty {
			return fmt.Errorf("refusing -dirty=false because the repository has tracked or untracked changes")
		}
		dirty = false
	case "auto":
	default:
		return fmt.Errorf("-dirty must be true, false, or auto")
	}
	artifacts, err := releaser.Build(ctx, releaser.Config{
		Root:       root,
		Output:     outputPath,
		Version:    *version,
		Commit:     *commit,
		BuildDate:  buildDate,
		Dirty:      dirty,
		Targets:    targets,
		CleanOwned: *cleanOwned,
	})
	if err != nil {
		return err
	}
	for _, artifact := range artifacts {
		fmt.Printf("%s  %s\n", artifact.Checksum, artifact.Name)
	}
	fmt.Printf("built and verified %d release archives in %s\n", len(artifacts), outputPath)
	return nil
}
