all: clean test build

build:
	@echo Building "./out/${APP}"...
	@mkdir -p ./out
	@packr2
	@GO111MODULE=on go build -o "./out/${APP}"
	@packr2 clean

test:
	GO111MODULE=on go test -count 1 -cover -v ./...

fmt:
	gofmt -l -s -w $(SOURCE_DIRS)

imports:
	go get -u golang.org/x/tools/cmd/goimports
	goimports -l -w -v $(SOURCE_DIRS)

cyclo:
	go get -u github.com/fzipp/gocyclo
	gocyclo -over 9 $(SOURCE_DIRS)

vet:
	GO111MODULE=on go vet ./...

lint:
	go get -u golang.org/x/lint/golint
	golint -set_exit_status ./...

