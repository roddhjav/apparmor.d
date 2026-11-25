// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"maps"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

var (
	ServicesFile  = paths.New("/etc/services")
	ProtocolsFile = paths.New("/etc/protocols")

	serviceNames  = sync.OnceValue(func() map[string]string { return loadNames(ServicesFile) })
	protocolNames = sync.OnceValue(func() map[string]string { return loadNames(ProtocolsFile) })
)

// loadNames returns the names from a services or protocols database,
// indexed by their second field: "<port>/<protocol>" or "<number>".
func loadNames(file *paths.Path) map[string]string {
	res := map[string]string{}
	if !file.Exist() {
		return res
	}
	for _, line := range file.MustReadFilteredFileAsLines() {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if _, ok := res[fields[1]]; !ok {
			res[fields[1]] = fields[0]
		}
	}
	return res
}

// networkProtocol returns the protocol (tcp, udp) of an inet/inet6 rule, or
// an empty string if the rule cannot be part of a network abstraction.
func networkProtocol(r *Network) string {
	if r.Domain != "inet" && r.Domain != "inet6" {
		return ""
	}
	switch r.Type {
	case "stream":
		return "tcp"
	case "dgram":
		return "udp"
	}
	return ""
}

// isEphemeralPort reports whether port is in the range the kernel uses for
// automatic local ports.
func isEphemeralPort(port string) bool {
	p, err := strconv.Atoi(port)
	return err == nil && p >= 32768 && p <= 60999
}

// networkService returns the service name of a port, or the port itself.
func networkService(protocol, port string) string {
	if name, ok := serviceNames()[port+"/"+protocol]; ok {
		return name
	}
	return port
}

// factorizeNetwork replaces the network rules by the network abstractions
// providing them: network/netlink, network/in/<service>, and
// network/out/<service>. Rules that do not match any abstraction are kept.
//
// Beware: an abstraction is proposed even if it does not exist yet.
//
// Logs do not identify sockets, so rules are matched against the whole set of
// network rules: an UDP port with send or receive is a real egress service,
// otherwise a UDP connect to a port without data is part of the network/udp
// base.
// A UDP connect without port is part of the network/udp base. A UDP socket on
// a fixed local port without peer (e.g. a DHCP client) is ingress. An
// unconnected UDP client (automatic local port, no peer logged) uses the peer
// port of its UDP connect. A raw socket is named after its protocol.
func (r Rules) factorizeNetwork() Rules {
	var (
		base    = MustAccess(NETWORK, "create", "getattr", "getopt", "setattr", "setopt", "shutdown")
		netlink = MustAccess(NETWORK, "create", "bind", "getattr", "getopt", "setopt", "send", "receive")
		connect = MustAccess(NETWORK, "connect")
		udpBase = base | connect
		ingress = udpBase | MustAccess(NETWORK, "accept", "bind", "listen", "send", "receive")
		egress  = base | connect | MustAccess(NETWORK, "send", "receive")
		bound   = MustAccess(NETWORK, "accept", "bind", "listen")
		data    = MustAccess(NETWORK, "send", "receive")
	)

	type socket struct {
		i        int
		n        *Network
		protocol string
	}
	type raw struct {
		idx    []int
		access AccessMask
	}
	includes := map[string]bool{}
	sockets := []socket{}
	raws := map[string]*raw{}

	// First pass: replace the netlink sockets, gather the raw sockets, and find
	// the local ports the profile serves and the remote ports it talks to.
	inPorts := map[string]bool{}
	outPorts := map[string]bool{}
	connectPorts := map[string]bool{}
	unconnected := false
	for i, rule := range r {
		n, ok := rule.(*Network)
		if !ok || n.Qualifier != (Qualifier{}) {
			continue
		}
		if n.Domain == "netlink" {
			if n.Type == "raw" && n.Access.Within(netlink) {
				includes["netlink"] = true
				r[i] = nil
			}
			continue
		}
		if n.Type == "raw" && (n.Domain == "inet" || n.Domain == "inet6") && n.Protocol != "" &&
			n.Access.Within(ingress) {
			if raws[n.Protocol] == nil {
				raws[n.Protocol] = &raw{}
			}
			raws[n.Protocol].idx = append(raws[n.Protocol].idx, i)
			raws[n.Protocol].access |= n.Access
			continue
		}
		protocol := networkProtocol(n)
		if protocol == "" {
			continue
		}
		sockets = append(sockets, socket{i, n, protocol})

		lport, fport := n.LocalAddress.Port, n.PeerAddress.Port
		if lport != "" && (n.Access.ContainsAny(bound) ||
			(protocol == "udp" && fport == "" && !isEphemeralPort(lport))) {
			inPorts[protocol+"/"+lport] = true
		}
		if fport == "" {
			// An unconnected UDP client sends from an automatic local port,
			// its peer is not logged.
			if protocol == "udp" && n.Access.ContainsAny(data) && isEphemeralPort(lport) {
				unconnected = true
			}
			continue
		}
		if protocol == "udp" && n.Access.Has(connect) {
			connectPorts[fport] = true
		}
		// A connected socket from an automatic local port is a client, even
		// when its connect was logged without address (failed af match). A TCP
		// socket with a peer is connected, whatever the logged access. An UDP one
		// without data is part of the network/udp base.
		client := (protocol == "tcp" || n.Access.ContainsAny(data)) &&
			isEphemeralPort(lport) && !isEphemeralPort(fport)
		if client || (protocol == "tcp" && n.Access.Has(connect)) ||
			(protocol == "udp" && n.Access.ContainsAny(data)) {
			outPorts[protocol+"/"+fport] = true
		}
	}

	// The UDP connect of an unconnected client reveals its peer port. It is
	// only used when there is a single one.
	clientPort := ""
	for port := range connectPorts {
		if outPorts["udp/"+port] {
			delete(connectPorts, port)
		}
	}
	if unconnected && len(connectPorts) == 1 {
		for port := range connectPorts {
			clientPort = port
		}
	}

	// Second pass: replace the tcp and udp sockets by the abstractions. A rule
	// only provided by a served abstraction is pending until all sockets are
	// known.
	served := map[string]bool{}   // By "<protocol>" and "<dir>/<protocol>"
	pending := map[string][]int{} // Same keys as served
	add := func(dir, protocol, port string) {
		includes[dir+"/"+networkService(protocol, port)] = true
		served[protocol], served[dir+"/"+protocol] = true, true
	}
	for _, s := range sockets {
		n, protocol := s.n, s.protocol
		lport, fport := n.LocalAddress.Port, n.PeerAddress.Port
		switch {
		case inPorts[protocol+"/"+lport] && n.Access.Within(ingress):
			add("in", protocol, lport)
			r[s.i] = nil
		case clientPort != "" && protocol == "udp" && fport == "" &&
			isEphemeralPort(lport) && n.Access.Within(udpBase|data):
			add("out", protocol, clientPort)
			r[s.i] = nil
		case outPorts[protocol+"/"+fport] && n.Access.Within(egress):
			add("out", protocol, fport)
			r[s.i] = nil
		case (protocol == "tcp" && n.Access.Within(base)) || (protocol == "udp" && n.Access.Within(udpBase)):
			// The local socket operations, provided by network/tcp and network/udp.
			pending[protocol] = append(pending[protocol], s.i)
		case protocol == "udp" && (lport == "" || lport == "0") &&
			n.Access.Within(udpBase|MustAccess(NETWORK, "bind")):
			// A bind to an automatic local port, provided by network/udp.
			pending[protocol] = append(pending[protocol], s.i)
		case protocol == "tcp" && fport == "" && n.Access.Within(egress):
			// A connect logged without address, provided by the out/ abstractions.
			pending["out/tcp"] = append(pending["out/tcp"], s.i)
		}
	}
	for key, idx := range pending {
		if !served[key] {
			continue
		}
		for _, i := range idx {
			r[i] = nil
		}
	}

	// A raw socket is named after its protocol. It is egress when it sends,
	// ingress when it only receives.
	for number, raw := range raws {
		name, ok := protocolNames()[number]
		if !ok {
			name = number
		}
		switch {
		case raw.access.ContainsAny(connect | MustAccess(NETWORK, "send")):
			includes["out/"+name] = true
		case raw.access.ContainsAny(MustAccess(NETWORK, "bind", "receive")):
			includes["in/"+name] = true
		default:
			continue
		}
		for _, i := range raw.idx {
			r[i] = nil
		}
	}

	r = slices.DeleteFunc(r, func(a Rule) bool { return a == nil })
	names := slices.Sorted(maps.Keys(includes))
	for _, name := range names {
		include := &Include{IsMagic: true, Path: "abstractions/network/" + name}
		if !MagicRoot.Join(include.Path).Exist() {
			include.IfExists = true
			include.Comment = " does not exist yet"
		}
		r = append(r, include)
	}
	return r
}
