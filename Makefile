.PHONY: montecarlo preprocess brute test clean all
all: montecarlo preprocess brute

GO     ?= go

montecarlo:
	GOPATH="$(CURDIR)" $(GO) build -o bin/montecarlo src/tools/montecarlo.go

preprocess:
	GOPATH="$(CURDIR)" $(GO) build -o bin/preprocess src/tools/preprocess.go

brute:
	GOPATH="$(CURDIR)" $(GO) build -o bin/brute src/tools/brute.go

test:
	GOPATH="$(CURDIR)" $(GO) test hstools -race -v -short

clean:
	- rm -r bin pkg
