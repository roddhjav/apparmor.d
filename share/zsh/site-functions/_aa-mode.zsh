#compdef aa-mode
#autoload

_aa-mode() {
	local IFS=$'\n'
	_arguments : \
		{-e,--enforce}'[set the profile in enforce mode]' \
		{-c,--complain}'[set the profile in complain mode]' \
		{-k,--kill}'[set the profile in kill mode]' \
		{-a,--default-allow}'[set the profile in default_allow mode]' \
		{-u,--unconfined}'[set the profile in unconfined mode]' \
		{-p,--prompt}'[set the profile in prompt mode]' \
		'--no-reload[do not reload the profile after modifying it]' \
		{-h,--help}'[display help information]'

	_values -C 'profile names' ${$(__aa_profiles):-""}
}

__aa_profiles() {
	find -L /etc/apparmor.d -maxdepth 1 -type f -printf '%P\n' | sort
}

_aa-mode
