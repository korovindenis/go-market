.PHONY: all

OS = linux
APP_BUILD_NAME = server
PATH_MAIN_GO = ./cmd/app/main.go


all: clean getbuild-server

build-server:
	@echo "  >  Building server"
	@CGO_ENABLED=0 GOOS=$(OS) go build -ldflags "-w" -a -o $(APP_BUILD_NAME) $(PATH_MAIN_GO)

gotest:
	go test `go list ./... | grep -v test` -count 1
	
gotestcover:
	go test -covermode=count -coverprofile=coverage.out $(shell go list ./... | egrep -v '(/test|/test/mock)$$') && go tool cover -func cover.out

get:
	@echo "  >  Checking dependencies"
	@go mod download
	@go install $(PATH_MAIN_GO)

clean:
	@echo "  >  Clearing folder"
	@rm -f ./$(APP_BUILD_NAME)

# MOCKS
.PHONY: mocks
mocks:
	go generate ./...