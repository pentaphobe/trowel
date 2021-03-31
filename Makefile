VERSION=$(shell cat VERSION)
SOURCE_FILES=$(shell find . -name "*.go")

all: build

.PHONY: update_version_tag
update_version_tag:
	@echo Updating ${VERSION} tag
	git push origin :refs/tags/${VERSION}
	git tag -f ${VERSION}
	git push origin --tags


build: $(SOURCE_FILES)
	go build ./...

test: $(SOURCE_FILE)
	go test -v ./...

coverage: OUTFILE:=$(shell mktemp -t "XXXXXX.out")
coverage: $(SOURCE_FILES)
	go test -v -coverprofile=${OUTFILE} ./...
	go tool cover -html=${OUTFILE}
	rm ${OUTFILE}
	
