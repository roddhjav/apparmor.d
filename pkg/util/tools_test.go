// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"reflect"
	"testing"
)

func TestDecodeHex(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Hexa",
			str:  "666F6F20626172",
			want: "foo bar",
		},
		{
			name: "Not Hexa",
			str:  "ALLOWED",
			want: "ALLOWED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeHex(tt.str); got != tt.want {
				t.Errorf("DecodeHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicate(t *testing.T) {
	tests := []struct {
		name   string
		inlist []string
		want   []string
	}{
		{
			name:   "Duplicate",
			inlist: []string{"foo", "bar", "foo", "bar", ""},
			want:   []string{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicate(tt.inlist); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
