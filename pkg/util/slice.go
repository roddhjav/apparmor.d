// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

// RemoveDuplicate filter out all duplicates from a slice. Also filter out empty element.
func RemoveDuplicate[T comparable](inlist []T) []T {
	var empty T
	list := []T{}
	seen := map[T]bool{}
	seen[empty] = true
	for _, item := range inlist {
		if _, ok := seen[item]; !ok {
			seen[item] = true
			list = append(list, item)
		}
	}
	return list
}

// Intersect returns the intersection between two collections.
// From https://github.com/samber/lo
func Intersect[T comparable](list1 []T, list2 []T) []T {
	result := []T{}
	seen := map[T]struct{}{}

	for _, elem := range list1 {
		seen[elem] = struct{}{}
	}

	for _, elem := range list2 {
		if _, ok := seen[elem]; ok {
			result = append(result, elem)
		}
	}

	return result
}

// Flatten returns an array a single level deep.
// From https://github.com/samber/lo
func Flatten[T comparable](collection [][]T) []T {
	totalLen := 0
	for i := range collection {
		totalLen += len(collection[i])
	}

	result := make([]T, 0, totalLen)
	for i := range collection {
		result = append(result, collection[i]...)
	}

	return result
}

// Invert creates a map composed of the inverted keys and values. If map
// contains duplicate values, subsequent values overwrite property assignments
// of previous values.
// Play: https://go.dev/play/p/rFQ4rak6iA1
func Invert[K comparable, V comparable](in map[K]V) map[V]K {
	out := make(map[V]K, len(in))

	for k := range in {
		out[in[k]] = k
	}

	return out
}

func InvertFlatten[V comparable](in map[V][]V) map[V]V {
	out := make(map[V]V, len(in))

	for k := range in {
		for _, v := range in[k] {
			out[v] = k
		}
	}

	return out
}
