version:=latest
commitId:=$(shell git rev-parse --short HEAD)
buildAt:=$(shell date "+%Y-%m-%d %H:%M:%S")
branch:=$(shell git symbolic-ref --short -q HEAD)

.phony: all
all: binary

.phony: binary
binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags " \
		-X 'github.com/skrbox/skrctl/main.version=${version}' \
		-X 'github.com/skrbox/skrctl/main.commitId=${commitId}' \
		-X 'github.com/skrbox/skrctl/main.branch=${branch}' \
		-X 'github.com/skrbox/skrctl/main.buildAt=${buildAt}' \
	" -o skrctl .
