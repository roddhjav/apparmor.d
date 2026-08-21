// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	cfg = tasks.NewTaskConfig(paths.New(".build"))
)

func TestBuilder_Apply(t *testing.T) {
	tests := []struct {
		name    string
		b       Builder
		profile string
		want    string
		wantErr bool
	}{
		{
			name: "abi4",
			b:    NewABI4(),
			profile: `
			  abi <abi/5.0>,
			  profile test {
			    if "gnome" in @{DE} {
			      /gnome r,
			    } else if "kde" in @{DE} {
			      /kde r,
			    } else {
			      /other r,
			    }
			  }`,
			want: `
			  abi <abi/4.0>,
			  profile test {
			      /gnome r,
			      /kde r,
			      /other r,
			  }`,
		},
		{
			name: "abi4 boolean variables",
			b:    NewABI4(),
			profile: `@{ABI} = 4
${FLATPAK_APPS} = true
#${RBAC} = false
`,
			want: `@{ABI} = 4
#${FLATPAK_APPS} = true
#${RBAC} = false
`,
		},
		{
			name: "complain-1",
			b:    NewDeployMode("complain", nil),
			profile: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
			want: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(complain) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
		},
		{
			name: "complain-2",
			b:    NewDeployMode("complain", nil),
			profile: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(complain) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
			want: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(complain) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
		},
		{
			name: "complain-3",
			b:    NewDeployMode("complain", nil),
			profile: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(attach_disconnected) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
			want: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(attach_disconnected,complain) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
		},
		{
			name: "complain-4",
			b:    NewDeployMode("complain", nil),
			profile: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(unconfined) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
			want: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(unconfined) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
		},
		{
			name: "enforce-1",
			b:    NewDeployMode("enforce", nil),
			profile: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
			want: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
		},
		{
			name: "enforce-2",
			b:    NewDeployMode("enforce", nil),
			profile: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(complain) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
			want: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
		},
		{
			name: "enforce-3",
			b:    NewDeployMode("enforce", nil),
			profile: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(attach_disconnected,complain) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
			want: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(attach_disconnected) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
		},
		{
			name:    "deploy-mode complain",
			b:       NewDeployMode("complain", nil),
			profile: "\nprofile foo /{,usr/}bin/foo {\n  include <abstractions/base>\n}",
			want:    "\nprofile foo /{,usr/}bin/foo flags=(complain) {\n  include <abstractions/base>\n}",
		},
		{
			name:    "deploy-mode enforce clears complain",
			b:       NewDeployMode("enforce", nil),
			profile: "\nprofile foo /{,usr/}bin/foo flags=(complain) {\n  include <abstractions/base>\n}",
			want:    "\nprofile foo /{,usr/}bin/foo {\n  include <abstractions/base>\n}",
		},
		{
			name:    "deploy-mode override wins",
			b:       NewDeployMode("complain", map[string]bool{"foo": true}),
			profile: "\nprofile foo /{,usr/}bin/foo {\n  include <abstractions/base>\n}",
			want:    "\nprofile foo /{,usr/}bin/foo {\n  include <abstractions/base>\n}",
		},
		{
			name:    "deploy-mode empty is noop",
			b:       NewDeployMode("", nil),
			profile: "\nprofile foo /{,usr/}bin/foo flags=(complain) {\n  include <abstractions/base>\n}",
			want:    "\nprofile foo /{,usr/}bin/foo flags=(complain) {\n  include <abstractions/base>\n}",
		},
		{
			name: "fsp",
			b:    NewFSP(),
			profile: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} {
			    include <abstractions/base>

			    @{exec_path} mr,
				@{bin}/* rPUx,
				@{lib}/* rUx,

				include if exists <local/foo>
			  }`,
			want: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} {
			    include <abstractions/base>

			    @{exec_path} mr,
				@{bin}/* rPx,
				@{lib}/* rPx,

				include if exists <local/foo>
			  }`,
		},
		{
			name: "userspace-1",
			b:    NewUserspace(),
			profile: `
			  @{exec_path}  = @{bin}/baloo_file @{lib}/{,kf6/}baloo_file
			  @{exec_path} += @{lib}/@{multiarch}/{,libexec/}baloo_file
			  profile baloo @{exec_path} {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/baloo>
			  }`,
			want: `
			  @{exec_path}  = @{bin}/baloo_file @{lib}/{,kf6/}baloo_file
			  @{exec_path} += @{lib}/@{multiarch}/{,libexec/}baloo_file
			  profile baloo /{{,usr/}bin/baloo_file,{,usr/}lib{,exec,32,64}/{,kf6/}baloo_file,{,usr/}lib{,exec,32,64}/*-linux-gnu*/{,libexec/}baloo_file} {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/baloo>
			  }`,
		},
		{
			name: "userspace-2",
			b:    NewUserspace(),
			profile: `
			  profile foo /usr/bin/foo {
			    include <abstractions/base>

			    /usr/bin/foo mr,

				include if exists <local/foo>
			  }`,
			want:    "",
			wantErr: true,
		},
		{
			name: "stacked-dbus-1",
			b:    NewStackedDbus(),
			profile: `
profile foo {
  dbus send bus=session path=/org/freedesktop/DBus
       interface=org.freedesktop.DBus
       member={Hello,AddMatch,RemoveMatch,GetNameOwner,NameHasOwner,StartServiceByName}
       peer=(name=org.freedesktop.DBus, label="@{p_dbus_session}"),

}`,
			want: `
profile foo {
dbus send bus=session path=/org/freedesktop/DBus
       interface=org.freedesktop.DBus
       member={Hello,AddMatch,RemoveMatch,GetNameOwner,NameHasOwner,StartServiceByName}
       peer=(name=org.freedesktop.DBus, label=dbus-session),
dbus send bus=session path=/org/freedesktop/DBus
       interface=org.freedesktop.DBus
       member={Hello,AddMatch,RemoveMatch,GetNameOwner,NameHasOwner,StartServiceByName}
       peer=(name=org.freedesktop.DBus, label=dbus-session//&unconfined),

}`,
		},
		{
			name: "base-strict-1",
			b:    NewBaseStrict(),
			profile: `
profile foo {
  include <abstractions/base>
}`,
			want: `
profile foo {
  include <abstractions/base-strict>
}`,
		},
		{
			name: "attach-1",
			b:    NewAttach(),
			profile: `
profile attach-1 flags=(attach_disconnected) {
  include <abstractions/base>
  include <abstractions/base-strict>
  include <abstractions/consoles>
}`,
			want: `
@{att} = /att/attach-1/
profile attach-1 flags=(attach_disconnected,attach_disconnected.path=@{att}) {
  include <abstractions/attached/base>
  include <abstractions/attached/base>
  include <abstractions/attached/consoles>
}`,
		},
		{
			name: "attach-2",
			b:    NewAttach(),
			profile: `
profile attach-2 flags=(complain) {
  include <abstractions/base>
  include <abstractions/base-strict>
  include <abstractions/consoles>
}`,
			want: `
@{att} = ""
profile attach-2 flags=(complain) {
  include <abstractions/base>
  include <abstractions/base-strict>
  include <abstractions/consoles>
}`,
		},
		{
			name: "attach-namespace",
			b:    NewAttach(),
			profile: `
profile :glycin:bwrap flags=(attach_disconnected) {
  include <abstractions/base>
}`,
			want: `
@{att} = /att/glycin/
profile :glycin:bwrap flags=(attach_disconnected,attach_disconnected.path=@{att}) {
  include <abstractions/attached/base>
}`,
		},
		{
			name: "debug-1",
			b:    NewDebug(),
			profile: `
profile debug-1 {
  include <abstractions/base>
  # @{exec_path} mr,
  audit @{bin}/ls Px,
  @{exec_path} mr,
  @{bin}/foo Px,
  @{bin}/bar ix,
}`,
			want: `
profile debug-1 {
  include <abstractions/base>
  # @{exec_path} mr,
  audit @{bin}/ls Px,
  @{exec_path} mr,
  audit @{bin}/foo Px,
  @{bin}/bar ix,
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := &Option{File: cfg.RootApparmor.Join(tt.name), Name: tt.name}
			tt.b.SetConfig(cfg)
			got, err := tt.b.Apply(opt, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOption(t *testing.T) {
	tests := []struct {
		name     string
		file     *paths.Path
		wantName string
		wantKind aa.FileKind
	}{
		{
			name:     "profile",
			file:     cfg.RootApparmor.Join("foo"),
			wantName: "foo",
			wantKind: aa.ProfileKind,
		},
		{
			name:     "profile-with-suffix",
			file:     cfg.RootApparmor.Join("bar.apparmor.d"),
			wantName: "bar",
			wantKind: aa.ProfileKind,
		},
		{
			name:     "abstraction",
			file:     cfg.RootApparmor.Join("abstractions", "app", "foo"),
			wantName: "foo",
			wantKind: aa.AbstractionKind,
		},
		{
			name:     "tunable",
			file:     cfg.RootApparmor.Join("tunables", "global"),
			wantName: "global",
			wantKind: aa.TunableKind,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOption(tt.file)
			if got.Name != tt.wantName {
				t.Errorf("NewOption() Name = %v, want %v", got.Name, tt.wantName)
			}
			if got.Kind != tt.wantKind {
				t.Errorf("NewOption() Kind = %v, want %v", got.Kind, tt.wantKind)
			}
			if got.File != tt.file {
				t.Errorf("NewOption() File = %v, want %v", got.File, tt.file)
			}
		})
	}
}

func TestBuilders_Run(t *testing.T) {
	tests := []struct {
		name     string
		builders []Builder
		profile  string
		want     string
		wantErr  bool
	}{
		{
			name:     "no-builder",
			builders: []Builder{},
			profile:  "profile foo /usr/bin/foo {\n}\n",
			want:     "profile foo /usr/bin/foo {\n}\n",
			wantErr:  false,
		},
		{
			name:     "complain-then-enforce",
			builders: []Builder{NewDeployMode("complain", nil), NewDeployMode("enforce", nil)},
			profile:  "profile foo /usr/bin/foo {\n}\n",
			want:     "profile foo /usr/bin/foo {\n}\n",
			wantErr:  false,
		},
		{
			name:     "complain",
			builders: []Builder{NewDeployMode("complain", nil)},
			profile:  "profile foo /usr/bin/foo {\n}\n",
			want:     "profile foo /usr/bin/foo flags=(complain) {\n}\n",
			wantErr:  false,
		},
		{
			name:     "error",
			builders: []Builder{NewUserspace()},
			profile:  "profile foo /usr/bin/foo {\n}\n",
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRunner(cfg)
			for _, b := range tt.builders {
				r.Add(b)
			}
			got, err := r.Run(cfg.RootApparmor.Join("foo"), tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builders.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), "userspace") {
				t.Errorf("Builders.Run() error = %v, want it to contain the builder name", err)
			}
			if got != tt.want {
				t.Errorf("Builders.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilders_Add(t *testing.T) {
	tests := []struct {
		name     string
		builders []Builder
		want     []string
	}{
		{
			name:     "add-builders",
			builders: []Builder{NewDeployMode("complain", nil), NewUserspace()},
			want:     []string{"deploy-mode", "userspace"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRunner(cfg)
			for _, b := range tt.builders {
				r.Add(b)
			}
			if len(r.Tasks) != len(tt.want) {
				t.Errorf("Builders.Add() len = %v, want %v", len(r.Tasks), len(tt.want))
			}
			for i, name := range tt.want {
				if r.Tasks[i].Name() != name {
					t.Errorf("Builders.Add() name = %v, want %v", r.Tasks[i].Name(), name)
				}
			}
		})
	}
}
