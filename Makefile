test:
	make lint
	make unit
	make integration

build:
	docker build -t pokegen:latest .

integration: build
	go test --tags=integration

unit:
	go test -v ./...

lint:
	docker run -t --rm -v $$(pwd):/app -v ~/.cache/golangci-lint/v1.53.3:/root/.cache -w /app golangci/golangci-lint:v1.53.3 golangci-lint run -v


