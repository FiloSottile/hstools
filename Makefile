.PHONY: montecarlo preprocess brute lookmeup announce grind test clean all
all: preprocess brute lookmeup announce grind

GO ?= go

montecarlo:
	GOPATH="$(CURDIR)" $(GO) build -o bin/montecarlo src/tools/montecarlo.go

preprocess:
	GOPATH="$(CURDIR)" $(GO) build -o bin/preprocess src/tools/preprocess.go

brute:
	GOPATH="$(CURDIR)" $(GO) build -o bin/brute src/tools/brute.go

lookmeup:
	GOPATH="$(CURDIR)" $(GO) build -o bin/lookmeup src/tools/lookmeup.go

announce:
	GOPATH="$(CURDIR)" $(GO) build -o bin/announce src/tools/announce.go

grind:
	GOPATH="$(CURDIR)" $(GO) build -o bin/grind src/tools/grind.go

test:
	GOPATH="$(CURDIR)" $(GO) test hstools -race -v -short

clean:
	- rm -r bin pkg
