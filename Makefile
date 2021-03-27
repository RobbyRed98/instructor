SHELL := /bin/bash 

VERSION=`cat version`

default: clean package

clean:
	rm instructor-${VERSION} -rf
	rm instructor_${VERSION}* -rf

package:
	mkdir instructor-${VERSION}
	cd instructor-*
	cp ./debian instructor-${VERSION}/ -r
	cp ./instructor.sh instructor-${VERSION}/
	tar -C instructor-${VERSION}/ -cvaf instructor_${VERSION}.orig.tar.xz ./
	cd instructor-${VERSION} && debuild -us -uc

install:
	mkdir /usr/local/share/instructor -p
	cp instructor.sh /usr/local/share/instructor
	ln -s /usr/local/share/instructor/instructor.sh /usr/bin/instructor -v
	ln -s /usr/local/share/instructor/instructor.sh /usr/bin/ins -v

uninstall:
	rm /usr/bin/instructor -f
	rm /usr/bin/ins -f
	rm /usr/share/instructor -rf
