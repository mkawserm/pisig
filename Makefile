

build:
	go build -v -o bin/pisig cmd/pisig/pisig.go

run:
	go run cmd/pisig/pisig.go run server --stderrthreshold INFO --logtostderr=true --v 3

download-module:
	go mod download
