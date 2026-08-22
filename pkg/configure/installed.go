// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	// skipSystemDirs filters out the top-level directories that are not
	// per-program profiles (shared building blocks, compiled out, mapping data).
	skipSystemDirs = paths.FilterOutNames("tunables", "abstractions", "disable", "mappings", "namespaces")

	// alwaysKeep lists profile-name globs installed even when they have no
	// attachment and their group is not otherwise installed.
	alwaysKeep = []string{
		"child-*", // groups/children
		"glycin",  // groups/children
		"dbus-*",  // groups/bus core buses
		":*",      // namespaces/* (merged as :<ns>:<name>)
	}

	// ignoredSuffixes lists file extensions that are not executable programs
	// to confine.
	ignoredSuffixes = []string{
		".so",
		".zst",
		".py",
		".pc",
		".service",
		".target",
		".socket",
		".mount",
		".timer",
	}
)

// isAlwaysKept reports whether a profile basename matches alwaysKeep.
func isAlwaysKept(base string) bool {
	return slices.ContainsFunc(alwaysKeep, func(pattern string) bool {
		ok, _ := filepath.Match(pattern, base)
		return ok
	})
}

func hasIgnoredSuffix(name string) bool {
	return slices.ContainsFunc(ignoredSuffixes, func(s string) bool {
		return strings.HasSuffix(name, s)
	})
}

// sbinGlob returns the filesystem-glob form of `@{sbin}` for the current
// distribution. On Arch, /usr/sbin is a symlink to /usr/bin, so both
// /sbin and /bin must match.
func sbinGlob() string {
	if tasks.Family == "pacman" {
		return "/{,usr/}{,s}bin"
	}
	return "/{,usr/}sbin"
}

type SelectInstalled struct {
	tasks.BaseTask
	include map[string]bool // profiles or groups kept even when not installed
}

// NewSelectInstalled creates a SelectInstalled task. The include entries
// (profile or group names) are kept even when their program is not installed.
func NewSelectInstalled(include ...string) *SelectInstalled {
	set := make(map[string]bool, len(include))
	for _, entry := range include {
		set[filepath.Base(entry)] = true
	}
	return &SelectInstalled{
		BaseTask: tasks.BaseTask{
			Keyword: "installed",
			Msg:     "Keep profiles for installed programs",
		},
		include: set,
	}
}

// profileState holds the install decision inputs gathered in the first pass.
type profileState struct {
	file      *paths.Path
	hasAtt    bool // profile has an attachment (i.e. is an attachable program)
	installed bool // an attachment resolves to an installed program
	force     bool // keep unconditionally (read/parse error: fail-safe)
}

func (p SelectInstalled) Apply() ([]string, error) {
	var res []string

	files, err := p.RootApparmor.ReadDirRecursiveFiltered(
		skipSystemDirs, paths.FilterOutDirectories(),
	)
	if err != nil {
		return res, err
	}

	// First pass: classify every profile and record which groups have at
	// least one installed attachable profile.
	states := make([]profileState, 0, len(files))
	installedGroups := map[string]bool{}
	for _, file := range files {
		st := profileState{file: file}
		// On parse/read errors, keep the profile (fail-safe: don't silently drop).
		profile, err := file.ReadFileAsString()
		if err != nil {
			st.force = true
			log.Printf("Warning: failed to read profile %s: %v", file, err)
			states = append(states, st)
			continue
		}
		att, err := getAttachments(profile)
		if err != nil {
			st.force = true
			states = append(states, st)
			continue
		}
		st.hasAtt = len(att) > 0
		st.installed = st.hasAtt && isInstalled(att)
		if st.installed {
			if group := p.Groups[file.Base()]; group != "" {
				installedGroups[group] = true
			}
		}
		states = append(states, st)
	}

	// Second pass: keep installed programs, keep attachment-less profiles that
	// are standalone or whose group is installed, drop the rest.
	var keptNames, removedNames []string
	for _, st := range states {
		base := st.file.Base()
		keep := st.force || st.installed || isAlwaysKept(base) ||
			p.include[base] || p.include[p.Groups[base]]
		if !keep && !st.hasAtt {
			group := p.Groups[st.file.Base()]
			keep = group == "" || installedGroups[group]
		}
		if !keep {
			if err := st.file.Remove(); err != nil {
				return res, err
			}
			removedNames = append(removedNames, st.file.Base())
			continue
		}
		keptNames = append(keptNames, st.file.Base())
	}

	res = append(res,
		fmt.Sprintf("Kept %d profiles", len(keptNames)),
		fmt.Sprintf("Ignored %d", len(removedNames)),
	)
	return res, nil
}

// getAttachments returns the resolved attachment paths from a
// post-prebuild profile header (e.g., "profile dolphin /{,usr/}bin/dolphin").
// Returns nil for child profiles and abstractions with no attachment.
func getAttachments(profile string) (att []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			att = nil
			err = fmt.Errorf("parse panic: %v", r)
		}
	}()

	f := aa.DefaultTunables()
	if _, err := f.Parse(profile); err != nil {
		return nil, err
	}
	return f.GetDefaultProfile().Attachments, nil
}

var regVariableRef = regexp.MustCompile(`@\{[^{}]+\}`)

// pathReplacer is the distro-aware AppArmor variable → filesystem glob
// substitution table. Order in strings.NewReplacer matters: more
// specific keys (e.g. @{user_share_dirs}) must precede generic ones
// (@{HOME}) so the specific match wins.
var pathReplacer = sync.OnceValue(func() *strings.Replacer {
	return strings.NewReplacer(
		"@{user_share_dirs}", "/home/*/.local/share",
		"@{user_config_dirs}", "/home/*/.config",
		"@{user_cache_dirs}", "/home/*/.cache",
		"@{user_state_dirs}", "/home/*/.local/state",
		"@{user_bin_dirs}", "/home/*/.local/bin",
		"@{user_lib_dirs}", "/home/*/.local/lib",
		"@{HOME}", "/home/*",
		"@{lib}", "/{,usr/}lib{,exec,32,64}",
		"@{sbin}", sbinGlob(),
		"@{bin}", "/{,usr/}bin",
		"@{etc_ro}", "/{,usr/}etc",
	)
})

// symlinkCache caches os.Lstat results for directory components to
// avoid repeated syscalls. A path that traverses any symlink directory
// is redundant (e.g., /bin/foo when /bin → /usr/bin) and can be skipped.
//
// Not safe for concurrent use: intentional for the single-pass
// prebuild CLI. If Apply ever parallelises profiles, guard this.
var symlinkCache = map[string]bool{}

// readDirCache memoises os.ReadDir results per directory. The host
// filesystem is static during install-detection, so a directory read once
// need never be read again across the ~1500 profiles. A nil slice caches a
// read error (missing/inaccessible dir).
//
// Not safe for concurrent use: same single-pass assumption as symlinkCache.
var readDirCache = map[string][]os.DirEntry{}

func cachedReadDir(dir string) []os.DirEntry {
	if entries, ok := readDirCache[dir]; ok {
		return entries
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		entries = nil
	}
	readDirCache[dir] = entries
	return entries
}

func hasSymlinkComponent(p string) bool {
	for i := 1; i < len(p); i++ {
		if p[i] != '/' {
			continue
		}
		dir := p[:i]
		isLink, cached := symlinkCache[dir]
		if !cached {
			info, err := os.Lstat(dir)
			isLink = err == nil && info.Mode()&os.ModeSymlink != 0
			symlinkCache[dir] = isLink
		}
		if isLink {
			return true
		}
	}
	return false
}

// pathIsInstalled returns whether a single resolved path indicates an
// installed program. Unresolvable variables are conservatively treated
// as installed; wildcard remnants and symlink-traversing paths are
// skipped (the non-symlinked duplicate will be stat'd instead).
func pathIsInstalled(p string) bool {
	if !strings.HasPrefix(p, "/") {
		return true
	}
	if strings.ContainsAny(p, "*?[{") {
		return false
	}
	if hasSymlinkComponent(p) {
		return false
	}
	_, err := os.Stat(p)
	return err == nil
}

// isInstalled reports whether any of a profile's attachments resolves to
// an existing filesystem path.
func isInstalled(attachments []string) bool {
	return slices.ContainsFunc(attachments, attachmentInstalled)
}

// attachmentInstalled reports whether a single AppArmor attachment matches
// an existing file. Known variables are substituted and unknown ones
// become wildcards; wildcard patterns (*, ?, [class], {alt}, **) are
// matched against the filesystem, plain paths are stat'd.
func attachmentInstalled(att string) bool {
	att = pathReplacer().Replace(att)
	att = regVariableRef.ReplaceAllString(att, "*")
	if strings.ContainsAny(att, "*?[{") {
		return globExists(att)
	}
	return pathIsInstalled(att)
}

// globExists reports whether any real file matches the AppArmor glob
// pattern. Structural brace alternations that span a path separator (e.g.
// {,usr/}) are expanded into concrete branches; intra-segment braces and
// character classes (e.g. the @{version} pattern [0-9]{[0-9],}...) are kept
// and matched as a regexp per path segment, so patterns like @{user} or a
// version glob never blow up combinatorially.
func globExists(pattern string) bool {
	return slices.ContainsFunc(expandPathBraces(pattern), func(p string) bool {
		return matchDescend("/", strings.Split(strings.TrimPrefix(p, "/"), "/"))
	})
}

// matchDescend walks a pattern segment by segment from dir, only reading
// the directories the pattern actually names. "**" matches any depth.
//
// ponytail: a non-terminal "**" (e.g. /a/**/b) recurses every subdir of a;
// in practice "**" is always terminal here, so no depth cap. Add one if a
// mid-pattern "**" over a huge tree ever shows up.
func matchDescend(dir string, segs []string) bool {
	for len(segs) > 0 && segs[0] == "" {
		segs = segs[1:]
	}
	if len(segs) == 0 {
		_, err := os.Lstat(dir)
		return err == nil
	}
	seg, rest := segs[0], segs[1:]

	if seg == "**" {
		if matchDescend(dir, rest) {
			return true
		}
		for _, e := range cachedReadDir(dir) {
			if e.IsDir() && matchDescend(filepath.Join(dir, e.Name()), segs) {
				return true
			}
		}
		return false
	}

	if !strings.ContainsAny(seg, "*?[{") {
		return matchDescend(filepath.Join(dir, seg), rest)
	}

	re, err := segRegexp(seg)
	if err != nil {
		return true // unparseable glob: keep the profile (fail-safe)
	}
	for _, e := range cachedReadDir(dir) {
		if re.MatchString(e.Name()) && matchDescend(filepath.Join(dir, e.Name()), rest) {
			return true
		}
	}
	return false
}

// segRegexp compiles a single AppArmor path segment (no '/') into an
// anchored regexp: * → any run of non-slash, ? → one char, {a,b} → (a|b),
// and [class] kept verbatim.
func segRegexp(seg string) (*regexp.Regexp, error) {
	var b strings.Builder
	b.WriteByte('^')
	depth := 0
	for i := 0; i < len(seg); i++ {
		c := seg[i]
		switch c {
		case '[':
			b.WriteByte('[')
			for i++; i < len(seg) && seg[i] != ']'; i++ {
				b.WriteByte(seg[i])
			}
			b.WriteByte(']')
		case '*':
			b.WriteString("[^/]*")
		case '?':
			b.WriteString("[^/]")
		case '{':
			depth++
			b.WriteByte('(')
		case '}':
			if depth > 0 {
				depth--
				b.WriteByte(')')
			} else {
				b.WriteString(`\}`)
			}
		case ',':
			if depth > 0 {
				b.WriteByte('|')
			} else {
				b.WriteByte(',')
			}
		case '.', '+', '(', ')', '|', '^', '$', '\\':
			b.WriteByte('\\')
			b.WriteByte(c)
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte('$')
	return regexp.Compile(b.String())
}

// expandPathBraces expands only brace groups that span a path separator
// (structural alternations like {,usr/} or {/bin/x,/opt/y}), leaving
// intra-segment braces verbatim for segRegexp to handle.
func expandPathBraces(pattern string) []string {
	searchFrom := 0
	for {
		rel := strings.IndexByte(pattern[searchFrom:], '{')
		if rel == -1 {
			return []string{pattern}
		}
		start := searchFrom + rel
		end := matchingBrace(pattern, start)
		if end == -1 {
			return []string{pattern}
		}
		content := pattern[start+1 : end]
		if !strings.Contains(content, "/") {
			searchFrom = end + 1 // intra-segment brace: leave for segRegexp
			continue
		}
		prefix, suffix := pattern[:start], pattern[end+1:]
		var out []string
		for _, alt := range splitAtCommas(content) {
			out = append(out, expandPathBraces(prefix+alt+suffix)...)
		}
		return out
	}
}

// matchingBrace returns the index of the '}' closing the '{' at start, or
// -1 if unbalanced.
func matchingBrace(s string, start int) int {
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// splitAtCommas splits a string at top-level commas, respecting nested braces.
func splitAtCommas(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
}
