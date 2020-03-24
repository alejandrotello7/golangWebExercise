GOPATH := ${PWD}:${GOPATH}
export GOPATH

default: build

build:
	go get github.com/talbright/go-zookeeper/zk
	go build -tags netgo -v -o ./bin/grproxy ./src/grproxy
