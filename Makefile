.PHONY: test vet lint format-check tidy-check check hooks

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run

format-check:
	@unformatted="$$(gofmt -l .)"; if [ -n "$$unformatted" ]; then echo "gofmt needed:"; echo "$$unformatted"; exit 1; fi

tidy-check:
	go mod tidy -diff

check: format-check tidy-check vet test lint

hooks:
	git config core.hooksPath .githooks
