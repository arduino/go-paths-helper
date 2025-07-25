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

//go:build windows

package paths

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
)

func shorten(t *testing.T, longPath string) string {
	var buf [4096]uint16
	shortPath := &buf[0]
	n, err := windows.GetShortPathName(windows.StringToUTF16Ptr(longPath), &buf[0], uint32(len(buf)))
	if n >= uint32(len(buf)) {
		buf2 := make([]uint16, n+1)
		shortPath = &buf2[0]
		_, err = windows.GetShortPathName(windows.StringToUTF16Ptr(longPath), &buf2[0], uint32(len(buf2)))
	}
	require.NoError(t, err, "GetShortPathName failed for %v", longPath)
	return windows.UTF16PtrToString(shortPath)
}
