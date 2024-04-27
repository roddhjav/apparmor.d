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
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReadDirRecursive(t *testing.T) {
	testPath := New("testdata", "fileset")

	list, err := testPath.ReadDirRecursive()
	require.NoError(t, err)
	require.Len(t, list, 16)

	pathEqualsTo(t, "testdata/fileset/anotherFile", list[0])
	pathEqualsTo(t, "testdata/fileset/file", list[1])
	pathEqualsTo(t, "testdata/fileset/folder", list[2])
	pathEqualsTo(t, "testdata/fileset/folder/.hidden", list[3])
	pathEqualsTo(t, "testdata/fileset/folder/file2", list[4])
	pathEqualsTo(t, "testdata/fileset/folder/file3", list[5])
	pathEqualsTo(t, "testdata/fileset/folder/subfolder", list[6])
	pathEqualsTo(t, "testdata/fileset/folder/subfolder/file4", list[7])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder", list[8])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/.hidden", list[9])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file2", list[10])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file3", list[11])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder", list[12])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder/file4", list[13])
	pathEqualsTo(t, "testdata/fileset/test.txt", list[14])
	pathEqualsTo(t, "testdata/fileset/test.txt.gz", list[15])
}

func TestReadDirRecursiveSymLinkLoop(t *testing.T) {
	// Test symlink loop
	tmp, err := MkTempDir("", "")
	require.NoError(t, err)
	defer tmp.RemoveAll()

	folder := tmp.Join("folder")
	err = os.Symlink(tmp.String(), folder.String())
	require.NoError(t, err)

	l, err := tmp.ReadDirRecursive()
	require.Error(t, err)
	fmt.Println(err)
	require.Nil(t, l)

	l, err = tmp.ReadDirRecursiveFiltered(nil)
	require.Error(t, err)
	fmt.Println(err)
	require.Nil(t, l)
}

func TestReadDirFiltered(t *testing.T) {
	folderPath := New("testdata/fileset/folder")
	list, err := folderPath.ReadDir()
	require.NoError(t, err)
	require.Len(t, list, 4)
	pathEqualsTo(t, "testdata/fileset/folder/.hidden", list[0])
	pathEqualsTo(t, "testdata/fileset/folder/file2", list[1])
	pathEqualsTo(t, "testdata/fileset/folder/file3", list[2])
	pathEqualsTo(t, "testdata/fileset/folder/subfolder", list[3])

	list, err = folderPath.ReadDir(FilterDirectories())
	require.NoError(t, err)
	require.Len(t, list, 1)
	pathEqualsTo(t, "testdata/fileset/folder/subfolder", list[0])

	list, err = folderPath.ReadDir(FilterOutPrefixes("file"))
	require.NoError(t, err)
	require.Len(t, list, 2)
	pathEqualsTo(t, "testdata/fileset/folder/.hidden", list[0])
	pathEqualsTo(t, "testdata/fileset/folder/subfolder", list[1])
}

func TestReadDirRecursiveFiltered(t *testing.T) {
	testdata := New("testdata", "fileset")
	l, err := testdata.ReadDirRecursiveFiltered(nil)
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 16)
	pathEqualsTo(t, "testdata/fileset/anotherFile", l[0])
	pathEqualsTo(t, "testdata/fileset/file", l[1])
	pathEqualsTo(t, "testdata/fileset/folder", l[2])
	pathEqualsTo(t, "testdata/fileset/folder/.hidden", l[3])
	pathEqualsTo(t, "testdata/fileset/folder/file2", l[4])
	pathEqualsTo(t, "testdata/fileset/folder/file3", l[5])
	pathEqualsTo(t, "testdata/fileset/folder/subfolder", l[6])
	pathEqualsTo(t, "testdata/fileset/folder/subfolder/file4", l[7])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder", l[8])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/.hidden", l[9])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file2", l[10])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file3", l[11])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder", l[12])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder/file4", l[13])
	pathEqualsTo(t, "testdata/fileset/test.txt", l[14])
	pathEqualsTo(t, "testdata/fileset/test.txt.gz", l[15])

	l, err = testdata.ReadDirRecursiveFiltered(FilterOutDirectories())
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 6)
	pathEqualsTo(t, "testdata/fileset/anotherFile", l[0])
	pathEqualsTo(t, "testdata/fileset/file", l[1])
	pathEqualsTo(t, "testdata/fileset/folder", l[2])          // <- this is listed but not traversed
	pathEqualsTo(t, "testdata/fileset/symlinktofolder", l[3]) // <- this is listed but not traversed
	pathEqualsTo(t, "testdata/fileset/test.txt", l[4])
	pathEqualsTo(t, "testdata/fileset/test.txt.gz", l[5])

	l, err = testdata.ReadDirRecursiveFiltered(nil, FilterOutDirectories())
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 12)
	pathEqualsTo(t, "testdata/fileset/anotherFile", l[0])
	pathEqualsTo(t, "testdata/fileset/file", l[1])
	pathEqualsTo(t, "testdata/fileset/folder/.hidden", l[2])
	pathEqualsTo(t, "testdata/fileset/folder/file2", l[3])
	pathEqualsTo(t, "testdata/fileset/folder/file3", l[4])
	pathEqualsTo(t, "testdata/fileset/folder/subfolder/file4", l[5])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/.hidden", l[6])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file2", l[7])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file3", l[8])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder/file4", l[9])
	pathEqualsTo(t, "testdata/fileset/test.txt", l[10])
	pathEqualsTo(t, "testdata/fileset/test.txt.gz", l[11])

	l, err = testdata.ReadDirRecursiveFiltered(FilterOutDirectories(), FilterOutDirectories())
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 4)
	pathEqualsTo(t, "testdata/fileset/anotherFile", l[0])
	pathEqualsTo(t, "testdata/fileset/file", l[1])
	pathEqualsTo(t, "testdata/fileset/test.txt", l[2])
	pathEqualsTo(t, "testdata/fileset/test.txt.gz", l[3])

	l, err = testdata.ReadDirRecursiveFiltered(FilterOutPrefixes("sub"), FilterOutSuffixes("3"))
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 12)
	pathEqualsTo(t, "testdata/fileset/anotherFile", l[0])
	pathEqualsTo(t, "testdata/fileset/file", l[1])
	pathEqualsTo(t, "testdata/fileset/folder", l[2])
	pathEqualsTo(t, "testdata/fileset/folder/.hidden", l[3])
	pathEqualsTo(t, "testdata/fileset/folder/file2", l[4])
	pathEqualsTo(t, "testdata/fileset/folder/subfolder", l[5]) // <- subfolder skipped by Prefix("sub")
	pathEqualsTo(t, "testdata/fileset/symlinktofolder", l[6])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/.hidden", l[7])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file2", l[8])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder", l[9]) // <- subfolder skipped by Prefix("sub")
	pathEqualsTo(t, "testdata/fileset/test.txt", l[10])
	pathEqualsTo(t, "testdata/fileset/test.txt.gz", l[11])

	l, err = testdata.ReadDirRecursiveFiltered(FilterOutPrefixes("sub"), AndFilter(FilterOutSuffixes("3"), FilterOutPrefixes("fil")))
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 9)
	pathEqualsTo(t, "testdata/fileset/anotherFile", l[0])
	pathEqualsTo(t, "testdata/fileset/folder", l[1])
	pathEqualsTo(t, "testdata/fileset/folder/.hidden", l[2])
	pathEqualsTo(t, "testdata/fileset/folder/subfolder", l[3])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder", l[4])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/.hidden", l[5])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder", l[6])
	pathEqualsTo(t, "testdata/fileset/test.txt", l[7])
	pathEqualsTo(t, "testdata/fileset/test.txt.gz", l[8])

	l, err = testdata.ReadDirRecursiveFiltered(FilterOutPrefixes("sub"), AndFilter(FilterOutSuffixes("3"), FilterOutPrefixes("fil"), FilterOutSuffixes(".gz")))
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 8)
	pathEqualsTo(t, "testdata/fileset/anotherFile", l[0])
	pathEqualsTo(t, "testdata/fileset/folder", l[1])
	pathEqualsTo(t, "testdata/fileset/folder/.hidden", l[2])
	pathEqualsTo(t, "testdata/fileset/folder/subfolder", l[3])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder", l[4])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/.hidden", l[5])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder", l[6])
	pathEqualsTo(t, "testdata/fileset/test.txt", l[7])

	l, err = testdata.ReadDirRecursiveFiltered(OrFilter(FilterPrefixes("sub"), FilterSuffixes("tofolder")))
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 11)
	pathEqualsTo(t, "testdata/fileset/anotherFile", l[0])
	pathEqualsTo(t, "testdata/fileset/file", l[1])
	pathEqualsTo(t, "testdata/fileset/folder", l[2])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder", l[3])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/.hidden", l[4])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file2", l[5])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file3", l[6])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder", l[7])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder/file4", l[8])
	pathEqualsTo(t, "testdata/fileset/test.txt", l[9])
	pathEqualsTo(t, "testdata/fileset/test.txt.gz", l[10])

	l, err = testdata.ReadDirRecursiveFiltered(nil, FilterNames("folder"))
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 1)
	pathEqualsTo(t, "testdata/fileset/folder", l[0])

	l, err = testdata.ReadDirRecursiveFiltered(FilterNames("symlinktofolder"), FilterOutNames(".hidden"))
	require.NoError(t, err)
	require.Len(t, l, 9)
	l.Sort()
	pathEqualsTo(t, "testdata/fileset/anotherFile", l[0])
	pathEqualsTo(t, "testdata/fileset/file", l[1])
	pathEqualsTo(t, "testdata/fileset/folder", l[2])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder", l[3])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file2", l[4])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/file3", l[5])
	pathEqualsTo(t, "testdata/fileset/symlinktofolder/subfolder", l[6])
	pathEqualsTo(t, "testdata/fileset/test.txt", l[7])
	pathEqualsTo(t, "testdata/fileset/test.txt.gz", l[8])
}

func TestReadDirRecursiveLoopDetection(t *testing.T) {
	loopsPath := New("testdata", "loops")
	unbuondedReaddir := func(testdir string) (PathList, error) {
		var files PathList
		var err error
		done := make(chan bool)
		go func() {
			files, err = loopsPath.Join(testdir).ReadDirRecursive()
			done <- true
		}()
		require.Eventually(
			t,
			func() bool {
				select {
				case <-done:
					return true
				default:
					return false
				}
			},
			5*time.Second,
			10*time.Millisecond,
			"Infinite symlink loop while loading sketch",
		)
		return files, err
	}

	for _, dir := range []string{"loop_1", "loop_2", "loop_3", "loop_4"} {
		l, err := unbuondedReaddir(dir)
		require.EqualError(t, err, "directories symlink loop detected", "loop not detected in %s", dir)
		require.Nil(t, l)
	}

	{
		l, err := unbuondedReaddir("regular_1")
		require.NoError(t, err)
		require.Len(t, l, 4)
		l.Sort()
		pathEqualsTo(t, "testdata/loops/regular_1/dir1", l[0])
		pathEqualsTo(t, "testdata/loops/regular_1/dir1/file1", l[1])
		pathEqualsTo(t, "testdata/loops/regular_1/dir2", l[2])
		pathEqualsTo(t, "testdata/loops/regular_1/dir2/file1", l[3])
	}

	{
		l, err := unbuondedReaddir("regular_2")
		require.NoError(t, err)
		require.Len(t, l, 6)
		l.Sort()
		pathEqualsTo(t, "testdata/loops/regular_2/dir1", l[0])
		pathEqualsTo(t, "testdata/loops/regular_2/dir1/file1", l[1])
		pathEqualsTo(t, "testdata/loops/regular_2/dir2", l[2])
		pathEqualsTo(t, "testdata/loops/regular_2/dir2/dir1", l[3])
		pathEqualsTo(t, "testdata/loops/regular_2/dir2/dir1/file1", l[4])
		pathEqualsTo(t, "testdata/loops/regular_2/dir2/file2", l[5])
	}

	{
		l, err := unbuondedReaddir("regular_3")
		require.NoError(t, err)
		require.Len(t, l, 7)
		l.Sort()
		pathEqualsTo(t, "testdata/loops/regular_3/dir1", l[0])
		pathEqualsTo(t, "testdata/loops/regular_3/dir1/file1", l[1])
		pathEqualsTo(t, "testdata/loops/regular_3/dir2", l[2])
		pathEqualsTo(t, "testdata/loops/regular_3/dir2/dir1", l[3])
		pathEqualsTo(t, "testdata/loops/regular_3/dir2/dir1/file1", l[4])
		pathEqualsTo(t, "testdata/loops/regular_3/dir2/file2", l[5])
		pathEqualsTo(t, "testdata/loops/regular_3/link", l[6]) // broken symlink is reported in files
	}

	if runtime.GOOS != "windows" {
		dir1 := loopsPath.Join("regular_4_with_permission_error", "dir1")

		l, err := unbuondedReaddir("regular_4_with_permission_error")
		require.NoError(t, err)
		require.NotEmpty(t, l)

		dir1Stat, err := dir1.Stat()
		require.NoError(t, err)
		err = dir1.Chmod(fs.FileMode(0)) // Enforce permission error
		require.NoError(t, err)
		t.Cleanup(func() {
			// Restore normal permission after the test
			dir1.Chmod(dir1Stat.Mode())
		})

		l, err = unbuondedReaddir("regular_4_with_permission_error")
		require.Error(t, err)
		require.Nil(t, l)
	}
}
