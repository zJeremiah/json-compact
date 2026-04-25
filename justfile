default: build

build:
  go build -o json-compact .

run *args:
  go run . {{args}}

clean:
  rm -f json-compact

vet:
  go vet ./...

test:
  go test ./...

tidy:
  go mod tidy