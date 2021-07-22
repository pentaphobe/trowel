VERSION=$(shell cat VERSION)
SOURCE_FILES=$(shell find . -name "*.go" -not -path "./example/*") 

all: build example_usage

.PHONY: update_version_tag
update_version_tag: DEFAULT_BRANCH:=$(shell git remote show origin | awk '/HEAD branch/ {print $$NF}')
update_version_tag: CURRENT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD)
update_version_tag: 		
	@if ! [ "$(CURRENT_BRANCH)" == "$(DEFAULT_BRANCH)" ]; then \
		echo Not on default branch; \
		false; \
	fi
	@echo Updating ${VERSION} tag
	git push origin :refs/tags/${VERSION}
	git tag -f ${VERSION}
	git push origin --tags

build: $(SOURCE_FILES)
	go build .

test: $(SOURCE_FILE) example_usage
	go test -v ./...

coverage: OUTFILE:=$(shell mktemp -t "XXXXXX.out")
coverage: $(SOURCE_FILES)
	go test -v -coverprofile=${OUTFILE} ./...
	go tool cover -html=${OUTFILE}
	rm ${OUTFILE}
	
.PHONY: example_usage
example_usage: OUTFILE:=$(shell mktemp -t "XXXXXX.out")
example_usage: $(shell find ./example -name "*.go")
	@go build -o ${OUTFILE} ./example/... || echo "Failed to build example"
	@# attempt to run the file (if it fails then we're unhappy)
	@${OUTFILE} || echo "Failed to run example"