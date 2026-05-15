// This file is part of PathsHelper library.
// Copyright (C) 2018-2025 Arduino AG (http://www.arduino.cc/)
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package paths

import (
	"sort"
)

// PathList is a list of Path
type PathList []*Path

// NewPathList creates a new PathList with the given paths
func NewPathList(paths ...string) PathList {
	res := PathList{}
	for _, path := range paths {
		res = append(res, New(path))
	}
	return res
}

// Clone returns a copy of the current PathList
func (p *PathList) Clone() PathList {
	res := PathList{}
	for _, path := range *p {
		res.Add(path.Clone())
	}
	return res
}

// AsStrings return this path list as a string array
func (p *PathList) AsStrings() []string {
	res := []string{}
	for _, path := range *p {
		res = append(res, path.String())
	}
	return res
}

// Equals returns true if the current PathList is equal to the
// PathList passed as argument
func (p *PathList) Equals(other PathList) bool {
	if len(*p) != len(other) {
		return false
	}
	for i, path := range *p {
		if !path.EqualsTo(other[i]) {
			return false
		}
	}
	return true
}

// FilterDirs remove all entries except directories
func (p *PathList) FilterDirs() {
	res := (*p)[:0]
	for _, path := range *p {
		if path.IsDir() {
			res = append(res, path)
		}
	}
	*p = res
}

// FilterOutDirs remove all directories entries
func (p *PathList) FilterOutDirs() {
	res := (*p)[:0]
	for _, path := range *p {
		if !path.IsDir() {
			res = append(res, path)
		}
	}
	*p = res
}

// FilterOutHiddenFiles remove all hidden files (files with the name
// starting with ".")
func (p *PathList) FilterOutHiddenFiles() {
	p.FilterOutPrefix(".")
}

// Filter will remove all the elements of the list that do not match
// the specified acceptor function
func (p *PathList) Filter(acceptorFunc func(*Path) bool) {
	res := (*p)[:0]
	for _, path := range *p {
		if acceptorFunc(path) {
			res = append(res, path)
		}
	}
	*p = res
}

// FilterOutPrefix remove all entries having one of the specified prefixes
func (p *PathList) FilterOutPrefix(prefixes ...string) {
	filterFunction := func(path *Path) bool {
		return !path.HasPrefix(prefixes...)
	}
	p.Filter(filterFunction)
}

// FilterPrefix remove all entries not having one of the specified prefixes
func (p *PathList) FilterPrefix(prefixes ...string) {
	filterFunction := func(path *Path) bool {
		return path.HasPrefix(prefixes...)
	}
	p.Filter(filterFunction)
}

// FilterOutSuffix remove all entries having one of the specified suffixes
func (p *PathList) FilterOutSuffix(suffixies ...string) {
	filterFunction := func(path *Path) bool {
		return !path.HasSuffix(suffixies...)
	}
	p.Filter(filterFunction)
}

// FilterSuffix remove all entries not having one of the specified suffixes
func (p *PathList) FilterSuffix(suffixies ...string) {
	filterFunction := func(path *Path) bool {
		return path.HasSuffix(suffixies...)
	}
	p.Filter(filterFunction)
}

// Add adds a Path to the PathList
func (p *PathList) Add(path *Path) {
	*p = append(*p, path)
}

// AddAll adds all Paths in the list passed as argument
func (p *PathList) AddAll(paths PathList) {
	*p = append(*p, paths...)
}

// AddIfMissing adds a Path to the PathList if the path is not already
// in the list
func (p *PathList) AddIfMissing(path *Path) {
	if (*p).Contains(path) {
		return
	}
	(*p).Add(path)
}

// AddAllMissing adds all paths to the PathList excluding the paths already
// in the list
func (p *PathList) AddAllMissing(pathsToAdd PathList) {
	for _, pathToAdd := range pathsToAdd {
		(*p).AddIfMissing(pathToAdd)
	}
}

// ToAbs calls Path.ToAbs() method on each path of the list.
// It stops at the first error and returns it. If all ToAbs calls
// are successful nil is returned.
func (p *PathList) ToAbs() error {
	for _, path := range *p {
		if err := path.ToAbs(); err != nil {
			return err
		}
	}
	return nil
}

// Contains check if the list contains a path that match
// exactly (EqualsTo) to the specified path
func (p *PathList) Contains(pathToSearch *Path) bool {
	for _, path := range *p {
		if path.EqualsTo(pathToSearch) {
			return true
		}
	}
	return false
}

// ContainsEquivalentTo check if the list contains a path
// that is equivalent (EquivalentTo) to the specified path
func (p *PathList) ContainsEquivalentTo(pathToSearch *Path) bool {
	for _, path := range *p {
		if path.EquivalentTo(pathToSearch) {
			return true
		}
	}
	return false
}

// Sort sorts this pathlist
func (p *PathList) Sort() {
	sort.Sort(p)
}

func (p *PathList) Len() int {
	return len(*p)
}

func (p *PathList) Less(i, j int) bool {
	return (*p)[i].path < (*p)[j].path
}

func (p *PathList) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}
