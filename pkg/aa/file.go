// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"strings"
)

const (
	LINK      Kind = "link"
	FILE      Kind = "file"
	tokOWNER       = "owner"
	tokSUBSET      = "subset"
)

var fileAccessBits = map[string]AccessMask{
	// access
	"m": 1 << 0,
	"r": 1 << 1,
	"w": 1 << 2,
	"a": 1 << 3,
	"l": 1 << 4,
	"k": 1 << 5,
	// transitions
	"ix":  1 << 6,
	"ux":  1 << 7,
	"Ux":  1 << 8,
	"px":  1 << 9,
	"Px":  1 << 10,
	"cx":  1 << 11,
	"Cx":  1 << 12,
	"pix": 1 << 13,
	"Pix": 1 << 14,
	"cix": 1 << 15,
	"Cix": 1 << 16,
	"pux": 1 << 17,
	"PUx": 1 << 18,
	"cux": 1 << 19,
	"CUx": 1 << 20,
	"x":   1 << 21,
	"Pux": 1 << 22,
	"pUx": 1 << 23,
}

func init() {
	requirements[FILE] = requirement{
		"access": {"m", "r", "w", "a", "l", "k"},
		"transition": {
			"ix", "ux", "Ux", "px", "Px", "cx", "Cx", "pix", "Pix", "cix",
			"Cix", "pux", "PUx", "cux", "CUx", "x",
			"Pux", "pUx",
		},
	}
}

func IsOwner(log map[string]string) bool {
	fsuid, hasFsUID := log["fsuid"]
	ouid, hasOuUID := log["ouid"]
	isDbus := strings.Contains(log["operation"], "dbus")
	if hasFsUID && hasOuUID && fsuid == ouid && ouid != "0" && !isDbus {
		return true
	}
	return false
}

type File struct {
	Base
	Qualifier
	Owner  bool
	Path   string
	Access AccessMask
	Target string
}

func newFile(q Qualifier, rule rule) (Rule, error) {
	path, access, target, owner := "", "", "", false
	if len(rule) > 0 {
		if rule.Get(0) == tokOWNER {
			owner = true
			rule = rule[1:]
		}
		if rule.Get(0) == FILE.Tok() {
			rule = rule[1:]
		}
		// Skip safe/unsafe modifiers (handled at exec time)
		if rule.Get(0) == "unsafe" || rule.Get(0) == "safe" {
			rule = rule[1:]
		}

		r := rule.GetSlice()
		size := len(r)
		if size < 2 {
			return nil, fmt.Errorf("missing file or access in rule: %s", rule)
		}

		// Determine format: "path access" vs "access path"
		// Try parsing first token as access - if valid, use "access path" format
		if testAccess, _ := toAccess(FILE, r[0]); testAccess != 0 {
			access, path = r[0], r[1]
		} else {
			path, access = r[0], r[1]
		}
		if size > 2 {
			if r[2] != tokARROW {
				return nil, fmt.Errorf("missing '%s' in rule: %s", tokARROW, rule)
			}
			target = r[3]
		}
	}
	accesses, err := toAccess(FILE, access)
	if err != nil {
		return nil, err
	}
	return &File{
		Base:      newBase(rule),
		Qualifier: q,
		Owner:     owner,
		Path:      path,
		Access:    accesses,
		Target:    target,
	}, rule.ValidateMapKeys([]string{})
}

func newFileFromLog(log map[string]string) Rule {
	if log["operation"] == "link" {
		log["requested_mask"] += "l"
	}
	accesses, err := toAccess("file-log", log["requested_mask"])
	if err != nil {
		panic(fmt.Errorf("newFileFromLog(%v): %w", log, err))
	}
	if accesses == accessBits[FILE]["l"] {
		return newLinkFromLog(log)
	}
	return &File{
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Owner:     IsOwner(log),
		Path:      log["name"],
		Access:    accesses,
		Target:    log["target"],
	}
}

func (r *File) mergeKey(b *strings.Builder) {
	writeQualifierKey(b, r.Qualifier)
	writeBool(b, r.Owner)
	b.WriteString(r.Path)
	b.WriteByte(0)
	b.WriteString(r.Target)
}

func (r *File) Kind() Kind {
	return FILE
}

func (r *File) AccessStrings() []string {
	return r.Access.Strings(FILE)
}

func (r *File) Constraint() Constraint {
	return BlockRule
}

func (r *File) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *File) Validate() error {
	if r.Path == "" && r.Target == "" && r.Access == 0 {
		return nil // rule: `file` or `owner file`
	}
	if !isAARE(r.Path) {
		return fmt.Errorf("'%s' is not a valid AARE", r.Path)
	}
	if r.Access == 0 {
		return fmt.Errorf("missing file access")
	}
	if err := validateAccess(r.Kind(), r.Access); err != nil {
		return err
	}
	if err := validateAAREPattern(r.Path); err != nil {
		return err
	}
	// Conflicting access: write (w) and append (a) cannot coexist
	hasW := r.Access.Has(accessBits[FILE]["w"])
	hasA := r.Access.Has(accessBits[FILE]["a"])
	if hasW && hasA {
		return fmt.Errorf("conflicting file access: 'w' and 'a' cannot coexist")
	}
	return nil
}

func (r *File) Compare(other Rule) int {
	o, _ := other.(*File)
	if o.AccessType == "deny" {
		return -1 // Deny file rules always come last
	}

	// Compare by file group - use pattern matching
	groupR := getGroup(fileWeights, r.Path)
	groupO := getGroup(fileWeights, o.Path)
	if groupR != "" && groupO != "" {
		weightR := fileWeights[groupR]
		weightO := fileWeights[groupO]
		if weightR != weightO {
			return weightR - weightO
		}
	} else if groupR != "" {
		return -1
	} else if groupO != "" {
		return 1
	}

	if res := r.Qualifier.Compare(o.Qualifier); res != 0 {
		return res
	}
	if res := compare(r.Owner, o.Owner); res != 0 {
		return res
	}
	if res := compare(r.Path, o.Path); res != 0 {
		return res
	}
	if res := compareAccessMask(r.Access, o.Access, FILE); res != 0 {
		return res
	}
	return compare(r.Target, o.Target)
}

func (r *File) Merge(other Rule) bool {
	o, _ := other.(*File)

	if !r.Equal(o.Qualifier) {
		return false
	}
	if r.Owner == o.Owner && r.Path == o.Path && r.Target == o.Target {
		r.Access |= o.Access
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *File) Lengths() []int {
	// Deny rules don't participate in padding alignment
	if r.AccessType == "deny" {
		return []int{0, 0, 0, 0}
	}

	// Add padding to align with other transition rule
	lenPath := 0
	ensureAccessMasks()
	if r.Access.ContainsAny(fileTransitionMask) {
		lenPath = length("", r.Path)
	}
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("owner", r.Owner),
		lenPath,
	}
}

func (r *File) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"owner", ""},
		[]any{r.Owner, r.Path})...,
	)
}

func (r *File) addLine(other Rule) bool {
	if other.Kind() != r.Kind() {
		return false
	}
	o := other.(*File)

	// Deny rules are all grouped together without blank lines
	if r.AccessType == "deny" && o.AccessType == "deny" {
		return false
	}

	patternI := getGroup(fileWeights, r.Path)
	patternJ := getGroup(fileWeights, o.Path)
	if patternI == "" || patternJ == "" {
		return patternI != patternJ
	}
	groupI, ok1 := fileAlphabetGroups[patternI]
	groupJ, ok2 := fileAlphabetGroups[patternJ]

	// Add newline if patterns differ and they're in different groups (or unrecognized)
	return patternI != patternJ && (!ok1 || !ok2 || groupI != groupJ)
}

type Link struct {
	Base
	Qualifier
	Owner  bool
	Subset bool
	Path   string
	Target string
}

func newLink(q Qualifier, rule rule) (Rule, error) {
	owner, subset, path, target := false, false, "", ""
	if len(rule) > 0 {
		if rule.Get(0) == tokOWNER {
			owner = true
			rule = rule[1:]
		}
		if len(rule) > 0 && rule.Get(0) == tokSUBSET {
			subset = true
			rule = rule[1:]
		}

		r := rule.GetSlice()
		size := len(r)
		if size > 0 {
			path = r[0]
		}
		if size > 2 {
			if r[1] != tokARROW {
				return nil, fmt.Errorf("missing '%s' in rule: %s", tokARROW, rule)
			}
			target = r[2]
		}
	}
	return &Link{
		Base:      newBase(rule),
		Qualifier: q,
		Owner:     owner,
		Subset:    subset,
		Path:      path,
		Target:    target,
	}, rule.ValidateMapKeys([]string{})
}

func newLinkFromLog(log map[string]string) Rule {
	return &Link{
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Owner:     IsOwner(log),
		Path:      log["name"],
		Target:    log["target"],
	}
}

func (r *Link) mergeKey(b *strings.Builder) {
	writeQualifierKey(b, r.Qualifier)
	writeBool(b, r.Owner)
	writeBool(b, r.Subset)
	b.WriteString(r.Path)
	b.WriteByte(0)
	b.WriteString(r.Target)
}

func (r *Link) Kind() Kind {
	return LINK
}

func (r *Link) Constraint() Constraint {
	return BlockRule
}

func (r *Link) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Link) Validate() error {
	if !isAARE(r.Path) {
		return fmt.Errorf("'%s' is not a valid AARE", r.Path)
	}
	return nil
}

func (r *Link) Compare(other Rule) int {
	o, _ := other.(*Link)
	if o.AccessType == "deny" {
		return -1 // Deny file rules always come last
	}

	// Compare by file group - use pattern matching
	groupR := getGroup(fileWeights, r.Path)
	groupO := getGroup(fileWeights, o.Path)
	if groupR != "" && groupO != "" {
		weightR := fileWeights[groupR]
		weightO := fileWeights[groupO]
		if weightR != weightO {
			return weightR - weightO
		}
	} else if groupR != "" {
		return -1
	} else if groupO != "" {
		return 1
	}

	if res := compare(r.Owner, o.Owner); res != 0 {
		return res
	}
	if res := compare(r.Path, o.Path); res != 0 {
		return res
	}
	if res := compare(r.Target, o.Target); res != 0 {
		return res
	}
	if res := compare(r.Subset, o.Subset); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Link) Merge(other Rule) bool {
	return false // Never merge link
}

func (r *Link) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("owner", r.Owner),
		length("subset", r.Subset),
		length("", r.Path),
		length("", r.Target),
	}
}

func (r *Link) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"owner", "subset", "", ""},
		[]any{r.Owner, r.Subset, r.Path, r.Target})...,
	)
}

// compareFileLink compares File and Link rules by their file group weight.
func compareFileLink(a, b Rule) int {
	pathA := ""
	switch r := a.(type) {
	case *File:
		pathA = r.Path
	case *Link:
		pathA = r.Path
	}

	pathB := ""
	switch r := b.(type) {
	case *File:
		pathB = r.Path
	case *Link:
		pathB = r.Path
	}

	groupA := getGroup(fileWeights, pathA)
	groupB := getGroup(fileWeights, pathB)
	if groupA != "" && groupB != "" {
		if res := fileWeights[groupA] - fileWeights[groupB]; res != 0 {
			return res
		}
	} else if groupA != "" {
		return -1
	} else if groupB != "" {
		return 1
	}
	return compare(pathA, pathB)
}
