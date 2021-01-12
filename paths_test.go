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
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPathNew(t *testing.T) {
	test1 := New("path")
	require.Equal(t, "path", test1.String())

	test2 := New("path", "path")
	require.Equal(t, filepath.Join("path", "path"), test2.String())

	test3 := New()
	require.Nil(t, test3)

	test4 := New("")
	require.Nil(t, test4)
}

func TestPath(t *testing.T) {
	testPath := New("_testdata")
	require.Equal(t, "_testdata", testPath.String())
	isDir, err := testPath.IsDirCheck()
	require.True(t, isDir)
	require.NoError(t, err)
	require.True(t, testPath.IsDir())
	require.False(t, testPath.IsNotDir())
	exist, err := testPath.ExistCheck()
	require.True(t, exist)
	require.NoError(t, err)
	require.True(t, testPath.Exist())
	require.False(t, testPath.NotExist())

	folderPath := testPath.Join("folder")
	require.Equal(t, "_testdata/folder", folderPath.String())
	isDir, err = folderPath.IsDirCheck()
	require.True(t, isDir)
	require.NoError(t, err)
	require.True(t, folderPath.IsDir())
	require.False(t, folderPath.IsNotDir())

	exist, err = folderPath.ExistCheck()
	require.True(t, exist)
	require.NoError(t, err)
	require.True(t, folderPath.Exist())
	require.False(t, folderPath.NotExist())

	filePath := testPath.Join("file")
	require.Equal(t, "_testdata/file", filePath.String())
	isDir, err = filePath.IsDirCheck()
	require.False(t, isDir)
	require.NoError(t, err)
	require.False(t, filePath.IsDir())
	require.True(t, filePath.IsNotDir())
	exist, err = filePath.ExistCheck()
	require.True(t, exist)
	require.NoError(t, err)
	require.True(t, filePath.Exist())
	require.False(t, filePath.NotExist())

	anotherFilePath := filePath.Join("notexistent")
	require.Equal(t, "_testdata/file/notexistent", anotherFilePath.String())
	isDir, err = anotherFilePath.IsDirCheck()
	require.False(t, isDir)
	require.Error(t, err)
	require.False(t, anotherFilePath.IsDir())
	require.False(t, anotherFilePath.IsNotDir())
	exist, err = anotherFilePath.ExistCheck()
	require.False(t, exist)
	require.Error(t, err)
	require.False(t, anotherFilePath.Exist())
	require.False(t, anotherFilePath.NotExist())

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
			isDir, err := file.IsDirCheck()
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
		require.True(t, in, "%s is inside %s", a, b)
	}

	notInside := func(a, b *Path) {
		in, err := a.IsInsideDir(b)
		require.NoError(t, err)
		require.False(t, in, "%s is inside %s", a, b)
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

	f4 := New("/home/megabug/aide/arduino-1.8.6/hardware/arduino/avr")
	f5 := New("/home/megabug/a15/packages")
	notInside(f5, f4)
	notInside(f4, f5)
}

func TestReadFileAsLines(t *testing.T) {
	lines, err := New("_testdata/anotherFile").ReadFileAsLines()
	require.NoError(t, err)
	require.Len(t, lines, 4)
	require.Equal(t, "line 1", lines[0])
	require.Equal(t, "line 2", lines[1])
	require.Equal(t, "", lines[2])
	require.Equal(t, "line 3", lines[3])
}

func TestCopyDir(t *testing.T) {
	tmp, err := MkTempDir("", "")
	require.NoError(t, err)
	defer tmp.RemoveAll()

	src := New("_testdata")
	err = src.CopyDirTo(tmp.Join("dest"))
	require.NoError(t, err, "copying dir")

	exist, err := tmp.Join("dest", "folder", "subfolder", "file4").ExistCheck()
	require.True(t, exist)
	require.NoError(t, err)

	isdir, err := tmp.Join("dest", "folder", "subfolder", "file4").IsDirCheck()
	require.False(t, isdir)
	require.NoError(t, err)

	err = src.CopyDirTo(tmp.Join("dest"))
	require.Error(t, err, "copying dir to already existing")

	err = src.Join("file").CopyDirTo(tmp.Join("dest2"))
	require.Error(t, err, "copying file as dir")
}

func TestParents(t *testing.T) {
	parents := New("/a/very/long/path").Parents()
	require.Len(t, parents, 5)
	require.Equal(t, "/a/very/long/path", parents[0].String())
	require.Equal(t, "/a/very/long", parents[1].String())
	require.Equal(t, "/a/very", parents[2].String())
	require.Equal(t, "/a", parents[3].String())
	require.Equal(t, "/", parents[4].String())

	parents2 := New("a/very/relative/path").Parents()
	require.Len(t, parents, 5)
	require.Equal(t, "a/very/relative/path", parents2[0].String())
	require.Equal(t, "a/very/relative", parents2[1].String())
	require.Equal(t, "a/very", parents2[2].String())
	require.Equal(t, "a", parents2[3].String())
	require.Equal(t, ".", parents2[4].String())
}

func TestReadDirRecursive(t *testing.T) {
	testPath := New("_testdata")

	list, err := testPath.ReadDirRecursive()
	require.NoError(t, err)
	require.Len(t, list, 14)

	require.Equal(t, "_testdata/anotherFile", list[0].String())
	require.Equal(t, "_testdata/file", list[1].String())
	require.Equal(t, "_testdata/folder", list[2].String())
	require.Equal(t, "_testdata/folder/.hidden", list[3].String())
	require.Equal(t, "_testdata/folder/file2", list[4].String())
	require.Equal(t, "_testdata/folder/file3", list[5].String())
	require.Equal(t, "_testdata/folder/subfolder", list[6].String())
	require.Equal(t, "_testdata/folder/subfolder/file4", list[7].String())
	require.Equal(t, "_testdata/symlinktofolder", list[8].String())
	require.Equal(t, "_testdata/symlinktofolder/.hidden", list[9].String())
	require.Equal(t, "_testdata/symlinktofolder/file2", list[10].String())
	require.Equal(t, "_testdata/symlinktofolder/file3", list[11].String())
	require.Equal(t, "_testdata/symlinktofolder/subfolder", list[12].String())
	require.Equal(t, "_testdata/symlinktofolder/subfolder/file4", list[13].String())
}

func TestFilterDirs(t *testing.T) {
	testPath := New("_testdata")

	list, err := testPath.ReadDir()
	require.NoError(t, err)
	require.Len(t, list, 4)

	require.Equal(t, "_testdata/anotherFile", list[0].String())
	require.Equal(t, "_testdata/file", list[1].String())
	require.Equal(t, "_testdata/folder", list[2].String())
	require.Equal(t, "_testdata/symlinktofolder", list[3].String())

	list.FilterDirs()
	require.Len(t, list, 2)
	require.Equal(t, "_testdata/folder", list[0].String())
	require.Equal(t, "_testdata/symlinktofolder", list[1].String())
}

func TestFilterOutDirs(t *testing.T) {
	testPath := New("_testdata")

	list, err := testPath.ReadDir()
	require.NoError(t, err)
	require.Len(t, list, 4)

	require.Equal(t, "_testdata/anotherFile", list[0].String())
	require.Equal(t, "_testdata/file", list[1].String())
	require.Equal(t, "_testdata/folder", list[2].String())
	require.Equal(t, "_testdata/symlinktofolder", list[3].String())

	list.FilterOutDirs()
	require.Len(t, list, 2)
	require.Equal(t, "_testdata/anotherFile", list[0].String())
	require.Equal(t, "_testdata/file", list[1].String())
}

func TestEquivalentPaths(t *testing.T) {
	wd, err := Getwd()
	require.NoError(t, err)
	require.True(t, New("file1").EquivalentTo(New("file1", "somethingelse", "..")))
	require.True(t, New("file1", "abc").EquivalentTo(New("file1", "abc", "def", "..")))
	require.True(t, wd.Join("file1").EquivalentTo(New("file1")))
	require.True(t, wd.Join("file1").EquivalentTo(New("file1", "abc", "..")))

	if runtime.GOOS == "windows" {
		q := New("_testdata", "anotherFile")
		r := New("_testdata", "ANOTHE~1")
		require.True(t, q.EquivalentTo(r))
		require.True(t, r.EquivalentTo(q))
	}
}

func TestCanonicalize(t *testing.T) {
	wd, err := Getwd()
	require.NoError(t, err)

	p := New("_testdata", "anotherFile").Canonical()
	require.Equal(t, wd.Join("_testdata", "anotherFile").String(), p.String())

	p = New("_testdata", "nonexistentFile").Canonical()
	require.Equal(t, wd.Join("_testdata", "nonexistentFile").String(), p.String())

	if runtime.GOOS == "windows" {
		q := New("_testdata", "ANOTHE~1").Canonical()
		require.Equal(t, wd.Join("_testdata", "anotherFile").String(), q.String())

		r := New("c:\\").Canonical()
		require.Equal(t, "C:\\", r.String())
	}
}

func TestRelativeTo(t *testing.T) {
	res, err := New("/my/abs/path/123/456").RelTo(New("/my/abs/path"))
	require.NoError(t, err)
	require.Equal(t, "../..", res.String())

	res, err = New("/my/abs/path").RelTo(New("/my/abs/path/123/456"))
	require.NoError(t, err)
	require.Equal(t, "123/456", res.String())

	res, err = New("my/path").RelTo(New("/other/path"))
	require.Error(t, err)
	require.Nil(t, res)

	res, err = New("/my/abs/path/123/456").RelFrom(New("/my/abs/path"))
	require.Equal(t, "123/456", res.String())
	require.NoError(t, err)

	res, err = New("/my/abs/path").RelFrom(New("/my/abs/path/123/456"))
	require.NoError(t, err)
	require.Equal(t, "../..", res.String())

	res, err = New("my/path").RelFrom(New("/other/path"))
	require.Error(t, err)
	require.Nil(t, res)
}

func TestWriteToTempFile(t *testing.T) {
	tmpDir := New("_testdata", "tmp")
	tmpData := []byte("test")
	tmp, err := WriteToTempFile(tmpData, tmpDir, "prefix")
	defer tmp.Remove()
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(tmp.Base(), "prefix"))
	inside, err := tmp.IsInsideDir(tmpDir)
	require.NoError(t, err)
	require.True(t, inside)
	data, err := tmp.ReadFile()
	require.NoError(t, err)
	require.Equal(t, tmpData, data)
}
