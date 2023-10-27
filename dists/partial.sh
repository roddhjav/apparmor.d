BUILD=.build
DESTDIR=/
        
for profile in "$@"
do
  echo "Installing profile $profile"
  cp $BUILD/apparmor.d/$profile $DESTDIR/etc/apparmor.d/
  grep "rPx," "${BUILD}/apparmor.d/${profile}" | while read line
  do
    dep=$(echo "$l1" | awk '{print $1}')
    dep=$(echo $dep | awk -F"/" '{print $NF}')
    dep=$(eval "ls ${BUILD}/apparmor.d/${dep} 2>/dev/null")
  	for i in $dep
  	do
  	  i=$(echo $i | awk -F"/" '{print $NF}')
  	  if [ ! -f "$DESTDIR/etc/apparmor.d/$i" ]; then
        bash "$0" "$i"
      fi
	  done
  done
  grep "rPx -> " "${BUILD}/apparmor.d/${profile}" | while read line
  do
    dep=${line%%#*}
    dep=$(echo $dep | awk '{print $NF}')
    dep=${dep::-1}
    dep=$(eval "ls ${BUILD}/apparmor.d/${dep} 2>/dev/null")
	  for i in $dep
	  do
	    i=$(echo $i | awk -F"/" '{print $NF}')
	    if [ ! -f "$DESTDIR/etc/apparmor.d/$i" ]; then
			  bash "$0" "$i"
      fi
    done
  done
done
