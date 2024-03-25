// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"slices"
	"testing"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		names       []string
		wantSuccess bool
	}{
		{
			name:        "test",
			names:       []string{"synchronise", "ignore"},
			wantSuccess: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			Register(tt.names...)
			for _, name := range tt.names {
				if got := slices.Contains(Prepares, Tasks[name]); got != tt.wantSuccess {
					t.Errorf("Register() = %v, want %v", got, tt.wantSuccess)
				}

			}
		})
	}
}
