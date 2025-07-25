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
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func tellCommandNotToSpawnShell(oscmd *exec.Cmd) {
	if oscmd.SysProcAttr == nil {
		oscmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	oscmd.SysProcAttr.HideWindow = true
}

func tellCommandToStartOnNewProcessGroup(_ *exec.Cmd) {
	// no op
}

func kill(oscmd *exec.Cmd) error {
	parentProcessMap, err := createParentProcessSnapshot()
	if err != nil {
		return err
	}
	return killPidTree(uint32(oscmd.Process.Pid), parentProcessMap)
}

// createParentProcessSnapshot returns a map that correlate a process
// with its parent process: childPid -> parentPid
func createParentProcessSnapshot() (map[uint32]uint32, error) {
	// Inspired by: https://stackoverflow.com/a/36089871/1655275

	// Make a snapshot of the current running processes
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("getting running processes snapshot: %w", err)
	}
	defer windows.CloseHandle(snapshot)

	// Iterate the result and extract the parent-child relationship
	processParentMap := map[uint32]uint32{}
	var processEntry windows.ProcessEntry32
	processEntry.Size = uint32(unsafe.Sizeof(processEntry))
	hasData := (windows.Process32First(snapshot, &processEntry) == nil)
	for hasData {
		processParentMap[processEntry.ProcessID] = processEntry.ParentProcessID
		hasData = (windows.Process32Next(snapshot, &processEntry) == nil)
	}
	return processParentMap, nil
}

func killPidTree(pid uint32, parentProcessMap map[uint32]uint32) error {
	for childPid, parentPid := range parentProcessMap {
		if parentPid == pid {
			// Descend process tree
			if err := killPidTree(childPid, parentProcessMap); err != nil {
				return fmt.Errorf("error killing child process: %w", err)
			}
		}
	}
	return killPid(pid)
}

func killPid(pid uint32) error {
	process, err := windows.OpenProcess(windows.PROCESS_ALL_ACCESS, false, pid)
	if err != nil {
		return fmt.Errorf("opening process for kill: %w", err)
	}
	defer windows.CloseHandle(process)
	return windows.TerminateProcess(process, 128)
}
