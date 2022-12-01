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
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDirRecursive(t *testing.T) {
	testPath := New("_testdata")

	list, err := testPath.ReadDirRecursive()
	require.NoError(t, err)
	require.Len(t, list, 16)

	pathEqualsTo(t, "_testdata/anotherFile", list[0])
	pathEqualsTo(t, "_testdata/file", list[1])
	pathEqualsTo(t, "_testdata/folder", list[2])
	pathEqualsTo(t, "_testdata/folder/.hidden", list[3])
	pathEqualsTo(t, "_testdata/folder/file2", list[4])
	pathEqualsTo(t, "_testdata/folder/file3", list[5])
	pathEqualsTo(t, "_testdata/folder/subfolder", list[6])
	pathEqualsTo(t, "_testdata/folder/subfolder/file4", list[7])
	pathEqualsTo(t, "_testdata/symlinktofolder", list[8])
	pathEqualsTo(t, "_testdata/symlinktofolder/.hidden", list[9])
	pathEqualsTo(t, "_testdata/symlinktofolder/file2", list[10])
	pathEqualsTo(t, "_testdata/symlinktofolder/file3", list[11])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder", list[12])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder/file4", list[13])
	pathEqualsTo(t, "_testdata/test.txt", list[14])
	pathEqualsTo(t, "_testdata/test.txt.gz", list[15])
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
	folderPath := New("_testdata/folder")
	list, err := folderPath.ReadDir()
	require.NoError(t, err)
	require.Len(t, list, 4)
	pathEqualsTo(t, "_testdata/folder/.hidden", list[0])
	pathEqualsTo(t, "_testdata/folder/file2", list[1])
	pathEqualsTo(t, "_testdata/folder/file3", list[2])
	pathEqualsTo(t, "_testdata/folder/subfolder", list[3])

	list, err = folderPath.ReadDir(FilterDirectories())
	require.NoError(t, err)
	require.Len(t, list, 1)
	pathEqualsTo(t, "_testdata/folder/subfolder", list[0])

	list, err = folderPath.ReadDir(FilterOutPrefixes("file"))
	require.NoError(t, err)
	require.Len(t, list, 2)
	pathEqualsTo(t, "_testdata/folder/.hidden", list[0])
	pathEqualsTo(t, "_testdata/folder/subfolder", list[1])
}

func TestReadDirRecursiveFiltered(t *testing.T) {
	testdata := New("_testdata")
	l, err := testdata.ReadDirRecursiveFiltered(nil)
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 16)
	pathEqualsTo(t, "_testdata/anotherFile", l[0])
	pathEqualsTo(t, "_testdata/file", l[1])
	pathEqualsTo(t, "_testdata/folder", l[2])
	pathEqualsTo(t, "_testdata/folder/.hidden", l[3])
	pathEqualsTo(t, "_testdata/folder/file2", l[4])
	pathEqualsTo(t, "_testdata/folder/file3", l[5])
	pathEqualsTo(t, "_testdata/folder/subfolder", l[6])
	pathEqualsTo(t, "_testdata/folder/subfolder/file4", l[7])
	pathEqualsTo(t, "_testdata/symlinktofolder", l[8])
	pathEqualsTo(t, "_testdata/symlinktofolder/.hidden", l[9])
	pathEqualsTo(t, "_testdata/symlinktofolder/file2", l[10])
	pathEqualsTo(t, "_testdata/symlinktofolder/file3", l[11])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder", l[12])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder/file4", l[13])
	pathEqualsTo(t, "_testdata/test.txt", l[14])
	pathEqualsTo(t, "_testdata/test.txt.gz", l[15])

	l, err = testdata.ReadDirRecursiveFiltered(FilterOutDirectories())
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 6)
	pathEqualsTo(t, "_testdata/anotherFile", l[0])
	pathEqualsTo(t, "_testdata/file", l[1])
	pathEqualsTo(t, "_testdata/folder", l[2])          // <- this is listed but not traversed
	pathEqualsTo(t, "_testdata/symlinktofolder", l[3]) // <- this is listed but not traversed
	pathEqualsTo(t, "_testdata/test.txt", l[4])
	pathEqualsTo(t, "_testdata/test.txt.gz", l[5])

	l, err = testdata.ReadDirRecursiveFiltered(nil, FilterOutDirectories())
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 12)
	pathEqualsTo(t, "_testdata/anotherFile", l[0])
	pathEqualsTo(t, "_testdata/file", l[1])
	pathEqualsTo(t, "_testdata/folder/.hidden", l[2])
	pathEqualsTo(t, "_testdata/folder/file2", l[3])
	pathEqualsTo(t, "_testdata/folder/file3", l[4])
	pathEqualsTo(t, "_testdata/folder/subfolder/file4", l[5])
	pathEqualsTo(t, "_testdata/symlinktofolder/.hidden", l[6])
	pathEqualsTo(t, "_testdata/symlinktofolder/file2", l[7])
	pathEqualsTo(t, "_testdata/symlinktofolder/file3", l[8])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder/file4", l[9])
	pathEqualsTo(t, "_testdata/test.txt", l[10])
	pathEqualsTo(t, "_testdata/test.txt.gz", l[11])

	l, err = testdata.ReadDirRecursiveFiltered(FilterOutDirectories(), FilterOutDirectories())
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 4)
	pathEqualsTo(t, "_testdata/anotherFile", l[0])
	pathEqualsTo(t, "_testdata/file", l[1])
	pathEqualsTo(t, "_testdata/test.txt", l[2])
	pathEqualsTo(t, "_testdata/test.txt.gz", l[3])

	l, err = testdata.ReadDirRecursiveFiltered(FilterOutPrefixes("sub"), FilterOutSuffixes("3"))
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 12)
	pathEqualsTo(t, "_testdata/anotherFile", l[0])
	pathEqualsTo(t, "_testdata/file", l[1])
	pathEqualsTo(t, "_testdata/folder", l[2])
	pathEqualsTo(t, "_testdata/folder/.hidden", l[3])
	pathEqualsTo(t, "_testdata/folder/file2", l[4])
	pathEqualsTo(t, "_testdata/folder/subfolder", l[5]) // <- subfolder skipped by Prefix("sub")
	pathEqualsTo(t, "_testdata/symlinktofolder", l[6])
	pathEqualsTo(t, "_testdata/symlinktofolder/.hidden", l[7])
	pathEqualsTo(t, "_testdata/symlinktofolder/file2", l[8])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder", l[9]) // <- subfolder skipped by Prefix("sub")
	pathEqualsTo(t, "_testdata/test.txt", l[10])
	pathEqualsTo(t, "_testdata/test.txt.gz", l[11])

	l, err = testdata.ReadDirRecursiveFiltered(FilterOutPrefixes("sub"), AndFilter(FilterOutSuffixes("3"), FilterOutPrefixes("fil")))
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 9)
	pathEqualsTo(t, "_testdata/anotherFile", l[0])
	pathEqualsTo(t, "_testdata/folder", l[1])
	pathEqualsTo(t, "_testdata/folder/.hidden", l[2])
	pathEqualsTo(t, "_testdata/folder/subfolder", l[3])
	pathEqualsTo(t, "_testdata/symlinktofolder", l[4])
	pathEqualsTo(t, "_testdata/symlinktofolder/.hidden", l[5])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder", l[6])
	pathEqualsTo(t, "_testdata/test.txt", l[7])
	pathEqualsTo(t, "_testdata/test.txt.gz", l[8])

	l, err = testdata.ReadDirRecursiveFiltered(FilterOutPrefixes("sub"), AndFilter(FilterOutSuffixes("3"), FilterOutPrefixes("fil"), FilterOutSuffixes(".gz")))
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 8)
	pathEqualsTo(t, "_testdata/anotherFile", l[0])
	pathEqualsTo(t, "_testdata/folder", l[1])
	pathEqualsTo(t, "_testdata/folder/.hidden", l[2])
	pathEqualsTo(t, "_testdata/folder/subfolder", l[3])
	pathEqualsTo(t, "_testdata/symlinktofolder", l[4])
	pathEqualsTo(t, "_testdata/symlinktofolder/.hidden", l[5])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder", l[6])
	pathEqualsTo(t, "_testdata/test.txt", l[7])

	l, err = testdata.ReadDirRecursiveFiltered(OrFilter(FilterPrefixes("sub"), FilterSuffixes("tofolder")))
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 11)
	pathEqualsTo(t, "_testdata/anotherFile", l[0])
	pathEqualsTo(t, "_testdata/file", l[1])
	pathEqualsTo(t, "_testdata/folder", l[2])
	pathEqualsTo(t, "_testdata/symlinktofolder", l[3])
	pathEqualsTo(t, "_testdata/symlinktofolder/.hidden", l[4])
	pathEqualsTo(t, "_testdata/symlinktofolder/file2", l[5])
	pathEqualsTo(t, "_testdata/symlinktofolder/file3", l[6])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder", l[7])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder/file4", l[8])
	pathEqualsTo(t, "_testdata/test.txt", l[9])
	pathEqualsTo(t, "_testdata/test.txt.gz", l[10])

	l, err = testdata.ReadDirRecursiveFiltered(nil, FilterNames("folder"))
	require.NoError(t, err)
	l.Sort()
	require.Len(t, l, 1)
	pathEqualsTo(t, "_testdata/folder", l[0])

	l, err = testdata.ReadDirRecursiveFiltered(FilterNames("symlinktofolder"), FilterOutNames(".hidden"))
	require.NoError(t, err)
	require.Len(t, l, 9)
	l.Sort()
	pathEqualsTo(t, "_testdata/anotherFile", l[0])
	pathEqualsTo(t, "_testdata/file", l[1])
	pathEqualsTo(t, "_testdata/folder", l[2])
	pathEqualsTo(t, "_testdata/symlinktofolder", l[3])
	pathEqualsTo(t, "_testdata/symlinktofolder/file2", l[4])
	pathEqualsTo(t, "_testdata/symlinktofolder/file3", l[5])
	pathEqualsTo(t, "_testdata/symlinktofolder/subfolder", l[6])
	pathEqualsTo(t, "_testdata/test.txt", l[7])
	pathEqualsTo(t, "_testdata/test.txt.gz", l[8])
}
