// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"reflect"
	"testing"
)

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

func TestIntersect(t *testing.T) {
	tests := []struct {
		name  string
		list1 []int
		list2 []int
		want  []int
	}{
		{
			name:  "1",
			list1: []int{0, 1, 2, 3, 4, 5},
			list2: []int{0, 2},
			want:  []int{0, 2},
		},
		{
			name:  "2",
			list1: []int{0, 1, 2, 3, 4, 5},
			list2: []int{0, 6},
			want:  []int{0},
		},
		{
			name:  "3",
			list1: []int{0, 1, 2, 3, 4, 5},
			list2: []int{-1, 6},
			want:  []int{},
		},
		{
			name:  "4",
			list1: []int{0, 6},
			list2: []int{0, 1, 2, 3, 4, 5},
			want:  []int{0},
		},
		{
			name:  "5",
			list1: []int{0, 6, 0},
			list2: []int{0, 1, 2, 3, 4, 5},
			want:  []int{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersect(tt.list1, tt.list2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name  string
		input [][]int
		want  []int
	}{
		{
			name:  "1",
			input: [][]int{{0, 1}, {2, 3, 4, 5}},
			want:  []int{0, 1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Flatten(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvert(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]int
		want  map[int]string
	}{
		{
			name:  "1",
			input: map[string]int{"a": 1, "b": 2},
			want:  map[int]string{1: "a", 2: "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Invert(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Invert() = %v, want %v", got, tt.want)
			}
		})
	}
}
