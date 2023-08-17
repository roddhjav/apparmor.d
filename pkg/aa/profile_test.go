// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"testing"
)

func TestAppArmorProfile_String(t *testing.T) {
	tests := []struct {
		name string
		p    AppArmorProfile
		want string
	}{
		{
			name: "empty",
			p:    AppArmorProfile{},
			want: `profile  {
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("AppArmorProfile.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
