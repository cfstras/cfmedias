
#TODO add .exe on windows boxen

GOPATH := $(CURDIR)
export GOPATH

all: build

.PHONY: build
build:
	mkdir -p bin
	go build -o bin/cfmedias main

.PHONY: clean
clean:
	rm -rf bin

run: build start

start:
	bin/cfmedias

deps:
	go get github.com/mattn/go-sqlite3
	go get github.com/go-contrib/uuid
	go get code.google.com/p/portaudio-go/portaudio
	go get github.com/coopernurse/gorp
	go get github.com/peterh/liner
	go get code.google.com/p/go.crypto/pbkdf2

	@echo please install portaudio1.9-dev and libtagc0-dev with your package manager
