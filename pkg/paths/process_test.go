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
	"context"
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
