export GO15VENDOREXPERIMENT=1

all: tools deps fmt build lint

tools:
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint

deps:
	glide install

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./...

build:
	 CGO_ENABLED=0 go build -o "terraform-provider-xenserver-`uname -s`-`uname -m`"
	 ln -sf "terraform-provider-xenserver-`uname -s`-`uname -m`" terraform-provider-xenserver

lint:
	golint

clean:
	rm -f ./terraform-provider-xenserver ./terraform-provider-xenserver-*

.PHONY: tools deps fmt build lint clean
