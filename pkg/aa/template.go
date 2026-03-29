// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"embed"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	// Default indentation for apparmor profile (2 spaces)
	Indentation = "  "

	// The current indentation level
	IndentationLevel = 0

	//go:embed templates/*.gotmpl
	//go:embed templates/rule/*.gotmpl
	tmplFiles embed.FS

	// The functions available in the template
	tmplFunctionMap = template.FuncMap{
		"kindof":     kindOf,
		"join":       join,
		"cjoin":      cjoin,
		"indent":     indent,
		"overindent": indentDbus,
		"setindent":  setindent,
	}

	// The apparmor templates
	tmpl = generateTemplates([]Kind{
		// Global templates
		"apparmor", PROFILE, HAT, "rules",

		// Preamble templates
		ABI, ALIAS, INCLUDE, VARIABLE, COMMENT,

		// Rules templates
		ALL, RLIMIT, USERNS, CAPABILITY, NETWORK,
		MOUNT, REMOUNT, UMOUNT, PIVOTROOT, CHANGEPROFILE,
		MQUEUE, IOURING, UNIX, PTRACE, SIGNAL, DBUS,
		FILE, LINK,
	})

	// convert apparmor requested mask to apparmor access mode
	maskToAccess = map[string]string{
		"a":  "w",
		"c":  "w",
		"d":  "w",
		"wc": "w",
		"x":  "ix",
	}

	// The order the apparmor rules should be sorted
	ruleAlphabet = []Kind{
		INCLUDE,
		ALL,
		RLIMIT,
		USERNS,
		CAPABILITY,
		NETWORK,
		MOUNT,
		REMOUNT,
		UMOUNT,
		PIVOTROOT,
		CHANGEPROFILE,
		MQUEUE,
		IOURING,
		SIGNAL,
		PTRACE,
		UNIX,
		DBUS,
		FILE,
		LINK,
		PROFILE,
		HAT,
		"include_if_exists",
	}
	ruleWeights = generateWeights(ruleAlphabet)

	// The order the apparmor file rules should be sorted. Some rules are sorted
	// in the same group and their order is determined by fileWeights.
	fileGroups = map[string][]string{
		// 1. entry point
		"attachment": {
			`@{exec_path}`,
		},
		// 2 Binaries
		"bin": {
			`@{sh_path}`,        // Shells
			`@{coreutils_path}`, // Coreutils
			`@{open_path}`,      // Binaries paths
			`@{bin}`,            // Binaries
		},
		// 3 Libraries
		"lib": {
			`@{lib}`,      // Libraries
			`@{lib_dirs}`, // Profile libraries
			`/opt`,        // opt binaries & libraries
		},
		// 4. System shared data
		"share": {
			`@{system_share_dirs}`,
			`/usr/share`,
		},
		// 5. System configuration
		"etc": {
			`/etc`, `@{etc_ro}`, `@{etc_rw}`,
			`/etc/machine-id$`,
			`/var/lib/dbus/machine-id`,
		},
		// 6. Boot files
		"boot": {
			`/boot`,
		},
		// 7. System read/write data
		"system-data": {
			`/$`,
			`/usr/`,
			`/usr/local/`,
			`/home`,
			`/var`,
		},
		// 8. System user data
		"system-user": {
			`@{(DESKTOP|GDM|SSDM|LIGHTDM)_HOME}`,
			`@{(desktop|gdm|ssdm|lightdm)_[a-z]*_dirs}`,
		},
		// 9. User data
		"user-data": {
			`@{MOUNTDIRS}`,
			`@{MOUNTS}`,
			`@{HOME}`,
			`@{[a-z]*_dirs}`,
			`@{user_[a-z]*_dirs}`,
		},
		// 10 Temporary data
		"tmp": {
			`/tmp`, `@{tmp}`, // Temporary data
			`/dev/shm`, // Shared memory
		},
		// 11 Runtime data
		"runtime": {
			`@{run}/user/@{uid}`,
			`@{run}/gdm`,
			`@{run}`,
		},
		// 12 Udev data
		"udev": {"@{run}/udev"},
		// 13. Sys files
		"sys": {"@{sys}"},
		// 14. Proc files
		"proc": {"@{PROC}"},
		// 15. Dev files
		"dev": {"/dev"},
	}
	fileAlphabetGroups = util.InvertFlatten(fileGroups)

	// The order the apparmor file group rules should be sorted
	fileAlphabetGroup = []string{
		"attachment", "bin", "lib", "share", "etc", "boot",
		"system-data", "system-user", "user-data",
		"tmp", "runtime", "udev", "sys", "proc", "dev",
	}

	fileWeights = generateFileWeights(fileAlphabetGroup, fileGroups)

	// Compiled regexps for matching file paths to their sort group
	fileReg = generateRegexp(fileWeights)

	// Memoization cache for file path to group label lookups
	groupCache = map[string]string{}

	// The order AARE should be sorted
	stringAlphabet = []byte(
		"!\"#$%&'*(){}[]@+,-./:;<=>?\\^_`|~0123456789abcdefghijklmnopqrstuvwxyz",
	)
	stringWeights = generateWeights(stringAlphabet)

	// The order the rule values (access, type, domains, etc) should be sorted
	requirements        = map[Kind]requirement{}
	requirementsWeights map[Kind]map[string]map[string]int

	// Pairs of mutually exclusive values that cannot coexist
	conflicts = map[Kind]map[string][][]string{}
)

func init() {
	requirementsWeights = generateRequirementsWeights(requirements)
}

// generateTemplates parses and clones the embedded gotmpl files for each rule kind.
func generateTemplates(names []Kind) map[Kind]*template.Template {
	res := make(map[Kind]*template.Template, len(names))
	base := template.New("").Funcs(tmplFunctionMap)
	base = template.Must(base.ParseFS(tmplFiles,
		"templates/*.gotmpl", "templates/rule/*.gotmpl",
	))
	for _, name := range names {
		t := template.Must(base.Clone())
		t = template.Must(t.Parse(
			fmt.Sprintf(`{{- template "%s" . -}}`, name),
		))
		res[name] = t
	}
	return res
}

// renderTemplate render the template for the given kind and data
func renderTemplate(name Kind, data any) string {
	var res strings.Builder
	template, ok := tmpl[name]
	if !ok {
		panic("template '" + name.String() + "' not found")
	}
	err := template.Execute(&res, data)
	if err != nil {
		panic(err)
	}
	return trimTrailingWhitespace(res.String())
}

// trimTrailingWhitespace removes trailing spaces and tabs from each line.
func trimTrailingWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}

// generateWeights assigns a sort weight to each element based on its position in the alphabet.
func generateWeights[T comparable](alphabet []T) map[T]int {
	res := make(map[T]int, len(alphabet))
	for i, r := range alphabet {
		res[r] = i
	}
	return res
}

// generateFileWeights assigns sort weights to file patterns based on their group order.
func generateFileWeights[T comparable](groupOrder []T, groups map[T][]T) map[T]int {
	totalLen := 0
	for i := range groups {
		totalLen += len(groups[i])
	}

	idx := 0
	res := make(map[T]int, totalLen)
	for _, group := range groupOrder {
		for _, r := range groups[group] {
			res[r] = idx
			idx++
		}
	}
	return res
}

// generateRequirementsWeights builds per-kind, per-key weight maps for sorting rule values.
func generateRequirementsWeights(requirements map[Kind]requirement) map[Kind]map[string]map[string]int {
	res := make(map[Kind]map[string]map[string]int, len(requirements))
	for rule, req := range requirements {
		res[rule] = make(map[string]map[string]int, len(req))
		for key, values := range req {
			res[rule][key] = generateWeights(values)
		}
	}
	return res
}

// generateRegexp compiles AARE file patterns into anchored regular expressions.
func generateRegexp(weights map[string]int) map[string]*regexp.Regexp {
	res := make(map[string]*regexp.Regexp, len(weights))
	for w := range weights {
		// Escape special regex chars that appear in AARE patterns
		// Note: $ at end of pattern is kept as regex end anchor for exact matching
		pattern := w
		pattern = strings.ReplaceAll(pattern, "{", `\{`)
		pattern = strings.ReplaceAll(pattern, "}", `\}`)
		// Only escape $ if not at end of pattern (end anchor)
		if !strings.HasSuffix(pattern, "$") {
			pattern = strings.ReplaceAll(pattern, "$", `\$`)
		}
		res[w] = regexp.MustCompile(`(?m)^` + pattern)
	}
	return res
}

// getGroup returns the file sort group label for the given path, using cached results.
func getGroup(weights map[string]int, in string) string {
	if result, ok := groupCache[in]; ok {
		return result
	}

	// Find the best (most specific) matching pattern
	// More specific patterns have higher weights within their group
	bestLabel := ""
	bestWeight := -1
	for w := range weights {
		if fileReg[w].MatchString(in) {
			// Choose the pattern with the highest weight (most specific)
			if weights[w] > bestWeight {
				bestWeight = weights[w]
				bestLabel = w
			}
		}
	}
	groupCache[in] = bestLabel
	return bestLabel
}

// join is a template function that joins slices with spaces or maps as key=value pairs.
func join(i any) string {
	switch i := i.(type) {
	case []string:
		return strings.Join(i, " ")
	case map[string]string:
		res := []string{}
		for k, v := range i {
			res = append(res, k+"="+v)
		}
		return strings.Join(res, " ")
	default:
		return i.(string)
	}
}

// cjoin is a template function that joins values with parentheses when there are multiple items.
func cjoin(i any) string {
	switch i := i.(type) {
	case []string:
		if len(i) == 1 {
			return i[0]
		}
		return "(" + strings.Join(i, " ") + ")"
	case map[string]string:
		res := []string{}
		for k, v := range i {
			res = append(res, k+"="+v)
		}
		return "(" + strings.Join(res, " ") + ")"
	default:
		return i.(string)
	}
}

// kindOf is a template function that returns the kind string of a rule.
func kindOf(i Rule) string {
	if i == nil {
		return ""
	}
	return i.Kind().String()
}

// setindent is a template function that increments or decrements the indentation level.
func setindent(i string) string {
	switch i {
	case "++":
		IndentationLevel++
	case "--":
		IndentationLevel--
	}
	return ""
}

// indent is a template function that prepends the current indentation to a string.
func indent(s string) string {
	return strings.Repeat(Indentation, IndentationLevel) + s
}

// indentDbus is a template function that indents dbus rule continuation lines with extra padding.
func indentDbus(s string) string {
	return strings.Join([]string{Indentation, s}, "     ")
}
