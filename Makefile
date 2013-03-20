all :
	go install -v ./particle
	go install -v ./twopt
	

clean :
	go clean -i ./...