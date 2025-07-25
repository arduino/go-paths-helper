// This file is part of PathsHelper library.
//
// Copyright 2018-2025 Arduino AG (http://www.arduino.cc/)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package paths

import (
	"os"
	"runtime"
)

// NullPath return the path to the /dev/null equivalent for the current OS
func NullPath() *Path {
	if runtime.GOOS == "windows" {
		return New("nul")
	}
	return New("/dev/null")
}

// TempDir returns the default path to use for temporary files
func TempDir() *Path {
	return New(os.TempDir()).Canonical()
}

// MkTempDir creates a new temporary directory in the directory
// dir with a name beginning with prefix and returns the path of
// the new directory. If dir is the empty string, TempDir uses the
// default directory for temporary files
func MkTempDir(dir, prefix string) (*Path, error) {
	path, err := os.MkdirTemp(dir, prefix)
	if err != nil {
		return nil, err
	}
	return New(path).Canonical(), nil
}

// MkTempFile creates a new temporary file in the directory dir with a name beginning with prefix,
// opens the file for reading and writing, and returns the resulting *os.File. If dir is nil,
// MkTempFile uses the default directory for temporary files (see paths.TempDir). Multiple programs
// calling TempFile simultaneously will not choose the same file. The caller can use f.Name() to
// find the pathname of the file. It is the caller's responsibility to remove the file when no longer needed.
func MkTempFile(dir *Path, prefix string) (*os.File, error) {
	tmpDir := ""
	if dir != nil {
		tmpDir = dir.String()
	}
	return os.CreateTemp(tmpDir, prefix)
}

// Getwd returns a rooted path name corresponding to the current
// directory.
func Getwd() (*Path, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return New(wd), nil
}
