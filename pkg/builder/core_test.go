// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"testing"

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
			name: "abi3",
			b:    NewABI3(),
			profile: `
			  abi <abi/4.0>,
			  profile test {
			    userns,
			    mqueue r type=posix /,
			  }`,
			want: `
			  abi <abi/3.0>,
			  profile test {
			    # userns,
			    # mqueue r type=posix /,
			  }`,
		},
		{
			name: "complain-1",
			b:    NewComplain(),
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
			b:    NewComplain(),
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
			b:    NewComplain(),
			profile: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(attach_disconnected) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
			want: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path}  flags=(attach_disconnected,complain) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
		},
		{
			name: "enforce-1",
			b:    NewEnforce(),
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
			b:    NewEnforce(),
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
			b:    NewEnforce(),
			profile: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path} flags=(attach_disconnected,complain) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
			want: `
			  @{exec_path} = @{bin}/foo
			  profile foo @{exec_path}  flags=(attach_disconnected) {
			    include <abstractions/base>

			    @{exec_path} mr,

				include if exists <local/foo>
			  }`,
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

func TestBuilders_Add(t *testing.T) {
	tests := []struct {
		name     string
		builders []Builder
		want     []string
	}{
		{
			name:     "add-builders",
			builders: []Builder{NewComplain(), NewEnforce()},
			want:     []string{"complain", "enforce"},
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
