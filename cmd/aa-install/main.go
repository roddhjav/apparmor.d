// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/builder"
	"github.com/roddhjav/apparmor.d/pkg/configure"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/runtime"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

const (
	nilConfig = ""
	nilMagic  = ""
	nilSrc    = ""
	usage     = `aa-install [-h] [--config DIR] [--src DIR] [--magic DIR] [ -i | -u | -s | -l ] [-e|-c]

    Install and manage apparmor profiles from apparmor.d.

    With no action flag, print the installation status summary.

Options:
    -h, --help         Show this help message and exit.
    -s, --status       Show the installation status summary (default).
    -l, --list         List installed profile paths from the manifest.
    -i, --install      Install the profiles.
    -c, --complain     Set complain flag on all the profiles.
    -e, --enforce      Set enforce flag on all the profiles.
    -u, --uninstall    Remove all profiles installed.
        --no-reload    Do not reload the profiles after modifying them.
        --config DIR   Select an alternate configuration directory (default: /etc/apparmor/).
        --magic DIR    Select an alternate apparmor.d directory (default: /etc/apparmor.d).
        --src DIR      Select an alternate source directory (default: /usr/share/apparmor.d).

Configuration files:
    modes              Modes of enforcement for all profiles.
    flags.d/*.conf     Set per-profile flags.
    ignore.d/*.conf    Set (group of) profiles to ignore.
    include.d/*.conf   Only install the given (group of) profiles.

Configuration directories:

    /usr/share/apparmor/ - Vendor defaults are read first
    /etc/apparmor/       - Local admin configuration, can be overwritten by --config DIR.

See man aa-install(1) for more information.
`
)

// Command line options
var (
	help      bool
	install   bool
	complain  bool
	enforce   bool
	uninstall bool
	status    bool
	list      bool
	noReload  bool
	config    string
	magic     string
	src       string
)

// reloadAppArmor is swappable so tests can skip the system reload.
var reloadAppArmor = util.ReloadAppArmor

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.BoolVar(&install, "i", false, "Install the profiles.")
	flag.BoolVar(&install, "install", false, "Install the profiles.")
	flag.BoolVar(&complain, "c", false, "Set complain flag on all the profiles.")
	flag.BoolVar(&complain, "complain", false, "Set complain flag on all the profiles.")
	flag.BoolVar(&enforce, "e", false, "Set enforce flag on all the profiles.")
	flag.BoolVar(&enforce, "enforce", false, "Set enforce flag on all the profiles.")
	flag.BoolVar(&uninstall, "u", false, "Remove all profiles installed.")
	flag.BoolVar(&uninstall, "uninstall", false, "Remove all profiles installed.")
	flag.BoolVar(&status, "s", false, "Show the installation status summary.")
	flag.BoolVar(&status, "status", false, "Show the installation status summary.")
	flag.BoolVar(&list, "l", false, "List installed profile paths from the manifest.")
	flag.BoolVar(&list, "list", false, "List installed profile paths from the manifest.")
	flag.BoolVar(&noReload, "no-reload", false, "Do not reload the profiles after modifying them.")
	flag.StringVar(&config, "config", nilConfig, "Select an alternate configuration directory (default: /etc/apparmor/).")
	flag.StringVar(&magic, "magic", nilMagic, "Select an alternate apparmor.d directory (default: /etc/apparmor.d).")
	flag.StringVar(&src, "src", nilSrc, "Select an alternate source directory (default: /usr/share/apparmor.d).")
}

// aaStatus prints a summary of the installation state: the number of
// profiles recorded in the manifest, and how many are missing or drifted
// (hash mismatch) from the install target.
func aaStatus(stateDir *paths.Path, targetDir *paths.Path) error {
	manifest := readManifest(stateDir)
	if len(manifest) == 0 {
		logging.Warning("No profiles installed")
		return nil
	}

	var missing, drifted, ok int
	for rel, ident := range manifest {
		target := targetDir.Join(rel)
		// Lstat so a symlink entry (possibly dangling) counts as present.
		if _, err := target.Lstat(); err != nil {
			missing++
			continue
		}
		current, err := fileIdentity(target)
		if err != nil {
			return err
		}
		if current != ident {
			drifted++
			continue
		}
		ok++
	}

	logging.Indent = ""
	logging.Success("%d profiles installed in %s", len(manifest), targetDir)
	logging.Indent = "   "
	logging.Bullet("%d up-to-date, %d drifted, %d missing", ok, drifted, missing)
	logging.Indent = ""
	return nil
}

// aaList prints one absolute installed-profile path per line, sorted.
// The paths are resolved against targetDir.
func aaList(stateDir *paths.Path, targetDir *paths.Path) error {
	manifest := readManifest(stateDir)
	if len(manifest) == 0 {
		return nil
	}
	keys := make([]string, 0, len(manifest))
	for k := range manifest {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, rel := range keys {
		fmt.Println(targetDir.Join(rel).String())
	}
	return nil
}

func aaInstall(configDir *paths.Path, srcDir *paths.Path, cfg *conf) (bool, error) {
	tmp, err := os.MkdirTemp("", "aa-install")
	if err != nil {
		return false, err
	}
	buildDir := paths.New(tmp)
	defer buildDir.RemoveAll()

	c := tasks.NewTaskConfig(buildDir)
	r := runtime.NewRunners(c)

	// Add default configure tasks
	r.Configures.

		// Initialize a new clean apparmor.d build directory
		Add(configure.NewSynchronise([]*paths.Path{srcDir}))

	// If not empty, only include profiles defined from the include.d dirs
	include := configure.NewInclude(cfg.includeDirs)
	if include.Active() {
		r.Configures.Add(include)
	}

	// Continue adding default configure tasks
	r.Configures.

		// Ignore profiles and files from the ignore.d dirs
		Add(configure.NewUserIgnore(cfg.ignoreDirs)).

		// Merge profiles (from group/, profiles-*-*/) to a unified apparmor.d directory
		Add(configure.NewMerge()).

		// Rename profiles that replace an upstream one (shipped disable/ links)
		Add(configure.NewOverwriteFromLinks()).

		// Set user-defined flags from the flags.d dirs
		Add(configure.NewSetFlags(cfg.flagDirs)).

		// Set detected system state in tunables/multiarch.d/state
		// Add(configure.NewSetState()).

		// Only keep profiles for installed programs
		Add(configure.NewSelectInstalled())

	// Apply the default deploy mode to every profile, except those a user
	// flags.d drop-in already assigns a mode to (SetFlags handled those).
	r.Builders.Add(builder.NewDeployMode(cfg.mode, userModeOverrides(cfg.flagDirs)))

	if err := r.Configure(); err != nil {
		return false, err
	}

	if err := r.Build(); err != nil {
		return false, err
	}
	return installProfiles(c.RootApparmor, aa.MagicRoot, configDir)
}

// aaUninstall removes all files recorded in the manifest, then the manifest
// itself. Returns whether any file was removed.
func aaUninstall(stateDir *paths.Path, targetDir *paths.Path) (bool, error) {
	manifest := readManifest(stateDir)
	if len(manifest) == 0 {
		logging.Warning("No manifest found, nothing to uninstall")
		return false, nil
	}

	removed := 0
	for rel := range manifest {
		target := targetDir.Join(rel)
		// Lstat so a dangling disable symlink still counts as present.
		if _, err := target.Lstat(); err == nil {
			if err := target.RemoveAll(); err != nil {
				return false, err
			}
			removed++
		}
	}

	manifestPath := stateDir.Join(manifestFile)
	if manifestPath.Exist() {
		if err := manifestPath.Remove(); err != nil {
			return false, err
		}
	}

	logging.Success("Uninstalled %d profiles from %s", removed, targetDir)
	return removed > 0, nil
}

func run() error {
	if complain && enforce {
		return fmt.Errorf("--complain and --enforce are mutually exclusive")
	}

	configDir := paths.New("/etc/apparmor/")
	if config != nilConfig {
		configDir = paths.New(config)
	}
	if magic != nilMagic {
		aa.MagicRoot = paths.New(magic)
	}
	srcRoot := paths.New("/usr/share/apparmor.d")
	if src != nilSrc {
		srcRoot = paths.New(src)
	}
	cfg, err := loadConfig(configDir)
	if err != nil {
		return err
	}

	var changed bool
	switch {
	case list:
		return aaList(configDir, aa.MagicRoot)

	case install || complain || enforce:
		changed, err = aaInstall(configDir, srcRoot, cfg)

	case uninstall:
		changed, err = aaUninstall(configDir, aa.MagicRoot)

	default:
		return aaStatus(configDir, aa.MagicRoot)

	}
	if err != nil {
		return err
	}
	if changed && !noReload {
		return reloadAppArmor()
	}
	return nil
}

func main() {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if err := run(); err != nil {
		logging.Fatal("%s", err.Error())
	}
}
