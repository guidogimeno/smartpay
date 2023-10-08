build:
	go build -o bin/smartpay cmd/main.go

run: build
	./bin/smartpay

test: go test -v ./...

clean:
	rm -rf bin
