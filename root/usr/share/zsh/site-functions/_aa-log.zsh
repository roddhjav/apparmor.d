#compdef aa-log
#autoload

_aa-log() {
	local IFS=$'\n'
	_arguments : \
		-f'[set a logfile or a prefix to the default log file]:FILE:__aa_files' \
		-s'[parse systemd dbus logs]' \
		-h'[display help information]'

	_values -C 'profile names' ${$(__aa_profiles):-""}
}

__aa_files() {
	find /var/log/audit/ -type f -printf '%P\n' | cut -d '.' -f 3
	_files
}

__aa_profiles() {
	find -L /etc/apparmor.d -maxdepth 1 -type f -printf '%P\n' | sort
}

_aa-log
