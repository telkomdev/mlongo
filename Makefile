.PHONY : test build clean format

build:
	go build -o bin/mlongo github.com/telkomdev/mlongo/cmd/mlongo

test:
	go test ./...

format:
	find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" | xargs gofmt -s -d -w

clean:
	rm bin/mlongo