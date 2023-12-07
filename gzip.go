/*
 * This file is part of PathsHelper library.
 *
 * Copyright 2021 Arduino AG (http://www.arduino.cc/)
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
	"compress/gzip"
	"fmt"
	"io"
)

// GZip compress src with gzip and writes the compressed file on dst
//func GZip(src, dst *Path) error {
//	return errors.New("gzip unimplemented")
//}

// GUnzip decompress src with gzip and writes the uncompressed file on dst
func GUnzip(src, dest *Path) error {
	gzIn, err := src.Open()
	if err != nil {
		return fmt.Errorf("opening %s: %w", src, err)
	}
	defer gzIn.Close()

	in, err := gzip.NewReader(gzIn)
	if err != nil {
		return fmt.Errorf("decoding %s: %w", src, err)
	}
	defer in.Close()

	out, err := dest.Create()
	if err != nil {
		return fmt.Errorf("creating %s: %w", dest, err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("uncompressing %s: %w", dest, err)
	}

	return nil
}
