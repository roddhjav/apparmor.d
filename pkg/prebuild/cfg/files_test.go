// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package cfg

import (
	"reflect"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func TestFlagger_Read(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    map[string][]string
	}{
		{
			name: "empty",
			content: `

`,
			want: map[string][]string{},
		},
		{
			name: "main",
			content: `
# Common profile flags definition for all distributions
# File format: one profile by line using the format: '<profile> <flags>'

bwrap attach_disconnected,mediate_deleted,complain
bwrap-app attach_disconnected,complain

akonadi_akonotes_resource complain # Dev
gnome-disks complain

`,
			want: map[string][]string{
				"akonadi_akonotes_resource": {"complain"},
				"bwrap":                     {"attach_disconnected", "mediate_deleted", "complain"},
				"bwrap-app":                 {"attach_disconnected", "complain"},
				"gnome-disks":               {"complain"},
			},
		},
	}
	FlagDir = paths.New("/tmp/")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FlagDir.Join(tt.name + ".flags").WriteFile([]byte(tt.content))
			if err != nil {
				return
			}
			if got := Flags.Read(tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flagger.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIgnore_Read(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name: "empty",
			content: `

`,
			want: []string{},
		},
		{
			name: "main",
			content: `
# Contains profiles and configuration for full system confinement, only included
# when built with 'make full'
apparmor.d/groups/_full

apparmor.d/groups/apps # should be sandboxed
code
`,
			want: []string{
				"apparmor.d/groups/_full",
				"apparmor.d/groups/apps",
				"code",
			},
		},
	}
	IgnoreDir = paths.New("/tmp/")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IgnoreDir.Join(tt.name + ".ignore").WriteFile([]byte(tt.content))
			if err != nil {
				return
			}
			if got := Ignore.Read(tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ignore.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
