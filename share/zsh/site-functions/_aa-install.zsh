#compdef aa-install
#autoload

_aa-install() {
	local IFS=$'\n'
	_arguments : \
		{-s,--status}'[show the installation status summary]' \
		{-l,--list}'[list installed profile paths from the manifest]' \
		{-i,--install}'[install the profiles]' \
		{-a,--all}'[install all profiles, even for programs not installed]' \
		{-c,--complain}'[set complain flag on all the profiles]' \
		{-e,--enforce}'[set enforce flag on all the profiles]' \
		{-u,--uninstall}'[remove all profiles installed]' \
		'--no-reload[do not reload the profiles after modifying them]' \
		'--config[select an alternate configuration directory]:DIR:_files -/' \
		'--magic[select an alternate apparmor.d directory]:DIR:_files -/' \
		'--src[select an alternate source directory]:DIR:_files -/' \
		{-h,--help}'[display help information]'
}

_aa-install
