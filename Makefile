.PHONY: run
run:
	go run ./cmd/api/main.go

.PHONY: docs
docs:
	swag init -g ./cmd/api/main.go
