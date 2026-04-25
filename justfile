version := `git describe --tags --always 2>/dev/null || echo "dev"`
ldflags  := "-s -w -X main.version=" + version

default: build

build:
    go build -ldflags '{{ldflags}}' -o json-compact .

run *args:
    go run -ldflags '{{ldflags}}' . {{args}}

clean:
    rm -f json-compact

vet:
    go vet ./...

test:
    go test ./...

tidy:
    go mod tidy
