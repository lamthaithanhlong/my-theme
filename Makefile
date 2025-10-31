APP_NAME := my-theme
RUN_FLAGS ?=
VOLUME ?=
v ?=

RUN_VOLUME := $(firstword $(filter-out ,$(v) $(VOLUME)))

.PHONY: run build tidy fmt clean

run:
	go run . $(if $(RUN_VOLUME),-v $(RUN_VOLUME)) $(RUN_FLAGS)

build:
	go build -o $(APP_NAME) ./...

tidy:
	go mod tidy

fmt:
	gofmt -w main.go

clean:
	rm -f $(APP_NAME)

