#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "busctl: Show all peers on the bus, by their service names" {
    busctl list
}

@test "busctl: Show process information and credentials of a bus service, a process, or the owner of the bus (if no parameter is specified)" {
    busctl status 1
    busctl status org.freedesktop.DBus
}

@test "busctl: Show an object tree of one or more services (or all services if no service is specified)" {
    busctl tree org.freedesktop.DBus
}

@test "busctl: Show interfaces, methods, properties and signals of the specified object on the specified service" {
    busctl introspect org.freedesktop.login1 /org/freedesktop/login1
}

@test "busctl: Retrieve the current value of one or more object properties" {
    busctl get-property org.freedesktop.login1 /org/freedesktop/login1 org.freedesktop.login1.Manager Docked
}
