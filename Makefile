# golang-client Makefile
# Follows the interface defined in the Golang CTI proposed
# in https://review.openstack.org/410355

#REPO_VERSION?=$(shell git describe --tags)

GIT_HOST = git.openstack.org

PWD := $(shell pwd)
BASE_DIR := $(shell basename $(PWD))
# Keep an existing GOPATH, make a private one if it is undefined
GOPATH_DEFAULT := $(PWD)/.go
export GOPATH ?= $(GOPATH_DEFAULT)
DEST := $(GOPATH)/src/$(GIT_HOST)/openstack/$(BASE_DIR).git

env:
	@echo "PWD: $(PWD)"
	@echo "BASE_DIR: $(BASE_DIR)"
	@echo "GOPATH: $(GOPATH)"
	@echo "DEST: $(DEST)"

work: $(GOPATH) $(DEST)

$(GOPATH):
	mkdir -p $(GOPATH)

$(DEST): $(GOPATH)
	mkdir -p $(shell dirname $(DEST))
	ln -s $(PWD) $(DEST)

get: work
	cd $(DEST); go get -tags=unit -t ./...

test: get
	cd $(DEST); go test -tags=unit ./...

fmt: work
	cd $(DEST) && go fmt ./...

cover:
	@echo "$@ not yet implemented"

docs:
	@echo "$@ not yet implemented"

relnotes:
	@echo "Reno not yet implemented for this repo"

translation:
	@echo "$@ not yet implemented"

.bindep:
	virtualenv .bindep
	.bindep/bin/pip install bindep

bindep: .bindep
	@.bindep/bin/bindep -b -f bindep.txt || true

install-distro-packages:
	tools/install-distro-packages.sh

clean:
	rm -rf .bindep
	if [ "$(GOPATH)" = "$(GOPATH_DEFAULT)" ]; then \
		rm -rf $(GOPATH); \
	fi

.PHONY: bindep clean
