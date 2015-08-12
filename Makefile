
GOIMPORTS := $(GOPATH)/bin/goimports
GO := go

.PHONY: all build build-debug compile bindata-final bindata-debug run start \
	clean fix bindata-dep check-gcc

FOLDERS := $(shell find * -maxdepth 2 -type d | grep -v web)

all: check-gcc build-debug

build: check-gcc bindata-final compile
build-debug: check-gcc bindata-debug compile


##########################

check-gcc:
	@# check for cygwin: sqlite will not build with cgo and cygwin.
	@sh -c 'gcc -v 2>&1 | grep cygwin >/dev/null || exit 0; echo && echo "Error: Please use mingw, not \
	cygwin."&& echo && exit 1'



compile:
	@echo -------------------------------------------------------------------
	@echo if you encounter include errors please install portaudio1.9-dev,
	@echo libtagc0-dev and libgpod-dev with your package manager
	@echo -------------------------------------------------------------------
	$(GO) get -d
	$(GO) build -v

bindata-final: bindata grunt
	$(GO) run cmd/compile-assets.go

bindata-debug: bindata grunt
	$(GO) run cmd/compile-assets.go -debug

run: build-debug start

test:
	$(GO) get -d ./...
	$(GO) test ./...

start:
	./cfmedias

clean:
	rm -f cfmedias
	rm -rf web/node_modules
	rm -rf web/bower_components
	rm -rf web/assets/vendor
	rm -rf web/bindata.go

fix: goimports
ifeq ($(FOLDERS),)
	@echo "Find is wrong!"
	@exit 1
endif
	$(GOIMPORTS) -l -w $(FOLDERS)
	for f in $$(find . -type f -name "*.go"); do \
		$(GO) fix "$$f"; \
		$(GO) tool vet -composites=false "$$f"; \
	done

.PHONY: grunt
grunt: web/assets/vendor
web/assets/vendor: web/bower_components web/Gruntfile.js
	cd web && node node_modules/grunt-cli/bin/grunt
	touch web/assets/vendor

web/bower_components: web/node_modules web/bower.json
	cd web && node node_modules/bower/bin/bower install
	touch web/bower_components

web/node_modules: web/package.json
	cd web && npm install --quiet
	touch web/node_modules

goimports:
	$(GO) get -v golang.org/x/tools/cmd/goimports

bindata:
	$(GO) get -v github.com/jteeuwen/go-bindata
