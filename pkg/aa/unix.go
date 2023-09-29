// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Unix struct {
	Qualifier
	Access   string
	Type     string
	Protocol string
	Address  string
	Label    string
	Attr     string
	Opt      string
	Peer     string
	PeerAddr string
}

func UnixFromLog(log map[string]string) ApparmorRule {
	return &Unix{
		Qualifier: NewQualifierFromLog(log),
		Access:    maskToAccess[log["requested_mask"]],
		Type:      log["sock_type"],
		Protocol:  log["protocol"],
		Address:   log["addr"],
		Label:     log["peer_label"],
		Attr:      log["attr"],
		Opt:       log["opt"],
		Peer:      log["peer"],
		PeerAddr:  log["peer_addr"],
	}
}

func (r *Unix) Less(other any) bool {
	o, _ := other.(*Unix)
	if r.Qualifier.Equals(o.Qualifier) {
		if r.Access == o.Access {
			if r.Type == o.Type {
				if r.Protocol == o.Protocol {
					if r.Address == o.Address {
						if r.Label == o.Label {
							if r.Attr == o.Attr {
								if r.Opt == o.Opt {
									if r.Peer == o.Peer {
										return r.PeerAddr < o.PeerAddr
									}
									return r.Peer < o.Peer
								}
								return r.Opt < o.Opt
							}
							return r.Attr < o.Attr
						}
						return r.Label < o.Label
					}
					return r.Address < o.Address
				}
				return r.Protocol < o.Protocol
			}
			return r.Type < o.Type
		}
		return r.Access < o.Access
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Unix) Equals(other any) bool {
	o, _ := other.(*Unix)
	return r.Access == o.Access && r.Type == o.Type &&
		r.Protocol == o.Protocol && r.Address == o.Address &&
		r.Label == o.Label && r.Attr == o.Attr && r.Opt == o.Opt &&
		r.Peer == o.Peer && r.PeerAddr == o.PeerAddr && r.Qualifier.Equals(o.Qualifier)
}
