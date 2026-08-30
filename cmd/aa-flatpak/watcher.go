// aa-flatpak - Confine flatpak applications
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	// retryInterval is how often a missing watched directory is checked for
	retryInterval = 10 * time.Second

	// debounceDelay is how long an event is held back to coalesce a burst
	debounceDelay = 500 * time.Millisecond
)

// FileWatcher manages file system notification events with debouncing
type FileWatcher struct {
	watcher     *fsnotify.Watcher
	debounce    map[string]*time.Timer
	lock        map[string]bool
	mu          sync.Mutex
	processFunc func(string) error
	removeFunc  func(string) error
}

// NewFileWatcher creates a new file watcher with debouncing
func NewFileWatcher(processFunc, removeFunc func(string) error) (*FileWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	return &FileWatcher{
		watcher:     w,
		debounce:    make(map[string]*time.Timer),
		lock:        make(map[string]bool),
		processFunc: processFunc,
		removeFunc:  removeFunc,
	}, nil
}

// Close stops the pending debounce timers and the underlying watcher.
func (fw *FileWatcher) Close() error {
	fw.mu.Lock()
	for _, timer := range fw.debounce {
		timer.Stop()
	}
	fw.mu.Unlock()
	return fw.watcher.Close()
}

// Watch starts watching for file changes in the specified directories
func (fw *FileWatcher) Watch(ctx context.Context, dirs ...string) error {
	// A directory may not exist yet (no flatpak application installed) or may be
	// removed later: (re)start watching it as soon as it shows up.
	watching := map[string]bool{}
	resync := func(initial bool) {
		for _, dir := range dirs {
			if err := fw.watcher.Add(dir); err != nil {
				if initial || watching[dir] {
					log.Printf("Waiting for %s to be created", dir)
				}
				watching[dir] = false
				continue
			}
			if watching[dir] {
				continue
			}
			watching[dir] = true
			log.Printf("Watching %s for changes", dir)
			if initial {
				continue // Already generated on startup
			}
			entries, _ := os.ReadDir(dir) // Content created alongside the directory
			for _, entry := range entries {
				if err := fw.processFunc(entry.Name()); err != nil {
					log.Printf("Failed to process %s: %v", entry.Name(), err)
				}
			}
		}
	}
	resync(true)

	retry := time.NewTicker(retryInterval)
	defer retry.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-retry.C:
			resync(false)

		case event, ok := <-fw.watcher.Events:
			if !ok {
				return fmt.Errorf("watcher event channel closed")
			}
			if strings.HasSuffix(event.Name, ".swp") || strings.HasSuffix(event.Name, "~") {
				continue
			}

			if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				name := filepath.Base(event.Name)
				if err := fw.removeFunc(name); err != nil {
					log.Printf("Failed to remove %s: %v", name, err)
				}
				continue
			}
			if !event.Has(fsnotify.Write) && !event.Has(fsnotify.Create) {
				continue
			}

			fw.mu.Lock()
			if fw.lock[event.Name] {
				fw.mu.Unlock()
				continue
			}
			if timer, exists := fw.debounce[event.Name]; exists {
				timer.Stop()
			}
			fw.debounce[event.Name] = time.AfterFunc(debounceDelay,
				func() {
					fw.mu.Lock()
					fw.lock[event.Name] = true
					fw.mu.Unlock()

					name := filepath.Base(event.Name)
					if err := fw.processFunc(name); err != nil {
						log.Printf("Failed to process %s: %v", name, err)
					}

					fw.mu.Lock()
					delete(fw.lock, event.Name)
					delete(fw.debounce, event.Name)
					fw.mu.Unlock()
				},
			)
			fw.mu.Unlock()

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return fmt.Errorf("watcher error channel closed")
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
