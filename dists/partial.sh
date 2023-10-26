BUILD=.build
DESTDIR=/
for profile in "$@"
do
  cp $BUILD/apparmor.d/$profile $DESTDIR/etc/apparmor.d/$profile
  grep "rPx," "$BUILD/apparmor.d/$profile" | while read l1
  do
    dep=$(echo "$l1" | awk '{print $1}')
      dep=$(echo $dep | awk -F"/" '{print $NF}')
    find . -type f  -name $dep | while read l2
        do
      if [ ! -f "$DESTDIR/etc/apparmor.d/$dep" ]; then
        install_seperate_with_depends $dep
      fi
        done
  done
  grep "rPx -> " $BUILD/apparmor.d/$profile | while read l1
  do
    dep=$(echo $l1 | awk '{print $NF}' | awk '{if (NR!=1) {print substr($2, 1, length($2)-1)}}')
    find . -type f  -name $dep | while read l2
    do
        if [ ! -f "$DESTDIR/etc/apparmor.d/$dep" ]; then
      install_seperate_with_depends $dep
        fi
    done
  done
done
