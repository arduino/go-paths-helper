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
	"io/ioutil"
)

// ReadDirFilter is a filter for Path.ReadDir and Path.ReadDirRecursive methods.
// The filter should return true to accept a file or false to reject it.
type ReadDirFilter func(file *Path) bool

// ReadDir returns a PathList containing the content of the directory
// pointed by the current Path. The resulting list is filtered by the given filters chained.
func (p *Path) ReadDir(filters ...ReadDirFilter) (PathList, error) {
	infos, err := ioutil.ReadDir(p.path)
	if err != nil {
		return nil, err
	}
	paths := PathList{}
fileLoop:
	for _, info := range infos {
		path := p.Join(info.Name())
		for _, filter := range filters {
			if !filter(path) {
				continue fileLoop
			}
		}
		paths.Add(path)
	}
	return paths, nil
}

// ReadDirRecursive returns a PathList containing the content of the directory
// and its subdirectories pointed by the current Path
func (p *Path) ReadDirRecursive() (PathList, error) {
	infos, err := ioutil.ReadDir(p.path)
	if err != nil {
		return nil, err
	}
	paths := PathList{}
	for _, info := range infos {
		path := p.Join(info.Name())
		paths.Add(path)

		if isDir, err := path.IsDirCheck(); err != nil {
			return nil, err
		} else if isDir {
			subPaths, err := path.ReadDirRecursive()
			if err != nil {
				return nil, err
			}
			paths.AddAll(subPaths)
		}

	}
	return paths, nil
}

// FilterDirectories is a ReadDirFilter that accepts only directories
func FilterDirectories() ReadDirFilter {
	return func(path *Path) bool {
		return path.IsDir()
	}
}

// FilterOutDirectories is a ReadDirFilter that rejects all directories
func FilterOutDirectories() ReadDirFilter {
	return func(path *Path) bool {
		return !path.IsDir()
	}
}

// OrFilter creates a ReadDirFilter that accepts all items that are accepted by x or by y
func OrFilter(x, y ReadDirFilter) ReadDirFilter {
	return func(path *Path) bool {
		return x(path) || y(path)
	}
}

// AndFilter creates a ReadDirFilter that accepts all items that are accepted by both x and y
func AndFilter(x, y ReadDirFilter) ReadDirFilter {
	return func(path *Path) bool {
		return x(path) && y(path)
	}
}

// NotFilter creates a ReadDifFilter that accepts all items rejected by x and viceversa
func NotFilter(x ReadDirFilter) ReadDirFilter {
	return func(path *Path) bool {
		return !x(path)
	}
}
