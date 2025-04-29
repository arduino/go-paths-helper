//
// This file is part of PathsHelper library.
//
// Copyright 2023 Arduino AG (http://www.arduino.cc/)
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
