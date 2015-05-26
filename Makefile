.PHONY: test clean all
.PHONY: preprocess brute lookmeup announce grind curiosity montecarlo
all:    preprocess brute lookmeup announce grind curiosity

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

curiosity:
	GOPATH="$(CURDIR)" $(GO) build -o bin/curiosity src/tools/curiosity.go

test:
	GOPATH="$(CURDIR)" $(GO) test hstools -race -v -short

clean:
	- rm -r bin pkg
