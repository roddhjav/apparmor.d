/*
 * This file is part of PathsHelper library.
 *
 * Copyright 2018 Arduino AG (http://www.arduino.cc/)
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
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/roddhjav/apparmor.d/pkg/util"
)

// Path represents a path
type Path struct {
	path string
}

// New creates a new Path object. If path is the empty string
// then nil is returned.
func New(path ...string) *Path {
	if len(path) == 0 {
		return nil
	}
	if len(path) == 1 && path[0] == "" {
		return nil
	}
	res := &Path{path: path[0]}
	if len(path) > 1 {
		return res.Join(path[1:]...)
	}
	return res
}

// NewFromFile creates a new Path object using the path name
// obtained from the File object (see os.File.Name function).
func NewFromFile(file *os.File) *Path {
	return New(file.Name())
}

// Stat returns a FileInfo describing the named file. The result is
// cached internally for next queries. To ensure that the cached
// FileInfo entry is updated just call Stat again.
func (p *Path) Stat() (fs.FileInfo, error) {
	return os.Stat(p.path)
}

// Lstat returns a FileInfo describing the named file. If the file is
// a symbolic link, the returned FileInfo describes the symbolic link.
// Lstat makes no attempt to follow the link. If there is an error, it
// will be of type *PathError.
func (p *Path) Lstat() (fs.FileInfo, error) {
	return os.Lstat(p.path)
}

// Clone create a copy of the Path object
func (p *Path) Clone() *Path {
	return New(p.path)
}

// Join create a new Path by joining the provided paths
func (p *Path) Join(paths ...string) *Path {
	return New(filepath.Join(p.path, filepath.Join(paths...)))
}

// JoinPath create a new Path by joining the provided paths
func (p *Path) JoinPath(paths ...*Path) *Path {
	res := p.Clone()
	for _, path := range paths {
		res = res.Join(path.path)
	}
	return res
}

// Base returns the last element of path
func (p *Path) Base() string {
	return filepath.Base(p.path)
}

// Ext returns the file name extension used by path
func (p *Path) Ext() string {
	return filepath.Ext(p.path)
}

// HasPrefix returns true if the file name has one of the
// given prefixes (the Base() method is used to obtain the
// file name used for the comparison)
func (p *Path) HasPrefix(prefixes ...string) bool {
	filename := p.Base()
	for _, prefix := range prefixes {
		if strings.HasPrefix(filename, prefix) {
			return true
		}
	}
	return false
}

// HasSuffix returns true if the file name has one of the
// given suffixies
func (p *Path) HasSuffix(suffixies ...string) bool {
	filename := p.String()
	for _, suffix := range suffixies {
		if strings.HasSuffix(filename, suffix) {
			return true
		}
	}
	return false
}

// RelTo returns a relative Path that is lexically equivalent to r when
// joined to the current Path.
//
// For example paths.New("/my/path/ab/cd").RelTo(paths.New("/my/path"))
// returns "../..".
func (p *Path) RelTo(r *Path) (*Path, error) {
	rel, err := filepath.Rel(p.path, r.path)
	if err != nil {
		return nil, err
	}
	return New(rel), nil
}

// RelFrom returns a relative Path that when joined with r is lexically
// equivalent to the current path.
//
// For example paths.New("/my/path/ab/cd").RelFrom(paths.New("/my/path"))
// returns "ab/cd".
func (p *Path) RelFrom(r *Path) (*Path, error) {
	rel, err := filepath.Rel(r.path, p.path)
	if err != nil {
		return nil, err
	}
	return New(rel), nil
}

// Abs returns the absolute path of the current Path
func (p *Path) Abs() (*Path, error) {
	abs, err := filepath.Abs(p.path)
	if err != nil {
		return nil, err
	}
	return New(abs), nil
}

// IsAbs returns true if the Path is absolute
func (p *Path) IsAbs() bool {
	return filepath.IsAbs(p.path)
}

// ToAbs transform the current Path to the corresponding absolute path
func (p *Path) ToAbs() error {
	abs, err := filepath.Abs(p.path)
	if err != nil {
		return err
	}
	p.path = abs
	return nil
}

// Clean Clean returns the shortest path name equivalent to path by
// purely lexical processing
func (p *Path) Clean() *Path {
	return New(filepath.Clean(p.path))
}

// IsInsideDir returns true if the current path is inside the provided
// dir
func (p *Path) IsInsideDir(dir *Path) (bool, error) {
	rel, err := filepath.Rel(dir.path, p.path)
	if err != nil {
		// If the dir cannot be made relative to this path it means
		// that it belong to a different filesystems, so it cannot be
		// inside this path.
		return false, nil
	}
	return !strings.Contains(rel, ".."+string(os.PathSeparator)) &&
		rel != ".." &&
		rel != ".", nil
}

// Parent returns all but the last element of path, typically the path's
// directory or the parent directory if the path is already a directory
func (p *Path) Parent() *Path {
	return New(filepath.Dir(p.path))
}

// Mkdir create a directory denoted by the current path
func (p *Path) Mkdir() error {
	return os.Mkdir(p.path, 0755)
}

// MkdirAll creates a directory named path, along with any necessary
// parents, and returns nil, or else returns an error
func (p *Path) MkdirAll() error {
	return os.MkdirAll(p.path, os.FileMode(0755))
}

// Remove removes the named file or directory
func (p *Path) Remove() error {
	return os.Remove(p.path)
}

// RemoveAll removes path and any children it contains. It removes
// everything it can but returns the first error it encounters. If
// the path does not exist, RemoveAll returns nil (no error).
func (p *Path) RemoveAll() error {
	return os.RemoveAll(p.path)
}

// Rename renames (moves) the path to newpath. If newpath already exists
// and is not a directory, Rename replaces it. OS-specific restrictions
// may apply when oldpath and newpath are in different directories. If
// there is an error, it will be of type *os.LinkError.
func (p *Path) Rename(newpath *Path) error {
	return os.Rename(p.path, newpath.path)
}

// MkTempDir creates a new temporary directory inside the path
// pointed by the Path object with a name beginning with prefix
// and returns the path of the new directory.
func (p *Path) MkTempDir(prefix string) (*Path, error) {
	return MkTempDir(p.path, prefix)
}

// FollowSymLink transforms the current path to the path pointed by the
// symlink if path is a symlink, otherwise it does nothing
func (p *Path) FollowSymLink() error {
	resolvedPath, err := filepath.EvalSymlinks(p.path)
	if err != nil {
		return err
	}
	p.path = resolvedPath
	return nil
}

// Exist return true if the file denoted by this path exists, false
// in any other case (also in case of error).
func (p *Path) Exist() bool {
	exist, err := p.ExistCheck()
	return exist && err == nil
}

// NotExist return true if the file denoted by this path DO NOT exists, false
// in any other case (also in case of error).
func (p *Path) NotExist() bool {
	exist, err := p.ExistCheck()
	return !exist && err == nil
}

// ExistCheck return true if the path exists or false if the path doesn't exists.
// In case the check fails false is returned together with the corresponding error.
func (p *Path) ExistCheck() (bool, error) {
	_, err := p.Stat()
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	if err.(*os.PathError).Err == syscall.ENOTDIR {
		return false, nil
	}
	return false, err
}

// IsDir returns true if the path exists and is a directory. In all the other
// cases (and also in case of any error) false is returned.
func (p *Path) IsDir() bool {
	isdir, err := p.IsDirCheck()
	return isdir && err == nil
}

// IsNotDir returns true if the path exists and is NOT a directory. In all the other
// cases (and also in case of any error) false is returned.
func (p *Path) IsNotDir() bool {
	isdir, err := p.IsDirCheck()
	return !isdir && err == nil
}

// IsDirCheck return true if the path exists and is a directory or false
// if the path exists and is not a directory. In all the other case false and
// the corresponding error is returned.
func (p *Path) IsDirCheck() (bool, error) {
	info, err := p.Stat()
	if err == nil {
		return info.IsDir(), nil
	}
	return false, err
}

// CopyTo copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func (p *Path) CopyTo(dst *Path) error {
	if p.EqualsTo(dst) {
		return fmt.Errorf("%s and %s are the same file", p.path, dst.path)
	}

	in, err := os.Open(p.path)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst.path)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	if err := out.Sync(); err != nil {
		return err
	}

	si, err := p.Stat()
	if err != nil {
		return err
	}

	err = os.Chmod(dst.path, si.Mode())
	if err != nil {
		return err
	}

	return nil
}

// CopyTo recursivelly copy all files from a source path to a destination path.
func CopyTo(src *Path, dst *Path) error {
	files, err := src.ReadDirRecursiveFiltered(nil,
		FilterOutDirectories(),
		FilterOutNames("README.md"),
	)
	if err != nil {
		return err
	}
	for _, file := range files {
		destination, err := file.RelFrom(src)
		if err != nil {
			return err
		}
		destination = dst.JoinPath(destination)
		if err := destination.Parent().MkdirAll(); err != nil {
			return err
		}
		if err := file.CopyTo(destination); err != nil {
			return err
		}
	}
	return nil
}

// CopyFS copies the file system fsys into the directory dir,
// creating dir if necessary. It is the exivalent of os.CopyFS with Path.
func (p *Path) CopyFS(dst *Path) error {
	err := os.CopyFS(dst.String(), os.DirFS(p.String()))
	if err != nil {
		return fmt.Errorf("copying %s to %s: %s", p, dst, err)
	}
	return nil
}

// CopyDirTo recursively copies the directory denoted by the current path to
// the destination path. The source directory must exist and the destination
// directory must NOT exist (no implicit destination name allowed).
// Symlinks are not copied, they will be supported in future versions.
func (p *Path) CopyDirTo(dst *Path) error {
	src := p.Clean()
	dst = dst.Clean()

	srcFiles, err := src.ReadDir()
	if err != nil {
		return fmt.Errorf("error reading source dir %s: %s", src, err)
	}

	if exist, err := dst.ExistCheck(); exist {
		return fmt.Errorf("destination %s already exists", dst)
	} else if err != nil {
		return fmt.Errorf("checking if %s exists: %s", dst, err)
	}

	if err := dst.MkdirAll(); err != nil {
		return fmt.Errorf("creating destination dir %s: %s", dst, err)
	}

	srcInfo, err := src.Stat()
	if err != nil {
		return fmt.Errorf("getting stat info for %s: %s", src, err)
	}
	if err := os.Chmod(dst.path, srcInfo.Mode()); err != nil {
		return fmt.Errorf("setting permission for dir %s: %s", dst, err)
	}

	for _, srcPath := range srcFiles {
		srcPathInfo, err := srcPath.Stat()
		if err != nil {
			return fmt.Errorf("getting stat info for %s: %s", srcPath, err)
		}
		dstPath := dst.Join(srcPath.Base())

		if srcPathInfo.IsDir() {
			if err := srcPath.CopyDirTo(dstPath); err != nil {
				return fmt.Errorf("copying %s to %s: %s", srcPath, dstPath, err)
			}
			continue
		}

		// Skip symlinks.
		if srcPathInfo.Mode()&os.ModeSymlink != 0 {
			// TODO
			continue
		}

		if err := srcPath.CopyTo(dstPath); err != nil {
			return fmt.Errorf("copying %s to %s: %s", srcPath, dstPath, err)
		}
	}
	return nil
}

// Chmod changes the mode of the named file to mode. If the file is a
// symbolic link, it changes the mode of the link's target. If there
// is an error, it will be of type *os.PathError.
func (p *Path) Chmod(mode fs.FileMode) error {
	return os.Chmod(p.path, mode)
}

// Chtimes changes the access and modification times of the named file,
// similar to the Unix utime() or utimes() functions.
func (p *Path) Chtimes(atime, mtime time.Time) error {
	return os.Chtimes(p.path, atime, mtime)
}

// ReadFile reads the file named by filename and returns the contents
func (p *Path) ReadFile() ([]byte, error) {
	return os.ReadFile(p.path)
}

// WriteFile writes data to a file named by filename. If the file
// does not exist, WriteFile creates it otherwise WriteFile truncates
// it before writing.
func (p *Path) WriteFile(data []byte) error {
	return os.WriteFile(p.path, data, os.FileMode(0644))
}

// WriteToTempFile writes data to a newly generated temporary file.
// dir and prefix have the same meaning for MkTempFile.
// In case of success the Path to the temp file is returned.
func WriteToTempFile(data []byte, dir *Path, prefix string) (res *Path, err error) {
	f, err := MkTempFile(dir, prefix)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if n, err := f.Write(data); err != nil {
		return nil, err
	} else if n < len(data) {
		return nil, fmt.Errorf("could not write all data (written %d bytes out of %d)", n, len(data))
	}
	return New(f.Name()), nil
}

// ReadFileAsString read a file and return its content as a string.
func (p *Path) ReadFileAsString() (string, error) {
	content, err := p.ReadFile()
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// MustReadFileAsString read a file and return its content as a string. Panic if an error occurs.
func (p *Path) MustReadFileAsString() string {
	content, err := p.ReadFile()
	if err != nil {
		panic(err)
	}
	return string(content)
}

// ReadFileAsLines reads the file named by filename and returns it as an
// array of lines. This function takes care of the newline encoding
// differences between different OS
func (p *Path) ReadFileAsLines() ([]string, error) {
	data, err := p.ReadFile()
	if err != nil {
		return nil, err
	}
	txt := string(data)
	txt = strings.ReplaceAll(txt, "\r\n", "\n")
	return strings.Split(txt, "\n"), nil
}

// MustReadFileAsLines read a file and return its content as a slice of string. Panic if an error occurs.
func (p *Path) MustReadFileAsLines() []string {
	lines, err := p.ReadFileAsLines()
	if err != nil {
		panic(err)
	}
	return lines
}

// MustReadFilteredFileAsLines read a file and return its content as a slice of string.
// It filter out comments and empty lines. Panic if an error occurs.
func (p *Path) MustReadFilteredFileAsLines() []string {
	data, err := p.ReadFile()
	if err != nil {
		panic(err)
	}
	txt := string(data)
	txt = strings.ReplaceAll(txt, "\r\n", "\n")
	txt = util.Filter(txt)
	res := strings.Split(txt, "\n")
	if slices.Contains(res, "") {
		idx := slices.Index(res, "")
		res = slices.Delete(res, idx, idx+1)
	}
	return res
}

// Truncate create an empty file named by path or if the file already
// exist it truncates it (delete all contents)
func (p *Path) Truncate() error {
	return p.WriteFile([]byte{})
}

// Open opens a file for reading. It calls os.Open on the
// underlying path.
func (p *Path) Open() (*os.File, error) {
	return os.Open(p.path)
}

// Create creates or truncates a file. It calls os.Create
// on the underlying path.
func (p *Path) Create() (*os.File, error) {
	return os.Create(p.path)
}

// Append opens a file for append or creates it if the file doesn't exist.
func (p *Path) Append() (*os.File, error) {
	return os.OpenFile(p.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
}

// EqualsTo return true if both paths are equal
func (p *Path) EqualsTo(other *Path) bool {
	return p.path == other.path
}

// EquivalentTo return true if both paths are equivalent (they points to the
// same file even if they are lexicographically different) based on the current
// working directory.
func (p *Path) EquivalentTo(other *Path) bool {
	if p.Clean().path == other.Clean().path {
		return true
	}

	if infoP, err := p.Stat(); err != nil {
		// go ahead with the next test...
	} else if infoOther, err := other.Stat(); err != nil {
		// go ahead with the next test...
	} else if os.SameFile(infoP, infoOther) {
		return true
	}

	if absP, err := p.Abs(); err != nil {
		return false
	} else if absOther, err := other.Abs(); err != nil {
		return false
	} else {
		return absP.path == absOther.path
	}
}

// Parents returns all the parents directories of the current path. If the path is absolute
// it starts from the current path to the root, if the path is relative is starts from the
// current path to the current directory.
// The path should be clean for this method to work properly (no .. or . or other shortcuts).
// This function does not performs any check on the returned paths.
func (p *Path) Parents() []*Path {
	res := []*Path{}
	dir := p
	for {
		res = append(res, dir)
		parent := dir.Parent()
		if parent.EquivalentTo(dir) {
			break
		}
		dir = parent
	}
	return res
}

func (p *Path) String() string {
	return p.path
}

// Canonical return a "canonical" Path for the given filename.
// The meaning of "canonical" is OS-dependent but the goal of this method
// is to always return the same path for a given file (factoring out all the
// possible ambiguities including, for example, relative paths traversal,
// symlinks, drive volume letter case, etc).
func (p *Path) Canonical() *Path {
	canonical := p.Clone()
	// https://github.com/golang/go/issues/17084#issuecomment-246645354
	if err := canonical.FollowSymLink(); err != nil {
		return nil
	}
	if absPath, err := canonical.Abs(); err == nil {
		canonical = absPath
	}
	return canonical
}
