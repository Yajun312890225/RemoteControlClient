.PHONY: start build

all: run

build:
	GOOS=linux GOARCH=arm go build main.go
run: 
	go run main.go