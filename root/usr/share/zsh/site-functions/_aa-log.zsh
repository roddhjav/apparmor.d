#compdef aa-log
#autoload

_aa-log () {
	local IFS=$'\n'
	_arguments : \
		-f'[set a logfile or a prefix to the default log file]:_files' \
		-h'[display help information]'

	_values -C 'profile names' ${$(__aa_profiles):-""}
}

__aa_profiles() {
	find -L /etc/apparmor.d -type f -printf '%P\n' \
		| sed -e '/abi/d' \
			  -e '/abstractions/d' \
			  -e '/local/d' \
			  -e '/tunables/d' \
		| sort
}

_aa-log
