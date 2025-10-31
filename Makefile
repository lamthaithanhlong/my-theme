APP_NAME := my-theme

.PHONY: run build tidy fmt clean

run:
	go run ./...

build:
	go build -o $(APP_NAME) ./...

tidy:
	go mod tidy

fmt:
	gofmt -w main.go

clean:
	rm -f $(APP_NAME)

