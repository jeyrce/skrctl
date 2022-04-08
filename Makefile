version:=latest
commitId:=$(shell git rev-parse --short HEAD)
buildAt:=$(shell date "+%Y-%m-%d %H:%M:%S")
branch:=$(shell git symbolic-ref --short -q HEAD)

.phony: all
all: binary

.phony: binary
binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags " \
		-X 'main.version=${version}' \
		-X 'main.commitId=${commitId}' \
		-X 'main.branch=${branch}' \
		-X 'main.buildAt=${buildAt}' \
	" -o skrctl .
