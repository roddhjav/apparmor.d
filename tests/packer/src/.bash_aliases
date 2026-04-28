#!/usr/bin/env bash

source /usr/share/bash-completion/bash_completion

function up() {
	for nb in $(seq "$1"); do
		cd ../
	done
}

function _ps() {
	LIBPROC_HIDE_KERNEL=1 ps "$@" | sed \
		-e "s/\bunconfined\b/\x1b[1;41;97munconfined\x1b[0m/g" \
		-e "s/\bcomplain\b/\x1b[1;33mcomplain\x1b[0m/g" \
		-e "s/\benforce\b/\x1b[1;32menforce\x1b[0m/g"
}

alias aa-log='sudo aa-log'
alias aa-status='sudo aa-status'
alias c='clear'
alias du='du -hs'
alias l='ll -h'
alias ll='ls -alFh'
alias p="_ps auxZ"
alias pf="_ps auxfZ"
alias pu="_ps auxZ | grep unconfined"
alias u='up 1'
alias uu='up 2'
alias uuu='up 3'
alias uuuu='up 4'
alias uuuuu='up 5'