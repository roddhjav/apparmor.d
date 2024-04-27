/*
 * This file is part of PathsHelper library.
 *
 * Copyright 2018-2022 Arduino AG (http://www.arduino.cc/)
 *
 * PathsHelper library is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 */

package paths

import (
	"errors"
	"os"
	"strings"
)

// ReadDirFilter is a filter for Path.ReadDir and Path.ReadDirRecursive methods.
// The filter should return true to accept a file or false to reject it.
type ReadDirFilter func(file *Path) bool

// ReadDir returns a PathList containing the content of the directory
// pointed by the current Path. The resulting list is filtered by the given filters chained.
func (p *Path) ReadDir(filters ...ReadDirFilter) (PathList, error) {
	infos, err := os.ReadDir(p.path)
	if err != nil {
		return nil, err
	}

	accept := func(p *Path) bool {
		for _, filter := range filters {
			if !filter(p) {
				return false
			}
		}
		return true
	}

	paths := PathList{}
	for _, info := range infos {
		path := p.Join(info.Name())
		if !accept(path) {
			continue
		}
		paths.Add(path)
	}
	return paths, nil
}

// ReadDirRecursive returns a PathList containing the content of the directory
// and its subdirectories pointed by the current Path
func (p *Path) ReadDirRecursive() (PathList, error) {
	return p.ReadDirRecursiveFiltered(nil)
}

// ReadDirRecursiveFiltered returns a PathList containing the content of the directory
// and its subdirectories pointed by the current Path, filtered by the given skipFilter
// and filters:
//   - `recursionFilter` is a filter that is checked to determine if the subdirectory must
//     by visited recursively (if the filter rejects the entry, the entry is not visited
//     but can still be added to the result)
//   - `filters` are the filters that are checked to determine if the entry should be
//     added to the resulting PathList
func (p *Path) ReadDirRecursiveFiltered(recursionFilter ReadDirFilter, filters ...ReadDirFilter) (PathList, error) {
	var search func(*Path) (PathList, error)

	explored := map[string]bool{}
	search = func(currPath *Path) (PathList, error) {
		canonical := currPath.Canonical().path
		if explored[canonical] {
			return nil, errors.New("directories symlink loop detected")
		}
		explored[canonical] = true
		defer delete(explored, canonical)

		infos, err := os.ReadDir(currPath.path)
		if err != nil {
			return nil, err
		}

		accept := func(p *Path) bool {
			for _, filter := range filters {
				if !filter(p) {
					return false
				}
			}
			return true
		}

		paths := PathList{}
		for _, info := range infos {
			path := currPath.Join(info.Name())

			if accept(path) {
				paths.Add(path)
			}

			if recursionFilter == nil || recursionFilter(path) {
				if path.IsDir() {
					subPaths, err := search(path)
					if err != nil {
						return nil, err
					}
					paths.AddAll(subPaths)
				}
			}
		}
		return paths, nil
	}

	return search(p)
}

// FilterDirectories is a ReadDirFilter that accepts only directories
func FilterDirectories() ReadDirFilter {
	return func(path *Path) bool {
		return path.IsDir()
	}
}

// FilterOutDirectories is a ReadDirFilter that rejects all directories
func FilterOutDirectories() ReadDirFilter {
	return func(path *Path) bool {
		return !path.IsDir()
	}
}

// FilterNames is a ReadDirFilter that accepts only the given filenames
func FilterNames(allowedNames ...string) ReadDirFilter {
	return func(file *Path) bool {
		for _, name := range allowedNames {
			if file.Base() == name {
				return true
			}
		}
		return false
	}
}

// FilterOutNames is a ReadDirFilter that rejects the given filenames
func FilterOutNames(rejectedNames ...string) ReadDirFilter {
	return func(file *Path) bool {
		for _, name := range rejectedNames {
			if file.Base() == name {
				return false
			}
		}
		return true
	}
}

// FilterSuffixes creates a ReadDirFilter that accepts only the given
// filename suffixes
func FilterSuffixes(allowedSuffixes ...string) ReadDirFilter {
	return func(file *Path) bool {
		for _, suffix := range allowedSuffixes {
			if strings.HasSuffix(file.String(), suffix) {
				return true
			}
		}
		return false
	}
}

// FilterOutSuffixes creates a ReadDirFilter that rejects all the given
// filename suffixes
func FilterOutSuffixes(rejectedSuffixes ...string) ReadDirFilter {
	return func(file *Path) bool {
		for _, suffix := range rejectedSuffixes {
			if strings.HasSuffix(file.String(), suffix) {
				return false
			}
		}
		return true
	}
}

// FilterPrefixes creates a ReadDirFilter that accepts only the given
// filename prefixes
func FilterPrefixes(allowedPrefixes ...string) ReadDirFilter {
	return func(file *Path) bool {
		name := file.Base()
		for _, prefix := range allowedPrefixes {
			if strings.HasPrefix(name, prefix) {
				return true
			}
		}
		return false
	}
}

// FilterOutPrefixes creates a ReadDirFilter that rejects all the given
// filename prefixes
func FilterOutPrefixes(rejectedPrefixes ...string) ReadDirFilter {
	return func(file *Path) bool {
		name := file.Base()
		for _, prefix := range rejectedPrefixes {
			if strings.HasPrefix(name, prefix) {
				return false
			}
		}
		return true
	}
}

// OrFilter creates a ReadDirFilter that accepts all items that are accepted
// by any (at least one) of the given filters
func OrFilter(filters ...ReadDirFilter) ReadDirFilter {
	return func(path *Path) bool {
		for _, f := range filters {
			if f(path) {
				return true
			}
		}
		return false
	}
}

// AndFilter creates a ReadDirFilter that accepts all items that are accepted
// by all the given filters
func AndFilter(filters ...ReadDirFilter) ReadDirFilter {
	return func(path *Path) bool {
		for _, f := range filters {
			if !f(path) {
				return false
			}
		}
		return true
	}
}

// NotFilter creates a ReadDifFilter that accepts all items rejected by x and viceversa
func NotFilter(x ReadDirFilter) ReadDirFilter {
	return func(path *Path) bool {
		return !x(path)
	}
}
