SHELL := /bin/bash 

VERSION=`cat version`

default: clean package

manpage:
	pandoc doc/instructor.1.md -s -t man | gzip > doc/instructor.1.gz
	pandoc doc/ins.1.md -s -t man | gzip > doc/ins.1.gz

clean:
	rm instructor-${VERSION} -rf
	rm instructor_${VERSION}* -rf
	rm doc/*.1.gz -rf	

package: manpage
	mkdir instructor-${VERSION}
	cd instructor-*
	cp ./doc/*.1.gz instructor-${VERSION}/
	cp ./debian instructor-${VERSION}/ -r
	cp ./instructor.sh instructor-${VERSION}/
	tar -C instructor-${VERSION}/ -cvaf instructor_${VERSION}.orig.tar.xz ./
	cd instructor-${VERSION} && debuild -us -uc

install:
	mkdir /usr/local/share/instructor -p
	cp instructor.sh /usr/local/share/instructor
	ln -s /usr/local/share/instructor/instructor.sh /usr/bin/instructor -v
	ln -s /usr/local/share/instructor/instructor.sh /usr/bin/ins -v
	cp doc/*.1.gz /usr/share/man/man1

uninstall:
	rm /usr/bin/instructor -f
	rm /usr/bin/ins -f
	rm /usr/share/instructor -rf
