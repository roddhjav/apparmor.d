// This file is part of PathsHelper library.
// Copyright (C) 2018-2025 Arduino AG (http://www.arduino.cc/)
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package paths

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestProcess_RunWithinContext(t *testing.T) {
	// Build `delay` helper inside testdata/delay
	builder, err := NewProcess(nil, "go", "build")
	if err != nil {
		t.Fatal(err)
	}
	builder.SetDir(testdataRoot + "/delay")
	if err := builder.Run(); err != nil {
		t.Fatal(err)
	}

	// Run delay and test if the process is terminated correctly due to context
	process, err := NewProcess(nil, testdataRoot+"/delay/delay")
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	err = process.RunWithinContext(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	if elapsed := time.Since(start); !(elapsed < 500*time.Millisecond) {
		t.Errorf("%v not less than %v", elapsed, 500*time.Millisecond)
	}
	cancel()
}

func TestProcess_KillProcessGroupOnLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping test on non-linux system")
	}

	p, err := NewProcess(nil, "bash", "-c", "sleep 5 ; echo -n 5")
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, _, err = p.RunAndCaptureOutput(ctx)
	if err == nil || err.Error() != "signal: killed" {
		t.Fatalf("got %v, want signal: killed", err)
	}
	// Assert that the process was killed within the timeout
	if elapsed := time.Since(start); !(elapsed < 2*time.Second) {
		t.Errorf("%v not less than %v", elapsed, 2*time.Second)
	}
}
