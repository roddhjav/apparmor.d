// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"reflect"
	"regexp"
	"testing"
)

func TestDecodeHexInString(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{
			name: "Hexa",
			str:  `apparmor="ALLOWED" operation="rename_dest" parent=6974 profile="/usr/sbin/httpd2-prefork//vhost_foo" name=2F686F6D652F7777772F666F6F2E6261722E696E2F68747470646F63732F61707061726D6F722F696D616765732F746573742F696D61676520312E6A7067 pid=20143 comm="httpd2-prefork" requested_mask="wc"`,
			want: `apparmor="ALLOWED" operation="rename_dest" parent=6974 profile="/usr/sbin/httpd2-prefork//vhost_foo" name="/home/www/foo.bar.in/httpdocs/apparmor/images/test/image 1.jpg" pid=20143 comm="httpd2-prefork" requested_mask="wc"`,
		},
		{
			name: "Hexa abstract addr",
			str:  `apparmor="DENIED" operation="bind" class="net" profile="gjs" comm="gst-plugin-scan" family="unix" sock_type="seqpacket" protocol=0 requested_mask="bind" denied_mask="bind" addr="@637564612D75766D66642D343032363533313833362D3332373100"`,
			want: `apparmor="DENIED" operation="bind" class="net" profile="gjs" comm="gst-plugin-scan" family="unix" sock_type="seqpacket" protocol=0 requested_mask="bind" denied_mask="bind" addr="@cuda-uvmfd-4026531836-3271?"`,
		},
		{
			name: "Hexa peer_addr",
			str:  `apparmor="DENIED" operation="connect" class="net" profile="foo" family="unix" sock_type="stream" addr=none peer_addr=2F746D702F6D7920736F636B`,
			want: `apparmor="DENIED" operation="connect" class="net" profile="foo" family="unix" sock_type="stream" addr=none peer_addr="/tmp/my sock"`,
		},
		{
			name: "Hexa looking abstract addr",
			str:  `apparmor="DENIED" operation="bind" class="net" profile="foo" family="unix" sock_type="stream" addr="@BEEF"`,
			want: `apparmor="DENIED" operation="bind" class="net" profile="foo" family="unix" sock_type="stream" addr="@BEEF"`,
		},
		{
			name: "Not Hexa addr",
			str:  `apparmor="DENIED" operation="connect" class="net" profile="foo" family="unix" sock_type="stream" addr="@/tmp/.X11-unix/X0" laddr=127.0.0.1`,
			want: `apparmor="DENIED" operation="connect" class="net" profile="foo" family="unix" sock_type="stream" addr="@/tmp/.X11-unix/X0" laddr=127.0.0.1`,
		},
		{
			name: "Not Hexa",
			str:  `type=AVC msg=audit(1424425690.883:716630): apparmor="ALLOWED" operation="file_mmap" info="Failed name lookup - disconnected path" error=-13 profile="/sbin/klogd" name="var/run/nscd/passwd" pid=25333 comm="id" requested_mask="r" denied_mask="r" fsuid=1002 ouid=0`,
			want: `type=AVC msg=audit(1424425690.883:716630): apparmor="ALLOWED" operation="file_mmap" info="Failed name lookup - disconnected path" error=-13 profile="/sbin/klogd" name="var/run/nscd/passwd" pid=25333 comm="id" requested_mask="r" denied_mask="r" fsuid=1002 ouid=0`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeHexInString(tt.str); got != tt.want {
				t.Errorf("DecodeHexInString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToRegexRepl(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want RegexReplList
	}{
		{
			name: "",
			in: []string{
				"^/foo/bar", "/foo/bar",
				"^/foo/bar", "/foo/bar",
			},
			want: []RegexRepl{
				{Regex: regexp.MustCompile("^/foo/bar"), Repl: "/foo/bar"},
				{Regex: regexp.MustCompile("^/foo/bar"), Repl: "/foo/bar"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToRegexRepl(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToRegexRepl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexReplList_Replace(t *testing.T) {
	tests := []struct {
		name string
		rr   RegexReplList
		str  string
		want string
	}{
		{
			name: "default",
			rr: []RegexRepl{
				{Regex: regexp.MustCompile(`^/foo`), Repl: "/bar"},
			},
			str:  "/foo",
			want: "/bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rr.Replace(tt.str); got != tt.want {
				t.Errorf("RegexReplList.Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

