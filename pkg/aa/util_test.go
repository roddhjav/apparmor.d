// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"
)

func Test_toAccess(t *testing.T) {
	tests := []struct {
		name     string
		kind     Kind
		inputs   []string
		wants    [][]string
		wantsErr []bool
	}{
		{
			name:     "empty",
			kind:     FILE,
			inputs:   []string{""},
			wants:    [][]string{nil},
			wantsErr: []bool{false},
		},
		{
			name: "file",
			kind: FILE,
			inputs: []string{
				"rPx", "rPUx", "mr", "rm", "rix", "rcx", "rCUx", "rmix", "rwlk",
				"mrwkl", "", "r", "x", "w", "wr", "px", "Px", "Ux", "mrwlkPix",
			},
			wants: [][]string{
				{"r", "Px"}, {"r", "PUx"}, {"m", "r"}, {"m", "r"}, {"r", "ix"},
				{"r", "cx"}, {"r", "CUx"}, {"m", "r", "ix"}, {"r", "w", "l", "k"},
				{"m", "r", "w", "l", "k"}, nil, {"r"}, {"x"}, {"w"}, {"r", "w"},
				{"px"}, {"Px"}, {"Ux"}, {"m", "r", "w", "l", "k", "Pix"},
			},
			wantsErr: []bool{
				false, false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false,
			},
		},
		{
			name: "file-log",
			kind: FILE + "-log",
			inputs: []string{
				"mr", "rm", "x", "rwlk", "mrwkl", "r", "c", "wc", "d", "wr",
			},
			wants: [][]string{
				{"m", "r"}, {"m", "r"}, {"ix"}, {"r", "w", "l", "k"},
				{"m", "r", "w", "l", "k"}, {"r"}, {"w"}, {"w"}, {"w"}, {"r", "w"},
			},
			wantsErr: []bool{
				false, false, false, false, false, false, false, false, false, false,
			},
		},
		{
			name:     "signal",
			kind:     SIGNAL,
			inputs:   []string{"send receive rw"},
			wants:    [][]string{{"rw", "send", "receive"}},
			wantsErr: []bool{false},
		},
		{
			name:   "ptrace",
			kind:   PTRACE,
			inputs: []string{"readby", "tracedby", "read readby", "r w", "rw", ""},
			wants: [][]string{
				{"readby"}, {"tracedby"}, {"read", "readby"}, {"r", "w"}, {"rw"}, {},
			},
			wantsErr: []bool{false, false, false, false, false, false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, input := range tt.inputs {
				got, err := toAccess(tt.kind, input)
				if (err != nil) != tt.wantsErr[i] {
					t.Errorf("toAccess() error = %v, wantErr %v", err, tt.wantsErr[i])
					return
				}
				if !reflect.DeepEqual(got, tt.wants[i]) {
					t.Errorf("toAccess() = %v, want %v", got, tt.wants[i])
				}
			}
		})
	}
}
