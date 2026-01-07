// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package tasks

import (
	"slices"
	"strings"
	"testing"
)

func TestBase_Helpers(t *testing.T) {
	tests := []struct {
		name string
		b    Base
		want string
	}{
		{
			name: "base",
			b:    Base{Keyword: "test", Help: []string{"test"}, Msg: "test"},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Name(); got != tt.want {
				t.Errorf("Base.Name() = %v, want %v", got, tt.want)
			}
			if got := tt.b.Usage(); !slices.Equal(got, []string{tt.want}) {
				t.Errorf("Base.Usage() = %v, want %v", got, tt.want)
			}
			if got := tt.b.Message(); got != tt.want {
				t.Errorf("Base.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelp(t *testing.T) {
	tests := []struct {
		name  string
		tasks map[string]Base
		want  string
	}{
		{
			name: "one",
			tasks: map[string]Base{
				"one": {Keyword: "one", Help: []string{"one"}, Msg: "one"},
				"two": {Keyword: "two", Help: []string{"two"}, Msg: "two"},
			},
			want: `one`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Help(tt.name, tt.tasks); !strings.Contains(got, tt.want) {
				t.Errorf("Help() = %v, want %v", got, tt.want)
			}
		})
	}
}
