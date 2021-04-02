#compdef aa-log
#autoload

_aa-log () {
	local IFS=$'\n'
	_values -C 'profile names' ${$(__aa_profiles):-""}
}

__aa_profiles() {
	find -L /etc/apparmor.d -type f \
		| sed -e 's#/etc/apparmor.d/##' \
			  -e '/abi/d' \
			  -e '/abstractions/d' \
			  -e '/local/d' \
		| sort
}

_aa-log
