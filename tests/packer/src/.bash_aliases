#!/usr/bin/env bash

source /usr/share/bash-completion/bash_completion

function up() {
for nb in $(seq "$1"); do
    cd ../
done
}

alias aa-log='sudo aa-log'
alias aa-status='sudo aa-status'
alias c='clear'
alias du='du -hs'
alias l='ll -h'
alias ll='ls -alFh'
alias p="ps auxZ | grep -v '\[.*\]'"
alias pf="ps auxfZ | grep -v '\[.*\]'"
alias pu="ps auxZ | grep -v '\[.*\]' | grep unconfined"
alias u='up 1'
alias uu='up 2'
alias uuu='up 3'
alias uuuu='up 4'
alias uuuuu='up 5'