// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"
)

func TestAppArmorProfileFile_resolveVariable(t *testing.T) {
	tests := []struct {
		name  string
		f     AppArmorProfileFile
		input string
		want  []string
	}{
		{
			name:  "nil",
			input: "@{newvar}",
			want:  []string{},
		},
		{
			name:  "empty",
			input: "@{}",
			want:  []string{"@{}"},
		},
		{
			name:  "default",
			input: "@{etc_ro}",
			want:  []string{"/{,usr/}etc/"},
		},
		{
			name:  "simple",
			input: "@{bin}/foo",
			want:  []string{"/{,usr/}{,s}bin/foo"},
		},
		{
			name:  "double",
			input: "@{lib}/@{multiarch}",
			want:  []string{"/{,usr/}lib{,exec,32,64}/*-linux-gnu*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := DefaultTunables()
			got := f.resolveVariable(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfileFile.resolveVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfileFile_Resolve(t *testing.T) {
	tests := []struct {
		name         string
		variables    Rules
		attachements []string
		want         *AppArmorProfileFile
		wantErr      bool
	}{
		{
			name: "firefox",
			variables: Rules{
				&Variable{Name: "firefox_name", Values: []string{"firefox{,-esr,-bin}"}, Define: true},
				&Variable{Name: "firefox_lib_dirs", Values: []string{"/{usr/,}/lib{,32,64}/@{firefox_name}", "/opt/@{firefox_name}"}, Define: true},
				&Variable{Name: "exec_path", Values: []string{"/{usr/,}bin/@{firefox_name}", "@{firefox_lib_dirs}/@{firefox_name}"}, Define: true},
			},
			attachements: []string{"@{exec_path}"},
			want: &AppArmorProfileFile{
				Preamble: Rules{
					&Variable{Name: "firefox_name", Values: []string{"firefox{,-esr,-bin}"}, Define: true},
					&Variable{
						Name: "firefox_lib_dirs", Define: true,
						Values: []string{
							"/{usr/,}/lib{,32,64}/firefox{,-esr,-bin}",
							"/opt/firefox{,-esr,-bin}",
						},
					},
					&Variable{
						Name: "exec_path", Define: true,
						Values: []string{
							"/{usr/,}bin/firefox{,-esr,-bin}",
							"/{usr/,}/lib{,32,64}/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
							"/opt/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
						},
					},
				},
				Profiles: []*Profile{
					{Header: Header{
						Attachments: []string{
							"/{usr/,}bin/firefox{,-esr,-bin}",
							"/{usr/,}/lib{,32,64}/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
							"/opt/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
						},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "chromium",
			variables: Rules{
				&Variable{Name: "name", Values: []string{"chromium"}, Define: true},
				&Variable{Name: "lib_dirs", Values: []string{"/{usr/,}lib/@{name}"}, Define: true},
				&Variable{Name: "path", Values: []string{"@{lib_dirs}/@{name}"}, Define: true},
			},
			attachements: []string{"@{path}/pass"},
			want: &AppArmorProfileFile{
				Preamble: Rules{
					&Variable{Name: "name", Values: []string{"chromium"}, Define: true},
					&Variable{Name: "lib_dirs", Values: []string{"/{usr/,}lib/chromium"}, Define: true},
					&Variable{Name: "path", Values: []string{"/{usr/,}lib/chromium/chromium"}, Define: true},
				},
				Profiles: []*Profile{
					{Header: Header{
						Attachments: []string{"/{usr/,}lib/chromium/chromium/pass"},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "geoclue",
			variables: Rules{
				&Variable{Name: "libexec", Values: []string{"/{usr/,}libexec"}, Define: true},
				&Variable{Name: "exec_path", Values: []string{"@{libexec}/geoclue", "@{libexec}/geoclue-2.0/demos/agent"}, Define: true},
			},
			attachements: []string{"@{exec_path}"},
			want: &AppArmorProfileFile{
				Preamble: Rules{
					&Variable{Name: "libexec", Values: []string{"/{usr/,}libexec"}, Define: true},
					&Variable{
						Name: "exec_path", Define: true,
						Values: []string{
							"/{usr/,}libexec/geoclue",
							"/{usr/,}libexec/geoclue-2.0/demos/agent",
						},
					},
				},
				Profiles: []*Profile{
					{Header: Header{
						Attachments: []string{
							"/{usr/,}libexec/geoclue",
							"/{usr/,}libexec/geoclue-2.0/demos/agent",
						},
					}},
				},
			},
			wantErr: false,
		},
		{
			name: "opera",
			variables: Rules{
				&Variable{Name: "multiarch", Values: []string{"*-linux-gnu*"}, Define: true},
				&Variable{Name: "name", Values: []string{"opera{,-beta,-developer}"}, Define: true},
				&Variable{Name: "lib_dirs", Values: []string{"/{usr/,}lib/@{multiarch}/@{name}"}, Define: true},
				&Variable{Name: "exec_path", Values: []string{"@{lib_dirs}/@{name}"}, Define: true},
			},
			attachements: []string{"@{exec_path}"},
			want: &AppArmorProfileFile{
				Preamble: Rules{
					&Variable{Name: "multiarch", Values: []string{"*-linux-gnu*"}, Define: true},
					&Variable{Name: "name", Values: []string{"opera{,-beta,-developer}"}, Define: true},
					&Variable{Name: "lib_dirs", Values: []string{"/{usr/,}lib/*-linux-gnu*/opera{,-beta,-developer}"}, Define: true},
					&Variable{Name: "exec_path", Values: []string{"/{usr/,}lib/*-linux-gnu*/opera{,-beta,-developer}/opera{,-beta,-developer}"}, Define: true},
				},
				Profiles: []*Profile{
					{Header: Header{
						Attachments: []string{
							"/{usr/,}lib/*-linux-gnu*/opera{,-beta,-developer}/opera{,-beta,-developer}",
						},
					}},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &AppArmorProfileFile{
				Profiles: []*Profile{{
					Header: Header{Attachments: tt.attachements},
				}},
			}
			got.Preamble = tt.variables
			if err := got.Resolve(); (err != nil) != tt.wantErr {
				t.Errorf("AppArmorProfileFile.Resolve() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfile.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
