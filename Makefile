PKG=github.com/shmel1k/exchangego
GOPATH:=$(PWD)/.root:$(GOPATH)
export GOPATH

all: build

clean:
		rm -rf bin/

copy:
	cp .root/src/$(PKG)/exchangego

build:
		go build -i -o bin/exchangego $(PKG)/exchangego
