all :
	go install -v ./twopt/...
	go install -v ./mesh/...
	

clean :
	go clean -i ./...


bulldogm :
	-rm -rf gcc
	mkdir -p gcc
	cp -R utils twopt mesh gcc/
	cp Makefile.gcc gcc
	cd gcc; make -f Makefile.gcc
