//
// This file is part of PathsHelper library.
//
// Copyright 2025 Arduino AG (http://www.arduino.cc/)
//
// PathsHelper library is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
//
// As a special exception, you may use this file as part of a free software
// library without restriction.  Specifically, if other files instantiate
// templates or use macros or inline functions from this file, or you compile
// this file and link it with other files to produce an executable, this
// file does not by itself cause the resulting executable to be covered by
// the GNU General Public License.  This exception does not however
// invalidate any other reasons why the executable file might be covered by
// the GNU General Public License.
//

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
