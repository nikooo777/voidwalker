BINARY=voidwalker

DIR = $(shell cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)
BIN_DIR = ${DIR}/bin
IMPORT_PATH = github.com/lbryio/voidwalker
GOARCH = amd64

VERSION = $(shell git --git-dir=${DIR}/.git describe --dirty --always --long --abbrev=7)
LDFLAGS = -ldflags "-X ${IMPORT_PATH}/meta.Version=${VERSION} -X ${IMPORT_PATH}/meta.Time=$(shell date +%s) -w"


.PHONY: build clean test lint dev
.DEFAULT_GOAL: build


build:
	mkdir -p ${BIN_DIR} && CGO_ENABLED=1 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -asmflags -trimpath=${DIR} -o ${BIN_DIR}/${BINARY} main.go

clean:
	if [ -f ${BIN_DIR}/${BINARY} ]; then rm ${BIN_DIR}/${BINARY}; fi

test:
	go test ./... -v -cover

lint:
	go get github.com/alecthomas/gometalinter && gometalinter --install && gometalinter ./...

dev:
	reflex --decoration=none --start-service=true --regex='main.go' go run .