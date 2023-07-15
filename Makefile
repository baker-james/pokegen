test:
	make integration_test

integration_test:
	docker build -t pokegen:latest .
	go test -v integration_test.go
