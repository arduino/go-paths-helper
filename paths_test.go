/*
 * This file is part of PathsHelper library.
 *
 * Copyright 2018 Arduino AG (http://www.arduino.cc/)
 *
 * PropertiesMap library is free software; you can redistribute it and/or modify
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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPathNew(t *testing.T) {
	test1 := New("path")
	require.Equal(t, "path", test1.String())

	test2 := New("path", "path")
	require.Equal(t, filepath.Join("path", "path"), test2.String())

	test3 := New()
	require.Equal(t, "", test3.String())

	test4 := New("")
	require.Equal(t, "", test4.String())
}

func TestPath(t *testing.T) {
	testPath := New("_testdata")
	require.Equal(t, "_testdata", testPath.String())
	isDir, err := testPath.IsDir()
	require.True(t, isDir)
	require.NoError(t, err)
	exist, err := testPath.Exist()
	require.True(t, exist)
	require.NoError(t, err)

	folderPath := testPath.Join("folder")
	require.Equal(t, "_testdata/folder", folderPath.String())
	isDir, err = folderPath.IsDir()
	require.True(t, isDir)
	require.NoError(t, err)
	exist, err = folderPath.Exist()
	require.True(t, exist)
	require.NoError(t, err)

	filePath := testPath.Join("file")
	require.Equal(t, "_testdata/file", filePath.String())
	isDir, err = filePath.IsDir()
	require.False(t, isDir)
	require.NoError(t, err)
	exist, err = filePath.Exist()
	require.True(t, exist)
	require.NoError(t, err)

	anotherFilePath := filePath.Join("notexistent")
	require.Equal(t, "_testdata/file/notexistent", anotherFilePath.String())
	isDir, err = anotherFilePath.IsDir()
	require.False(t, isDir)
	require.Error(t, err)
	exist, err = anotherFilePath.Exist()
	require.False(t, exist)
	require.Error(t, err)

	list, err := folderPath.ReadDir()
	require.NoError(t, err)
	require.Len(t, list, 4)
	require.Equal(t, "_testdata/folder/.hidden", list[0].String())
	require.Equal(t, "_testdata/folder/file2", list[1].String())
	require.Equal(t, "_testdata/folder/file3", list[2].String())
	require.Equal(t, "_testdata/folder/subfolder", list[3].String())

	list2 := list.Clone()
	list2.FilterDirs()
	require.Len(t, list2, 1)
	require.Equal(t, "_testdata/folder/subfolder", list2[0].String())

	list2 = list.Clone()
	list2.FilterOutHiddenFiles()
	require.Len(t, list2, 3)
	require.Equal(t, "_testdata/folder/file2", list2[0].String())
	require.Equal(t, "_testdata/folder/file3", list2[1].String())
	require.Equal(t, "_testdata/folder/subfolder", list2[2].String())

	list2 = list.Clone()
	list2.FilterOutPrefix("file")
	require.Len(t, list2, 2)
	require.Equal(t, "_testdata/folder/.hidden", list2[0].String())
	require.Equal(t, "_testdata/folder/subfolder", list2[1].String())
}

func TestResetStatCacheWhenFollowingSymlink(t *testing.T) {
	testdata := New("_testdata")
	files, err := testdata.ReadDir()
	require.NoError(t, err)
	for _, file := range files {
		if file.Base() == "symlinktofolder" {
			err = file.FollowSymLink()
			require.NoError(t, err)
			isDir, err := file.IsDir()
			require.NoError(t, err)
			require.True(t, isDir)
			break
		}
	}
}

func TestIsInsideDir(t *testing.T) {
	inside := func(a, b *Path) {
		in, err := a.IsInsideDir(b)
		require.NoError(t, err)
		require.True(t, in)
	}

	notInside := func(a, b *Path) {
		in, err := a.IsInsideDir(b)
		require.NoError(t, err)
		require.False(t, in)
	}

	f1 := New("/a/b/c")
	f2 := New("/a/b/c/d")
	f3 := New("/a/b/c/d/e")

	notInside(f1, f1)
	notInside(f1, f2)
	inside(f2, f1)
	notInside(f1, f3)
	inside(f3, f1)

	r1 := New("a/b/c")
	r2 := New("a/b/c/d")
	r3 := New("a/b/c/d/e")
	r4 := New("f/../a/b/c/d/e")
	r5 := New("a/b/c/d/e/f/..")

	notInside(r1, r1)
	notInside(r1, r2)
	inside(r2, r1)
	notInside(r1, r3)
	inside(r3, r1)
	inside(r4, r1)
	notInside(r1, r4)
	inside(r5, r1)
	notInside(r1, r5)
}
