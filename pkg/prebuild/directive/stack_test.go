// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
)

func TestStack_Apply(t *testing.T) {
	tests := []struct {
		name          string
		rootApparmord *paths.Path
		opt           *Option
		profile       string
		want          string
		wantErr       bool
	}{
		{
			name:          "stack",
			rootApparmord: paths.New("../../../apparmor.d/groups/freedesktop/"),
			opt: &Option{
				Name:    "stack",
				ArgMap:  map[string]string{"plymouth": ""},
				ArgList: []string{"plymouth"},
				File:    nil,
				Raw:     "  #aa:stack plymouth",
			},
			profile: `
profile parent @{exec_path} {
  include <abstractions/base>

  @{exec_path} mr,

  #aa:stack plymouth
  @{bin}/plymouth rPx -> parent//&plymouth,

  @{PROC}/cmdline r,

  include if exists <local/parent>
}`,
			want: `
profile parent @{exec_path} {
  include <abstractions/base>

  @{exec_path} mr,


  @{bin}/plymouth rPx -> parent//&plymouth,

  @{PROC}/cmdline r,

  # Stacked profile: plymouth
  include <abstractions/fonts>
  include <abstractions/fontconfig-cache-read>
  include <abstractions/consoles>
  unix (send, receive, connect) type=stream peer=(addr="@/org/freedesktop/plymouthd"),
  @{PROC}/cmdline r,
  include if exists <local/plymouth>

  include if exists <local/parent>
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.RootApparmord = tt.rootApparmord
			got, err := Directives["stack"].Apply(tt.opt, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stack.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Stack.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
