commitID := $(shell git rev-parse --short HEAD)
buildTime := $(shell date "+%F_%T")
buildName := $(shell whoami)@$(shell hostname)
buildVersion := $(shell cat VERSION)
account := qdm:Qdm12345
branch := $(CI_COMMIT_REF_NAME)

x86_bin := cloudctl-$(buildVersion)-x86-$(branch)-${commitID}
arm_bin := cloudctl-$(buildVersion)-arm-$(branch)-${commitID}

.PHONY: build_x86
build_x86:
	@echo 构建x86可执行文件
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=adm64
	@go build -ldflags \
 	" \
 	-X 'woqutech.com/cloudctl/cmd.BuildVersion=${buildVersion}' \
 	-X 'woqutech.com/cloudctl/cmd.BuildTime=${buildTime}' \
 	-X 'woqutech.com/cloudctl/cmd.BuildName=${buildName}' \
 	-X 'woqutech.com/cloudctl/cmd.CommitID=${commitID}' \
 	" \
 	-o ${x86_bin} cloudctl.go

.PHONY: build_arm
build_arm:
	@echo 构建arm架构可执行文件
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=arm64
	@go build -ldflags \
 	" \
 	-X 'woqutech.com/cloudctl/cmd.BuildVersion=${buildVersion}' \
 	-X 'woqutech.com/cloudctl/cmd.BuildTime=${buildTime}' \
 	-X 'woqutech.com/cloudctl/cmd.BuildName=${buildName}' \
 	-X 'woqutech.com/cloudctl/cmd.CommitID=${commitID}' \
 	" \
 	-o ${arm_bin} cloudctl.go


.PHONY: fmt
fmt:
	go clean
	go fmt

.PHONY: upload_x86
upload_x86:
	echo 上传至大文件存储平台
	curl -u ${account} -T ${x86_bin} http://mirrors.woqutech.com/remote.php/dav/files/qdm/cloudctl/

.PHONY: upload_arm
upload_arm:
	echo 上传至大文件存储平台
	curl -u ${account} -T ${arm_bin} http://mirrors.woqutech.com/remote.php/dav/files/qdm/cloudctl/