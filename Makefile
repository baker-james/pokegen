test: build
	make lint
	go test -v ./...

build:
	docker build -t pokegen:latest .

lint:
	docker run -t --rm -v $$(pwd):/app -v ~/.cache/golangci-lint/v1.53.3:/root/.cache -w /app golangci/golangci-lint:v1.53.3 golangci-lint run -v


