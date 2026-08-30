// aa-flatpak - Confine flatpak applications
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/flatpak"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

const usage = `aa-flatpak [-h] [--daemon] [-c | -e] [application-id...]

    Confine Flatpak applications with AppArmor profiles generated from their
    metadata and permissions.

    Tailor profiles for each application and integrate them with apparmor.d.

    It can be given optional application IDs to generate profiles for. If no ID
    is provided, it generates profiles for all installed applications.

Options:
    -h, --help     Show this help message and exit.
    -d, --daemon   Run as a daemon and watch for Flatpak changes.
    -c, --complain Generate profiles in complain mode (default).
    -e, --enforce  Generate profiles in enforce mode.

See https://apparmor.pujol.io/flatpak/ for more details.
`

var (
	// Command line options
	help     bool
	daemon   bool
	complain bool
	enforce  bool
)

var (
	profilesDir  = aa.MagicRoot.Join("flatpak-apps")
	appDir       = paths.New("/var/lib/flatpak/app")
	overridesDir string
)

// metadataPath returns the metadata file of an installed application.
func metadataPath(name string) *paths.Path {
	return appDir.Join(name, "current/active/metadata")
}

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.BoolVar(&daemon, "daemon", false, "Run as daemon and watch for profile changes.")
	flag.BoolVar(&daemon, "d", false, "Run as daemon and watch for profile changes.")
	flag.BoolVar(&complain, "complain", false, "Generate profiles in complain mode (default).")
	flag.BoolVar(&complain, "c", false, "Generate profiles in complain mode (default).")
	flag.BoolVar(&enforce, "enforce", false, "Generate profiles in enforce mode.")
	flag.BoolVar(&enforce, "e", false, "Generate profiles in enforce mode.")
}

func generateProfile(name string) (*paths.Path, error) {
	file := metadataPath(name)
	if !file.Exist() {
		return nil, fmt.Errorf("%s is not an installed flatpak application", name)
	}

	meta, err := flatpak.LoadFlatpakMetadata(file)
	if err != nil {
		return nil, err
	}
	if err := meta.LoadOverride(paths.New(overridesDir)); err != nil {
		return nil, err
	}
	meta.Requirements()

	mode := "complain"
	if enforce {
		mode = "enforce"
	}
	profile := flatpak.NewFlatpakAppArmorProfile(meta, mode)
	if err := profile.Generate(); err != nil {
		return nil, fmt.Errorf("failed to generate profile: %v", err)
	}

	profile.Format()
	if err := profilesDir.MkdirAll(); err != nil {
		return nil, err
	}
	out := profilesDir.Join(profile.FileName)
	if err := out.WriteFile([]byte(profile.String())); err != nil {
		return nil, err
	}
	return out, nil
}

func aaFlatpakGenerate() error {
	names := flag.Args()
	if len(names) == 0 {
		apps, err := appDir.ReadDir(paths.FilterDirectories())
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to list installed flatpak applications: %v", err)
		}
		for _, app := range apps {
			names = append(names, app.Base())
		}
	}

	profiles := paths.PathList{}
	for _, name := range names {
		log.Printf("Generating profile for %s", name)
		out, err := generateProfile(name)
		if err != nil {
			return err
		}
		profiles = append(profiles, out)
	}

	return util.ReloadProfiles(profiles)
}

func aaFlatpakDaemon() error {
	if err := aaFlatpakGenerate(); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fw, err := NewFileWatcher(
		func(name string) error {
			if !metadataPath(name).Exist() {
				return nil
			}
			log.Printf("Generating profile for %s", name)
			out, err := generateProfile(name)
			if err != nil {
				return err
			}
			log.Printf("Reloading profile for %s", name)
			return util.ReloadProfiles(paths.PathList{out})
		},
		func(name string) error {
			profile := profilesDir.Join(flatpak.ProfileName(name))
			if !profile.Exist() {
				return nil
			}
			log.Printf("Removing profile for %s", name)
			return profile.Remove()
		},
	)
	if err != nil {
		return err
	}
	defer fw.Close()

	return fw.Watch(ctx, appDir.String(), overridesDir)
}

func run() error {
	if complain && enforce {
		return fmt.Errorf("--complain and --enforce are mutually exclusive")
	}

	homes, _ := paths.New("/home").ReadDir(paths.FilterDirectories())
	if len(homes) > 0 {
		overridesDir = homes[0].Join(".local/share/flatpak/overrides").String()
	}

	if daemon {
		return aaFlatpakDaemon()
	}
	return aaFlatpakGenerate()
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if err := run(); err != nil {
		log.Fatalf("failed: %v", err)
	}
}
