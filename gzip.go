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
