
#TODO add .exe on windows boxen

GOPATH := $(CURDIR)

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
	# portaudio1.9-dev