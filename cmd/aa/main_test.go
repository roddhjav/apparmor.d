// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import "testing"

func Test_getIndentationLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "empty",
			input: "",
			want:  0,
		},
		{
			name:  "one-level",
			input: "  /etc/foo r,\n  /etc/bar r,\n\n",
			want:  1,
		},
		{
			name:  "padded-audit-owner",
			input: "              /etc/foo r,\n  audit owner /etc/foo r,\n\n",
			want:  1,
		},
		{
			name:  "padded-audit-link",
			input: "        link /a -> /b,\n  audit link /a -> /b,\n\n",
			want:  1,
		},
		{
			name:  "deep-owner",
			input: "        owner /etc/foo r,\n\n",
			want:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIndentationLevel(tt.input); got != tt.want {
				t.Errorf("getIndentationLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
