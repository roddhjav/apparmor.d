// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"fmt"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

type Package struct {
	Name     string
	Mode     string
	Required []string
	Profiles []string
	Ignores  []string
	Ignored  []string
	builddir *paths.Path
}

func NewPackage(name string) *Package {
	path := PkgDir.Join(name + ".conf")
	if !path.Exist() {
		panic(fmt.Sprintf("Unknown package: %s", name))
	}
	lines := path.MustReadFilteredFileAsLines()
	mode := ""
	profiles := make([]string, 0, len(lines))
	ignores := []string{}
	dependencies := []string{}
	ignored := getFilesIgnoredByDistribution()
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "mode="):
			mode = strings.TrimPrefix(line, "mode=")
		case strings.HasPrefix(line, "require="):
			dependencies = strings.Split(strings.TrimPrefix(line, "require="), ",")
		case strings.HasPrefix(line, "!"):
			ignores = append(ignores, strings.TrimPrefix(line, "!"))
		default:
			profiles = append(profiles, line)
		}
	}
	return &Package{
		Name:     name,
		Mode:     mode,
		Required: dependencies,
		Profiles: profiles,
		Ignores:  ignores,
		Ignored:  ignored,
		builddir: Root.Join(name),
	}
}

func getFilesIgnoredByDistribution() []string {
	res := []string{}
	for _, iname := range []string{"main", Distribution} {
		for _, ignore := range Ignore.Read(iname) {
			if !strings.HasPrefix(ignore, Src) {
				continue
			}
			profile := strings.TrimPrefix(ignore, Src+"/")
			path := SrcApparmord.Join(profile)
			if path.IsDir() {
				files, err := path.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
				if err != nil {
					panic(err)
				}
				for _, file := range files {
					res = append(res, file.Base())
				}
			} else if path.Exist() {
				res = append(res, path.Base())
			} else {
				panic(fmt.Errorf("%s.ignore: no files found for '%s'", iname, profile))
			}
		}
	}
	return res
}

func (p *Package) Generate() ([]string, error) {
	var res []string

	if err := p.builddir.RemoveAll(); err != nil {
		return res, err
	}
	if err := p.builddir.MkdirAll(); err != nil {
		return res, err
	}

	explode := paths.PathList{
		paths.New("groups"), paths.New("profiles-a-f"),
		paths.New("profiles-m-r"), paths.New("profiles-s-z"),
	}
	for _, name := range p.Profiles {
		originalPath := SrcApparmord.Join(name)

		if originalPath.IsDir() {
			originFiles, err := originalPath.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
			if err != nil {
				return res, err
			}
			for _, originFile := range originFiles {
				file, err := originFile.RelFrom(SrcApparmord)
				if err != nil {
					return res, err
				}

				if slices.Contains(p.Ignores, file.String()) {
					continue
				}

				done := false
				for _, e := range explode {
					if ok, _ := file.IsInsideDir(e); ok {
						base := file.Base()
						msg, err := p.move(base)
						if err != nil {
							return res, err
						}
						res = append(res, msg)
						done = true
						break
					}
				}

				if !done {
					msg, err := p.move(file)
					if err != nil {
						return res, err
					}
					res = append(res, msg)
				}
			}

		} else if originalPath.Exist() {
			base := originalPath.Base()
			if slices.Contains(p.Ignores, base) {
				continue
			}
			msg, err := p.move(base)
			if err != nil {
				return res, err
			}
			res = append(res, msg)

		} else {
			return res, fmt.Errorf("No %s", originalPath)
		}
	}
	return res, nil
}

func (p *Package) move(origin any) (string, error) {
	var src *paths.Path
	var dst *paths.Path
	var srcOverridden *paths.Path
	var dstOverridden *paths.Path
	var srcSymlink *paths.Path
	var dstSymlink *paths.Path
	const ext = ".apparmor.d"

	switch value := any(origin).(type) {
	case string:
		src = RootApparmord.Join(value)
		dst = p.builddir.Join(value)
		srcOverridden = RootApparmord.Join(value + ext)
		dstOverridden = p.builddir.Join(value + ext)
		srcSymlink = RootApparmord.Join("disable", value)
		dstSymlink = p.builddir.Join("disable", value)

	case *paths.Path:
		src = RootApparmord.JoinPath(value)
		dst = p.builddir.JoinPath(value)
		srcOverridden = RootApparmord.JoinPath(value.Parent()).Join(value.Base() + ext)
		dstOverridden = p.builddir.JoinPath(value.Parent()).Join(value.Base() + ext)
		srcSymlink = RootApparmord.Join("disable").JoinPath(value)
		dstSymlink = p.builddir.Join("disable").JoinPath(value)

	default:
		panic("Package.move: unsupported type")
	}

	if src.Exist() {
		if err := dst.Parent().MkdirAll(); err != nil {
			return "", nil
		}
		if err := src.Rename(dst); err != nil {
			return "", nil
		}
		// fmt.Printf("%s -> %s\n", src, dst)

	} else if srcOverridden.Exist() {
		if err := dstOverridden.Parent().MkdirAll(); err != nil {
			return "", nil
		}
		if err := dstSymlink.Parent().MkdirAll(); err != nil {
			return "", nil
		}
		if err := srcOverridden.Rename(dstOverridden); err != nil {
			return "", nil
		}
		if err := srcSymlink.Rename(dstSymlink); err != nil {
			return "", nil
		}
		// fmt.Printf("%s -> %s\n", srcOverridden, dstOverridden)

	} else {
		srcRltv, err := src.RelFrom(RootApparmord)
		if err != nil {
			return "", nil
		}
		if !slices.Contains(p.Ignored, srcRltv.String()) {
			fmt.Printf("Warning: No %s\n", src)
			// return "", fmt.Errorf("No %s", src)
		}

	}
	return "", nil
}

// Validate ensures a package has its required dependencies
func (p *Package) Validate() error {
	return nil
}

