default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 -run .*Queue.* ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m  ./...

.PHONY: fmt lint test testacc build install generate

build-generator:
	go build -C ./tools/generator/ -o ./enablelagen ./main.go


generate-terraform: build-generator
	go generate ./internal/generate/