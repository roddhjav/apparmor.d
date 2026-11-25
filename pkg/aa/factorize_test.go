// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func newTestNetwork(t *testing.T, domain, nType, lport, fport string, access ...string) *Network {
	t.Helper()
	return &Network{
		LocalAddress: LocalAddress{Port: lport},
		PeerAddress:  PeerAddress{Port: fport},
		Access:       MustAccess(NETWORK, access...),
		Domain:       domain,
		Type:         nType,
	}
}

func newTestNetworkInclude(t *testing.T, name string) *Include {
	t.Helper()
	return &Include{
		IsMagic: true,
		Path:    "abstractions/network/" + name,
	}
}

func TestRules_Factorize(t *testing.T) {
	tests := []struct {
		name  string
		rules func(t *testing.T) Rules
		want  func(t *testing.T) Rules
	}{
		{
			name: "netlink",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "netlink", "raw", "", "", "create"),
					newTestNetwork(t, "netlink", "raw", "", "", "bind"),
					newTestNetwork(t, "netlink", "raw", "", "", "send"),
					newTestNetwork(t, "netlink", "raw", "", "", "receive"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "netlink"),
				}
			},
		},
		{
			name: "netlink-kept",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "netlink", "raw", "", "", "listen"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "netlink", "raw", "", "", "listen"),
				}
			},
		},
		{
			name: "netlink-dgram-kept",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "netlink", "dgram", "", "", "receive"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "netlink", "dgram", "", "", "receive"),
				}
			},
		},
		{
			name: "netlink-udp-connect-kept",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "netlink", "raw", "", "", "create"),
					newTestNetwork(t, "netlink", "raw", "", "", "getattr"),
					newTestNetwork(t, "inet6", "dgram", "", "", "create"),
					newTestNetwork(t, "inet6", "dgram", "", "", "connect"),
					newTestNetwork(t, "inet6", "dgram", "53880", "22", "getattr"),
					newTestNetwork(t, "inet6", "dgram", "53880", "22", "connect"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet6", "dgram", "", "", "create"),
					newTestNetwork(t, "inet6", "dgram", "", "", "connect"),
					newTestNetwork(t, "inet6", "dgram", "53880", "22", "getattr"),
					newTestNetwork(t, "inet6", "dgram", "53880", "22", "connect"),
					newTestNetworkInclude(t, "netlink"),
				}
			},
		},
		{
			name: "udp-connect-with-udp-service",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "netlink", "raw", "", "", "create"),
					newTestNetwork(t, "inet", "dgram", "40001", "22", "connect"),
					newTestNetwork(t, "inet", "dgram", "40000", "53", "send"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "netlink"),
					newTestNetworkInclude(t, "out/domain"),
				}
			},
		},
		{
			name: "in-ssh",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet6", "stream", "", "", "create"),
					newTestNetwork(t, "inet6", "stream", "", "", "setopt"),
					newTestNetwork(t, "inet6", "stream", "22", "", "bind"),
					newTestNetwork(t, "inet6", "stream", "22", "", "listen"),
					newTestNetwork(t, "inet6", "stream", "22", "41562", "accept"),
					newTestNetwork(t, "inet6", "stream", "22", "41562", "receive"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "in/ssh"),
				}
			},
		},
		{
			name: "in-port",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "stream", "8080", "", "bind"),
					newTestNetwork(t, "inet", "stream", "8080", "", "listen"),
				}
			},
			want: func(t *testing.T) Rules {
				include := newTestNetworkInclude(t, "in/8080")
				include.IfExists = true
				include.Comment = " does not exist yet"
				return Rules{
					include,
				}
			},
		},
		{
			name: "out-https",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "stream", "", "", "create"),
					newTestNetwork(t, "inet", "stream", "", "443", "connect"),
					newTestNetwork(t, "inet", "stream", "51234", "443", "send"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "out/https"),
				}
			},
		},
		{
			name: "out-http-connect-without-address",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "stream", "", "", "create"),
					newTestNetwork(t, "inet", "stream", "", "", "connect"),
					newTestNetwork(t, "inet", "stream", "34888", "80", "send"),
					newTestNetwork(t, "inet", "stream", "34888", "80", "receive"),
					newTestNetwork(t, "inet", "stream", "34888", "80", "getopt"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "out/http"),
				}
			},
		},
		{
			name: "out-https-without-data",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "stream", "", "", "create"),
					newTestNetwork(t, "inet", "stream", "", "", "setopt"),
					newTestNetwork(t, "inet", "stream", "", "", "connect"),
					newTestNetwork(t, "inet", "stream", "37232", "443", "getattr"),
					newTestNetwork(t, "inet", "stream", "37232", "443", "getopt"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "out/https"),
				}
			},
		},
		{
			name: "out-dhcp",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "dgram", "", "", "create"),
					newTestNetwork(t, "inet", "dgram", "", "", "bind"),
					newTestNetwork(t, "inet", "dgram", "68", "", "connect"),
					newTestNetwork(t, "inet6", "dgram", "", "", "create"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "in/bootpc"),
				}
			},
		},
		{
			name: "udp-bind-kept",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "dgram", "", "", "bind"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "dgram", "", "", "bind"),
				}
			},
		},
		{
			name: "out-ntp-unconnected",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "dgram", "", "", "create"),
					newTestNetwork(t, "inet", "dgram", "", "", "connect"),
					newTestNetwork(t, "inet", "dgram", "56483", "123", "connect"),
					newTestNetwork(t, "inet", "dgram", "56483", "123", "getattr"),
					newTestNetwork(t, "inet", "dgram", "", "", "bind"),
					newTestNetwork(t, "inet", "dgram", "60368", "", "setopt"),
					newTestNetwork(t, "inet", "dgram", "60368", "", "send"),
					newTestNetwork(t, "inet", "dgram", "60368", "", "receive"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "out/ntp"),
				}
			},
		},
		{
			name: "unconnected-ambiguous-kept",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "dgram", "56483", "123", "connect"),
					newTestNetwork(t, "inet", "dgram", "56484", "53", "connect"),
					newTestNetwork(t, "inet", "dgram", "60368", "", "send"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "dgram", "56483", "123", "connect"),
					newTestNetwork(t, "inet", "dgram", "56484", "53", "connect"),
					newTestNetwork(t, "inet", "dgram", "60368", "", "send"),
				}
			},
		},
		{
			name: "out-dns-udp-connect",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "dgram", "", "", "create"),
					newTestNetwork(t, "inet", "dgram", "", "53", "connect"),
					newTestNetwork(t, "inet", "dgram", "40000", "53", "send"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "out/domain"),
				}
			},
		},
		{
			name: "out-dns-connect-without-port",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "dgram", "", "", "connect"),
					newTestNetwork(t, "inet", "dgram", "", "53", "connect"),
					newTestNetwork(t, "inet", "dgram", "40000", "53", "send"),
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "out/domain"),
				}
			},
		},
		{
			name: "out-ipv6-icmp",
			rules: func(t *testing.T) Rules {
				return Rules{
					&Network{Access: MustAccess(NETWORK, "create"), Domain: "inet6", Type: "raw", Protocol: "58"},
					&Network{LocalAddress: LocalAddress{Port: "58"}, Access: MustAccess(NETWORK, "setopt"), Domain: "inet6", Type: "raw", Protocol: "58"},
					&Network{LocalAddress: LocalAddress{Port: "58"}, Access: MustAccess(NETWORK, "send"), Domain: "inet6", Type: "raw", Protocol: "58"},
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetworkInclude(t, "out/ipv6-icmp"),
				}
			},
		},
		{
			name: "raw-without-traffic-kept",
			rules: func(t *testing.T) Rules {
				return Rules{
					&Network{Access: MustAccess(NETWORK, "create"), Domain: "inet", Type: "raw", Protocol: "1"},
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					&Network{Access: MustAccess(NETWORK, "create"), Domain: "inet", Type: "raw", Protocol: "1"},
				}
			},
		},
		{
			name: "kept",
			rules: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "stream", "", "", "create"),
					newTestNetwork(t, "inet", "stream", "22", "41562", "send"),
					&Capability{Names: []string{"net_admin"}},
				}
			},
			want: func(t *testing.T) Rules {
				return Rules{
					newTestNetwork(t, "inet", "stream", "", "", "create"),
					newTestNetwork(t, "inet", "stream", "22", "41562", "send"),
					&Capability{Names: []string{"net_admin"}},
				}
			},
		},
	}
	magicRoot, servicesFile, protocolsFile := MagicRoot, ServicesFile, ProtocolsFile
	defer func() { MagicRoot, ServicesFile, ProtocolsFile = magicRoot, servicesFile, protocolsFile }()
	MagicRoot = paths.New("../../apparmor.d")
	ServicesFile = paths.New("../../tests/testdata/services")
	ProtocolsFile = paths.New("../../tests/testdata/protocols")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rules(t).Factorize()
			want := tt.want(t)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Rules.Factorize() = %v, want %v", got, want)
			}
		})
	}
}
