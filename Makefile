.PHONY: build
build:
	go build -ldflags="-s -w \
		-X ${VERSION_PKG}.GitVersion=${GIT_VERSION} \
		-X ${VERSION_PKG}.GitCommit=${GIT_COMMIT} \
		-X ${VERSION_PKG}.BuildDate=${BUILD_DATE}" \
		-buildmode=pie -a cmd/s3selecgo/main.go