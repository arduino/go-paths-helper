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
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProcessWithinContext(t *testing.T) {
	// Build `delay` helper inside testdata/delay
	builder, err := NewProcess(nil, "go", "build")
	require.NoError(t, err)
	builder.SetDir("testdata/delay")
	require.NoError(t, builder.Run())

	// Run delay and test if the process is terminated correctly due to context
	process, err := NewProcess(nil, "testdata/delay/delay")
	require.NoError(t, err)
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	err = process.RunWithinContext(ctx)
	require.Error(t, err)
	require.Less(t, time.Since(start), 500*time.Millisecond)
	cancel()
}

func TestKillProcessGroupOnLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping test on non-linux system")
	}

	p, err := NewProcess(nil, "bash", "-c", "sleep 5 ; echo -n 5")
	require.NoError(t, err)
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, _, err = p.RunAndCaptureOutput(ctx)
	require.EqualError(t, err, "signal: killed")
	// Assert that the process was killed within the timeout
	require.Less(t, time.Since(start), 2*time.Second)
}
