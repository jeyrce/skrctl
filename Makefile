version:=latest
ARCH:=amd64

commitId:=$(shell git rev-parse --short HEAD)
buildAt:=$(shell date "+%Y-%m-%d %H:%M:%S")
branch:=$(shell git symbolic-ref --short -q HEAD)

.phony: all
all: x86_64 aarch64

.phony: x86_64
x86_64:
	make binary ARCH=amd64

.phony: aarch64
aarch64:
	make binary ARCH=arm64

.phony: binary
binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -ldflags " \
		-X 'main.version=${version}' \
		-X 'main.commitId=${commitId}' \
		-X 'main.branch=${branch}' \
		-X 'main.buildAt=${buildAt}' \
	" -o _out/skrctl-${ARCH} .

