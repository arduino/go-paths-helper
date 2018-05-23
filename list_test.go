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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListConstructors(t *testing.T) {
	list0 := NewPathList()
	require.Len(t, list0, 0)

	list1 := NewPathList("test")
	require.Len(t, list1, 1)
	require.Equal(t, "[test]", fmt.Sprintf("%s", list1))

	list3 := NewPathList("a", "b", "c")
	require.Len(t, list3, 3)
	require.Equal(t, "[a b c]", fmt.Sprintf("%s", list3))

	require.False(t, list3.Contains(New("d")))
	require.True(t, list3.Contains(New("a")))
	require.False(t, list3.Contains(New("d/../a")))

	require.False(t, list3.ContainsEquivalentTo(New("d")))
	require.True(t, list3.ContainsEquivalentTo(New("a")))
	require.True(t, list3.ContainsEquivalentTo(New("d/../a")))
}
