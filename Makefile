all: cover

build:
	mkdir -p .build
	cd .build && \
	GOOS=linux GOARCH=amd64 go build -o main ../cmd/main.go

deps:dependencies
	@go mod tidy

cover:deps
	$(HOME)/go/bin/ginkgo -r --progress --failFast  --randomizeAllSpecs --randomizeSuites --failOnPending --trace --reportFile=./junit.xml -coverpkg=./... -coverprofile=coverage.out -outputdir=./test

report:cover
	@go tool cover -html=./test/coverage.out -o ./test/coverage.html

lint:
	@go vet ./...

clean:
	@rm -fr **/*.{out,xml,html}

clean-cache:
	@go clean -cache
	@go clean -testcache
	@go clean -modcache

dependencies:
	@go get -u github.com/onsi/ginkgo/ginkgo
	@go get -u github.com/onsi/gomega/...

test:deps
	@go test ./...

.PHONY: cover