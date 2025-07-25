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
	"os/exec"
	"syscall"
)

func tellCommandNotToSpawnShell(_ *exec.Cmd) {
	// no op
}

func tellCommandToStartOnNewProcessGroup(oscmd *exec.Cmd) {
	// https://groups.google.com/g/golang-nuts/c/XoQ3RhFBJl8

	// Start the process in a new process group.
	// This is needed to kill the process and its children
	// if we need to kill the process.
	if oscmd.SysProcAttr == nil {
		oscmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	oscmd.SysProcAttr.Setpgid = true
}

func kill(oscmd *exec.Cmd) error {
	// https://groups.google.com/g/golang-nuts/c/XoQ3RhFBJl8

	// Kill the process group
	pgid, err := syscall.Getpgid(oscmd.Process.Pid)
	if err != nil {
		return err
	}
	return syscall.Kill(-pgid, syscall.SIGKILL)
}
