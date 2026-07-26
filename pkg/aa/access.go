// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"sync"
)

type AccessMask uint32

var (
	// accessBits maps each rule kind to its access token bit assignments.
	// The per-kind map literals live in their rule file alongside requirements.
	// Bits may be reused across kinds because they live on different rule types.
	accessBits = map[Kind]map[string]AccessMask{
		FILE:    fileAccessBits,
		NETWORK: networkAccessBits,
		DBUS:    dbusAccessBits,
		UNIX:    unixAccessBits,
		PTRACE:  ptraceAccessBits,
		SIGNAL:  signalAccessBits,
		MQUEUE:  mqueueAccessBits,
		IOURING: ioUringAccessBits,
	}

	// Lazily-derived masks. Built on first use because they read requirements,
	// which each rule file populates in its own init(); a package-level build
	// here could run before those inits.
	accessMasksOnce    sync.Once
	allAccessMasks     map[Kind]AccessMask // every bit defined for the kind
	localOnlyMasks     map[Kind]AccessMask // bits flagged "local-only" in requirements
	fileTransitionMask AccessMask          // file transition bits plus "m"
)

func ensureAccessMasks() {
	accessMasksOnce.Do(func() {
		allAccessMasks = make(map[Kind]AccessMask, len(accessBits))
		localOnlyMasks = make(map[Kind]AccessMask, len(accessBits))
		for kind, bits := range accessBits {
			var all, lo AccessMask
			for _, b := range bits {
				all |= b
			}
			for _, tok := range requirements[kind]["local-only"] {
				lo |= bits[tok]
			}
			allAccessMasks[kind] = all
			localOnlyMasks[kind] = lo
		}
		// File transition mask: transition tokens plus "m", mirroring file.go.
		fbits := accessBits[FILE]
		fileTransitionMask = fbits["m"]
		for _, tok := range requirements[FILE]["transition"] {
			fileTransitionMask |= fbits[tok]
		}
	})
}

// toAccessMask converts a slice of access tokens into a bitmask for the given kind.
// Returns an error if any token is not defined for the kind.
func toAccessMask(kind Kind, tokens []string) (AccessMask, error) {
	bits, ok := accessBits[kind]
	if !ok {
		return 0, fmt.Errorf("unknown rule kind for access: %s", kind)
	}
	var m AccessMask
	for _, t := range tokens {
		if t == "" {
			continue
		}
		b, ok := bits[t]
		if !ok {
			return 0, fmt.Errorf("unrecognized access for rule %s: %s", kind, t)
		}
		m |= b
	}
	return m, nil
}

// Strings returns the tokens set in the mask, in the canonical order held by
// requirements[kind] (access tokens, then file transition tokens).
func (m AccessMask) Strings(kind Kind) []string {
	if m == 0 {
		return nil
	}
	bits, ok := accessBits[kind]
	if !ok {
		return nil
	}
	var res []string
	for _, group := range [...]string{"access", "transition"} {
		for _, tok := range requirements[kind][group] {
			if m&bits[tok] != 0 {
				res = append(res, tok)
			}
		}
	}
	return res
}

// Has reports whether m has bit set.
func (m AccessMask) Has(bit AccessMask) bool {
	return m&bit != 0
}

// ContainsAny reports whether m has any bit in common with other.
func (m AccessMask) ContainsAny(other AccessMask) bool {
	return m&other != 0
}

// validateAccess checks that m only has bits defined for kind.
func validateAccess(kind Kind, m AccessMask) error {
	ensureAccessMasks()
	all, ok := allAccessMasks[kind]
	if !ok {
		return fmt.Errorf("unknown rule kind for access: %s", kind)
	}
	if extra := m &^ all; extra != 0 {
		return fmt.Errorf("invalid mode '%#x'", uint32(extra))
	}
	return nil
}

// compareAccessMask compares two access masks for the given kind, returning the
// same lexicographic ordering as comparing their rendered string slices. This
// preserves Rules.Sort() ordering across the bitmask refactor. It is only ever
// a sort tie-breaker (files) or runs on a handful of rules (signal, dbus, ...),
// so the two short-lived slices it allocates are not worth hand-optimizing away.
func compareAccessMask(a, b AccessMask, kind Kind) int {
	if a == b {
		return 0
	}
	return compare(a.Strings(kind), b.Strings(kind))
}

// allLocalOnly returns true if every set bit in access corresponds to a
// local-only access token for kind.
func allLocalOnly(access AccessMask, kind Kind) bool {
	ensureAccessMasks()
	lo, ok := localOnlyMasks[kind]
	if !ok {
		return false
	}
	return access&^lo == 0
}

// MustAccess parses a slice of access tokens and panics on error. Intended for
// use in tests where the tokens are known to be valid.
func MustAccess(kind Kind, tokens ...string) AccessMask {
	m, err := toAccessMask(kind, tokens)
	if err != nil {
		panic(err)
	}
	return m
}
