// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"slices"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
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
			b:    Builders["abi3"],
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
			b:    Builders["complain"],
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
			b:    Builders["complain"],
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
			b:    Builders["complain"],
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
			b:    Builders["enforce"],
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
			b:    Builders["enforce"],
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
			name: "complain-3",
			b:    Builders["enforce"],
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
			b:    Builders["fsp"],
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
			b:    Builders["userspace"],
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
			b:    Builders["userspace"],
			profile: `
			  profile foo /usr/bin/foo {
			    include <abstractions/base>

			    /usr/bin/foo mr,

				include if exists <local/foo>
			  }`,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := &Option{File: prebuild.RootApparmord.Join(tt.name)}
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

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		names       []string
		wantSuccess bool
	}{
		{
			name:        "test",
			names:       []string{"complain", "enforce"},
			wantSuccess: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.names...)
			for _, name := range tt.names {
				if got := slices.Contains(Builds, Builders[name]); got != tt.wantSuccess {
					t.Errorf("Register() = %v, want %v", got, tt.wantSuccess)
				}

			}
		})
	}
}
