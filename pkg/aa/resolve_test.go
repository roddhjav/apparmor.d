// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func TestAppArmorProfileFile_resolveInclude(t *testing.T) {
	tests := []struct {
		name    string
		include *Include
		want    *AppArmorProfileFile
		wantErr bool
	}{
		{
			name:    "empty",
			include: &Include{Path: "", IsMagic: true},
			want:    &AppArmorProfileFile{Preamble: Rules{&Include{Path: "", IsMagic: true}}},
			wantErr: true,
		},
		{
			name:    "tunables",
			include: &Include{Path: "tunables/global", IsMagic: true},
			want: &AppArmorProfileFile{
				Preamble: Rules{
					&Alias{Path: "/usr/", RewrittenPath: "/User/"},
					&Alias{Path: "/lib/", RewrittenPath: "/Libraries/"},
					&Comment{Base: Base{IsLineRule: true, Comment: " variable declarations for inclusion"}},
					&Variable{
						Name: "FOO", Define: true,
						Values: []string{
							"/foo", "/bar", "/baz", "/biff", "/lib", "/tmp",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	MagicRoot = paths.New("../../tests/testdata")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &AppArmorProfileFile{}
			got.Preamble = append(got.Preamble, tt.include)
			if err := got.resolveInclude(tt.include); (err != nil) != tt.wantErr {
				t.Errorf("AppArmorProfileFile.resolveInclude() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfileFile.resolveValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfileFile_resolveValues(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:    "not-defined",
			input:   "@{newvar}",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no-name",
			input:   "@{}",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "default",
			input: "@{etc_ro}",
			want:  []string{"/{,usr/}etc/"},
		},
		{
			name:  "simple",
			input: "@{bin}/foo",
			want:  []string{"/{,usr/}bin/foo"},
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
			got, err := f.resolveValues(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppArmorProfileFile.resolveValues() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfileFile.resolveValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfileFile_Resolve(t *testing.T) {
	tests := []struct {
		name         string
		preamble     Rules
		attachements []string
		want         *AppArmorProfileFile
		wantErr      bool
	}{
		{
			name: "variables/append",
			preamble: Rules{
				&Variable{Name: "lib", Values: []string{"/{usr/,}lib"}, Define: true},
				&Variable{Name: "multiarch", Values: []string{"*-linux-gnu*"}, Define: true},
				&Variable{Name: "exec_path", Values: []string{"@{lib}/DiscoverNotifier"}, Define: true},
				&Variable{Name: "exec_path", Values: []string{"@{lib}/@{multiarch}/{,libexec/}DiscoverNotifier"}, Define: false},
			},
			want: &AppArmorProfileFile{
				Preamble: Rules{
					&Variable{Name: "lib", Values: []string{"/{usr/,}lib"}, Define: true},
					&Variable{Name: "multiarch", Values: []string{"*-linux-gnu*"}, Define: true},
					&Variable{
						Name: "exec_path", Define: true,
						Values: []string{
							"/{usr/,}lib/DiscoverNotifier",
							"/{usr/,}lib/*-linux-gnu*/{,libexec/}DiscoverNotifier",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "attachment/firefox",
			preamble: Rules{
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
			name: "attachment/chromium",
			preamble: Rules{
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
			name: "attachment/geoclue",
			preamble: Rules{
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
			name: "attachment/opera",
			preamble: Rules{
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
			got := &AppArmorProfileFile{Preamble: tt.preamble}
			if tt.attachements != nil {
				got.Profiles = append(got.Profiles, &Profile{Header: Header{Attachments: tt.attachements}})
			}

			if err := got.Resolve(); (err != nil) != tt.wantErr {
				t.Errorf("AppArmorProfileFile.Resolve() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfile.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
