// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
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
			want: `  /{,usr/}lib{,exec,32,64}/*-linux-gnu*/{,libexec/}DiscoverNotifier Px,
  /{,usr/}lib{,exec,32,64}/DiscoverNotifier Px,`,
		},
		{
			name:          "exec-unconfined",
			rootApparmord: paths.New("../../../apparmor.d/groups/polkit/"),
			opt: &Option{
				Name:    "exec",
				ArgMap:  map[string]string{"U": "", "polkit-agent-helper": ""},
				ArgList: []string{"U", "polkit-agent-helper"},
				File:    nil,
				Raw:     "  #aa:exec U polkit-agent-helper",
			},
			profile: `  #aa:exec U polkit-agent-helper`,
			want: `  /{,usr/}lib{,exec,32,64}/polkit-[0-9]/polkit-agent-helper-[0-9] Ux,
  /{,usr/}lib{,exec,32,64}/polkit-agent-helper-[0-9] Ux,`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prebuild.RootApparmord = tt.rootApparmord
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
