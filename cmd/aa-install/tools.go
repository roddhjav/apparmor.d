// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
)

const (
	// manifestFile is the name of the manifest file that tracks installed profiles.
	manifestFile = "install.db"

	// maxDetailedFiles is the number of changed files below which they are
	// listed one by one instead of only counted.
	maxDetailedFiles = 20

	// symlinkPrefix marks a manifest entry as a symlink; the remainder is the
	// link target. Disable links (disable/<name> -> ../<name>) must be installed
	// as symlinks, not dereferenced. The upstream profile they disable only
	// exists on the target, so the link is dangling in the build directory.
	symlinkPrefix = "symlink:"
)

// hashFile returns the hex-encoded SHA-256 hash of a file's contents.
func hashFile(path *paths.Path) (string, error) {
	data, err := path.ReadFile()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// fileIdentity returns a change-detection identity for a build/target entry:
// the (prefixed) link target for a symlink, otherwise the content hash.
func fileIdentity(path *paths.Path) (string, error) {
	isLink, err := path.IsSymlink()
	if err != nil {
		return "", err
	}
	if isLink {
		target, err := os.Readlink(path.String())
		if err != nil {
			return "", err
		}
		return symlinkPrefix + target, nil
	}
	return hashFile(path)
}

// copyEntry writes the build entry to dst, recreating symlinks as symlinks.
func copyEntry(file, dst *paths.Path, ident string) error {
	if target, ok := strings.CutPrefix(ident, symlinkPrefix); ok {
		return os.Symlink(target, dst.String())
	}
	return file.CopyTo(dst)
}

// readManifest reads the manifest as a map of relative path → hash.
func readManifest(stateDir *paths.Path) map[string]string {
	res := map[string]string{}
	path := stateDir.Join(manifestFile)
	if !path.Exist() {
		return res
	}
	data, err := path.ReadFile()
	if err != nil {
		logging.Warning("Cannot read manifest %s: %s", path, err)
		return res
	}
	for line := range strings.SplitSeq(strings.TrimSpace(string(data)), "\n") {
		if fields := strings.SplitN(line, " ", 2); len(fields) == 2 {
			res[fields[1]] = fields[0]
		}
	}
	return res
}

// writeManifest writes the manifest as "hash path" lines, sorted by path.
func writeManifest(stateDir *paths.Path, entries map[string]string) error {
	if err := stateDir.MkdirAll(); err != nil {
		return err
	}
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	var buf strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&buf, "%s %s\n", entries[k], k)
	}
	return stateDir.Join(manifestFile).WriteFile([]byte(buf.String()))
}

// install copies files from the build directory to the install target
// and removes files that were previously installed but are no longer present.
// It uses a manifest file in stateDir to track installed files.
// Returns whether any file was added, updated, or removed.
func installProfiles(buildDir *paths.Path, targetDir *paths.Path, stateDir *paths.Path) (bool, error) {
	previous := readManifest(stateDir)

	files, err := buildDir.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	if err != nil {
		return false, err
	}

	newManifest := make(map[string]string, len(files))
	var added, updated, skipped []string
	var unchanged int
	for _, file := range files {
		rel, err := file.RelFrom(buildDir)
		if err != nil {
			return false, err
		}
		relStr := rel.String()

		ident, err := fileIdentity(file)
		if err != nil {
			return false, err
		}

		_, tracked := previous[relStr]
		delete(previous, relStr)

		// Compare against the target itself, not the manifest, so that
		// drifted or missing files are repaired on reinstall. Lstat so a
		// symlink (possibly dangling) counts as present.
		dst := targetDir.JoinPath(rel)
		_, lerr := dst.Lstat()
		existed := lerr == nil

		// Never take over a file we do not own: an untracked file on the
		// target belongs to the admin or another package. Leave it alone
		// and keep it out of the manifest so uninstall does not remove it.
		if existed && !tracked {
			// An existing disable/<name> link, whatever it points to, means
			// the profile is already disabled: our link is redundant, not in
			// conflict. Keep theirs and stay quiet about it.
			if !strings.HasPrefix(relStr, "disable/") {
				skipped = append(skipped, relStr)
			}
			continue
		}
		newManifest[relStr] = ident

		if existed {
			current, err := fileIdentity(dst)
			if err != nil {
				return false, err
			}
			if current == ident {
				unchanged++
				continue
			}
			// Remove first: the entry type (file vs symlink) may change.
			if err := dst.RemoveAll(); err != nil {
				return false, err
			}
		}

		if err := dst.Parent().MkdirAll(); err != nil {
			return false, err
		}
		if err := copyEntry(file, dst, ident); err != nil {
			return false, err
		}

		if existed {
			updated = append(updated, relStr)
		} else {
			added = append(added, relStr)
		}
	}

	var removed []string
	for old := range previous {
		target := targetDir.Join(old)
		if _, err := target.Lstat(); err == nil {
			if err := target.RemoveAll(); err != nil {
				return false, err
			}
			removed = append(removed, old)
		}
	}
	slices.Sort(removed)

	if err := writeManifest(stateDir, newManifest); err != nil {
		return false, err
	}

	logging.Indent = ""
	logging.Success("Installed %d profiles to %s", len(newManifest), targetDir)
	logging.Indent = "   "
	logging.Bullet("%d added, %d updated, %d unchanged", len(added), len(updated), unchanged)
	if len(added)+len(updated) < maxDetailedFiles {
		bulletList("Added", added)
		bulletList("Updated", updated)
	}
	bulletList(fmt.Sprintf("Removed %d stale files", len(removed)), removed)
	logging.Indent = ""
	if len(skipped) > 0 {
		slices.Sort(skipped)
		logging.Warning("Skipped %d untracked files already present in %s:", len(skipped), targetDir)
		logging.Indent = "   "
		for _, rel := range skipped {
			logging.Bullet("%s", rel)
		}
		logging.Indent = ""
	}
	return len(added)+len(updated)+len(removed) > 0, nil
}

// bulletList prints the sorted files under title, indented one level deeper.
func bulletList(title string, files []string) {
	if len(files) == 0 {
		return
	}
	slices.Sort(files)
	logging.Bullet("%s:", title)
	logging.Indent = "      "
	for _, rel := range files {
		logging.Bullet("%s", rel)
	}
	logging.Indent = "   "
}
