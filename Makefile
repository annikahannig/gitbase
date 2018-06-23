

deps:
	go get .


test:
	go test -v

all: deps test
	@echo "Fetchted deps and ran tests."	

