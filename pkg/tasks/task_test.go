// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package tasks

import (
	"slices"
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func TestBaseTask_Name(t *testing.T) {
	tests := []struct {
		name string
		b    BaseTask
		want string
	}{
		{
			name: "simple",
			b:    BaseTask{Keyword: "test"},
			want: "test",
		},
		{
			name: "with-dashes",
			b:    BaseTask{Keyword: "test-task"},
			want: "test-task",
		},
		{
			name: "empty",
			b:    BaseTask{Keyword: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Name(); got != tt.want {
				t.Errorf("BaseTask.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseTask_Usage(t *testing.T) {
	tests := []struct {
		name string
		b    BaseTask
		want []string
	}{
		{
			name: "single",
			b:    BaseTask{Help: []string{"test"}},
			want: []string{"test"},
		},
		{
			name: "multiple",
			b:    BaseTask{Help: []string{"line1", "line2", "line3"}},
			want: []string{"line1", "line2", "line3"},
		},
		{
			name: "empty",
			b:    BaseTask{Help: []string{}},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Usage(); !slices.Equal(got, tt.want) {
				t.Errorf("BaseTask.Usage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseTask_Message(t *testing.T) {
	tests := []struct {
		name string
		b    BaseTask
		want string
	}{
		{
			name: "simple",
			b:    BaseTask{Msg: "test message"},
			want: "test message",
		},
		{
			name: "empty",
			b:    BaseTask{Msg: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Message(); got != tt.want {
				t.Errorf("BaseTask.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseTask_SetConfig(t *testing.T) {
	tests := []struct {
		name   string
		root   string
		wantAA string
	}{
		{
			name:   "standard",
			root:   "/tmp/build",
			wantAA: "/tmp/build/apparmor.d",
		},
		{
			name:   "root",
			root:   "/",
			wantAA: "/apparmor.d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := paths.New(tt.root)
			config := NewTaskConfig(root)
			task := &BaseTask{}
			task.SetConfig(config)

			if task.Root.String() != tt.root {
				t.Errorf("BaseTask.SetConfig() Root = %v, want %v", task.Root, tt.root)
			}
			if task.RootApparmor.String() != tt.wantAA {
				t.Errorf("BaseTask.SetConfig() RootApparmor = %v, want %v", task.RootApparmor, tt.wantAA)
			}
		})
	}
}

func TestNewTaskConfig(t *testing.T) {
	tests := []struct {
		name     string
		root     string
		wantRoot string
		wantAA   string
	}{
		{
			name:     "standard",
			root:     "/tmp/build",
			wantRoot: "/tmp/build",
			wantAA:   "/tmp/build/apparmor.d",
		},
		{
			name:     "root",
			root:     "/",
			wantRoot: "/",
			wantAA:   "/apparmor.d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := paths.New(tt.root)
			got := NewTaskConfig(root)
			if got.Root.String() != tt.wantRoot {
				t.Errorf("NewTaskConfig().Root = %v, want %v", got.Root, tt.wantRoot)
			}
			if got.RootApparmor.String() != tt.wantAA {
				t.Errorf("NewTaskConfig().RootApparmor = %v, want %v", got.RootApparmor, tt.wantAA)
			}
		})
	}
}

func TestNewBaseRunner(t *testing.T) {
	tests := []struct {
		name string
		root string
	}{
		{
			name: "standard",
			root: "/tmp/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := paths.New(tt.root)
			config := NewTaskConfig(root)
			runner := NewBaseRunner[*BaseTask](config)

			if runner == nil {
				t.Fatal("NewBaseRunner() returned nil")
			}
			if runner.Root.String() != tt.root {
				t.Errorf("NewBaseRunner().Root = %v, want %v", runner.Root, tt.root)
			}
			if len(runner.Tasks) != 0 {
				t.Errorf("NewBaseRunner().Tasks length = %v, want 0", len(runner.Tasks))
			}
		})
	}
}

func TestBaseRunner_Add(t *testing.T) {
	tests := []struct {
		name      string
		tasks     []*BaseTask
		wantCount int
	}{
		{
			name: "single",
			tasks: []*BaseTask{
				{Keyword: "task1", Help: []string{"help1"}, Msg: "msg1"},
			},
			wantCount: 1,
		},
		{
			name: "multiple",
			tasks: []*BaseTask{
				{Keyword: "task1", Help: []string{"help1"}, Msg: "msg1"},
				{Keyword: "task2", Help: []string{"help2"}, Msg: "msg2"},
				{Keyword: "task3", Help: []string{"help3"}, Msg: "msg3"},
			},
			wantCount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := paths.New("/tmp/test")
			config := NewTaskConfig(root)
			runner := NewBaseRunner[*BaseTask](config)

			for _, task := range tt.tasks {
				runner.Add(task)
			}

			if len(runner.Tasks) != tt.wantCount {
				t.Errorf("BaseRunner.Add() tasks length = %v, want %v", len(runner.Tasks), tt.wantCount)
			}

			// Verify tasks received config
			for i, task := range runner.Tasks {
				if task.Root.String() != root.String() {
					t.Errorf("Task[%d] config not set, Root = %v, want %v", i, task.Root, root)
				}
			}
		})
	}
}

func TestBaseRunner_Help(t *testing.T) {
	tests := []struct {
		name        string
		runnerName  string
		tasks       []*BaseTask
		wantStrings []string
	}{
		{
			name:       "single",
			runnerName: "test-runner",
			tasks: []*BaseTask{
				{Keyword: "build", Help: []string{"build help"}, Msg: "Build the project"},
			},
			wantStrings: []string{"test-runner tasks:", "build", "Build the project"},
		},
		{
			name:       "multiple",
			runnerName: "suite",
			tasks: []*BaseTask{
				{Keyword: "build", Help: []string{"build help"}, Msg: "Build the project"},
				{Keyword: "test", Help: []string{"test help"}, Msg: "Run tests"},
			},
			wantStrings: []string{"suite tasks:", "build", "Build the project", "test", "Run tests"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := paths.New("/tmp/test")
			config := NewTaskConfig(root)
			runner := NewBaseRunner[*BaseTask](config)

			for _, task := range tt.tasks {
				runner.Add(task)
			}

			got := runner.Help(tt.runnerName)

			for _, want := range tt.wantStrings {
				if !strings.Contains(got, want) {
					t.Errorf("BaseRunner.Help() missing expected string %q\nGot: %s", want, got)
				}
			}
		})
	}
}
