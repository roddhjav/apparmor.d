#compdef aa-log
#autoload

_aa-log() {
	local IFS=$'\n'
	_arguments : \
		{-f,--file}'[set a logfile or a suffix to the default log file]:FILE:__aa_files' \
		{-s,--systemd}'[parse systemd logs from journalctl]' \
		{-r,--rules}'[convert the log into AppArmor rules]' \
		{-R,--raw}'[print the raw log without any formatting]' \
		{-h,--help}'[display help information]'

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
