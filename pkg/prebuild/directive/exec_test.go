// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
)

func TestExec_Apply(t *testing.T) {
	tests := []struct {
		name          string
		rootApparmord *paths.Path
		opt           *Option
		profile       string
		want          string
		wantErr       bool
	}{
		{
			name:          "exec",
			rootApparmord: paths.New("../../../apparmor.d/groups/kde/"),
			opt: &Option{
				Name:    "exec",
				ArgMap:  map[string]string{"DiscoverNotifier": ""},
				ArgList: []string{"DiscoverNotifier"},
				File:    nil,
				Raw:     "  #aa:exec DiscoverNotifier",
			},
			profile: `  #aa:exec DiscoverNotifier`,
			want: `  @{lib}/@{multiarch}/{,libexec/}DiscoverNotifier Px,
  @{lib}/DiscoverNotifier Px,`,
		},
		{
			name:          "exec-unconfined",
			rootApparmord: paths.New("../../../apparmor.d/groups/freedesktop/"),
			opt: &Option{
				Name:    "exec",
				ArgMap:  map[string]string{"U": "", "polkit-agent-helper": ""},
				ArgList: []string{"U", "polkit-agent-helper"},
				File:    nil,
				Raw:     "  #aa:exec U polkit-agent-helper",
			},
			profile: `  #aa:exec U polkit-agent-helper`,
			want: `  @{lib}/polkit-[0-9]/polkit-agent-helper-[0-9] Ux,
  @{lib}/polkit-agent-helper-[0-9] Ux,`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.RootApparmord = tt.rootApparmord
			got, err := Directives["exec"].Apply(tt.opt, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exec.Apply() = |%v|, want |%v|", got, tt.want)
			}
		})
	}
}
