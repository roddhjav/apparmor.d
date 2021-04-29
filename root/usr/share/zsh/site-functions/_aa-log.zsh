#compdef aa-log
#autoload

_aa-log () {
	local IFS=$'\n'
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
