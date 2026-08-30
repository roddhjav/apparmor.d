// aa-flatpak - Confine flatpak applications
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// eventTimeout is how long a test waits for a debounced callback.
const eventTimeout = 5 * time.Second

// wantEvent asserts the next callback argument received on ch.
func wantEvent(t *testing.T, ch <-chan string, want string) {
	t.Helper()
	select {
	case got := <-ch:
		if got != want {
			t.Errorf("Watch() = %v, want %v", got, want)
		}
	case <-time.After(eventTimeout):
		t.Errorf("Watch() = no event, want %v", want)
	}
}

// wantNoEvent asserts no callback fires within a few debounce windows.
func wantNoEvent(t *testing.T, ch <-chan string) {
	t.Helper()
	select {
	case got := <-ch:
		t.Errorf("Watch() = %v, want no event", got)
	case <-time.After(2 * debounceDelay):
	}
}

func TestFileWatcher_Watch(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		remove        bool
		wantProcessed string // "" means no process callback is expected
		wantRemoved   string // "" means no remove callback is expected
	}{
		{
			name:          "created application is processed",
			file:          "com.google.Chrome",
			wantProcessed: "com.google.Chrome",
		},
		{
			name: "editor swap file is ignored",
			file: "com.google.Chrome.swp",
		},
		{
			name: "editor backup file is ignored",
			file: "com.google.Chrome~",
		},
		{
			name:        "removed application is removed",
			file:        "com.google.Chrome",
			remove:      true,
			wantRemoved: "com.google.Chrome",
		},
	}

	setupTmp(t)
	dir := t.TempDir()
	processed := make(chan string, 8)
	removed := make(chan string, 8)
	fw, err := NewFileWatcher(
		func(name string) error {
			processed <- name
			return nil
		},
		func(name string) error {
			removed <- name
			return nil
		},
	)
	if err != nil {
		t.Fatalf("NewFileWatcher() error = %v, wantErr %v", err, false)
	}
	defer fw.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- fw.Watch(ctx, dir) }()

	// Let the watcher register its inotify watch before touching the directory.
	time.Sleep(200 * time.Millisecond)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.file)
			if tt.remove {
				if err := os.Remove(path); err != nil {
					t.Fatalf("remove %s: %v", tt.file, err)
				}
			} else {
				if err := os.WriteFile(path, []byte("[Application]\n"), 0o644); err != nil {
					t.Fatalf("write %s: %v", tt.file, err)
				}
			}

			switch {
			case tt.wantProcessed != "":
				wantEvent(t, processed, tt.wantProcessed)
			case tt.wantRemoved != "":
				wantEvent(t, removed, tt.wantRemoved)
			default:
				wantNoEvent(t, processed)
			}
		})
	}

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Watch() error = %v, wantErr %v", err, false)
		}
	case <-time.After(eventTimeout):
		t.Errorf("Watch() did not return on context cancellation")
	}
}

// TestFileWatcher_WatchMissingDir covers a directory that does not exist yet:
// the watcher must keep running instead of failing.
func TestFileWatcher_WatchMissingDir(t *testing.T) {
	setupTmp(t)
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	fw, err := NewFileWatcher(
		func(string) error { return nil },
		func(string) error { return nil },
	)
	if err != nil {
		t.Fatalf("NewFileWatcher() error = %v, wantErr %v", err, false)
	}
	defer fw.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	if err := fw.Watch(ctx, missing); err != nil {
		t.Errorf("Watch() error = %v, wantErr %v", err, false)
	}
}
