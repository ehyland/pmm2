package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

func newUpdateSelfCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "update-self",
		Short: "Update pmm2 itself",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateSelf(version)
		},
	}
}

func runUpdateSelf(version string) error {
	updater, err := selfupdate.NewUpdater(selfupdate.Config{})
	if err != nil {
		return err
	}
	latest, found, err := updater.DetectLatest(context.Background(), selfupdate.ParseSlug("ehyland/pmm2"))
	if err != nil {
		return err
	}
	if !found {
		fmt.Println("No updates found.")
		return nil
	}

	if latest.LessOrEqual(version) {
		fmt.Printf("pmm2 %s is already up to date\n", version)
		return nil
	}

	fmt.Printf("Updating to %s...\n", latest.Version())
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	if absPath, err := filepath.Abs(exePath); err == nil {
		exePath = absPath
	}
	if err := updater.UpdateTo(context.Background(), latest, exePath); err != nil {
		return err
	}

	fmt.Println("Self-update successful. Synchronizing shims...")

	// Use syscall.Exec to run the NEW binary with the 'setup' command
	// This ensures the new config.Shims list from the updated binary is used.
	return syscall.Exec(exePath, []string{filepath.Base(exePath), "setup"}, os.Environ())
}
