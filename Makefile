

build:
	go build -v -o bin/pisig cmd/pisig/pisig.go

run:
	go run cmd/pisig/pisig.go run server --stderrthreshold INFO --logtostderr=true --v 3

module-download:
	go mod download

module-update:
	go get github.com/spf13/cobra
	go get github.com/spf13/pflag
	go get github.com/golang/glog
	go get github.com/gobwas/ws
