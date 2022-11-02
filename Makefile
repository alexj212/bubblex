-include .env
export $(shell [ -f ".env" ] && sed 's/=.*//' .env)

export BIN_DIR=./bin
export DIST_DIR=./dist
export PROJ_PATH=github.com/alexj212/bubblex
export DOCKER_APP_NAME=bubblex


export DATE := $(shell date +%Y.%m.%d-%H%M)
export BUILT_ON_IP := $(shell [ $$(uname) = Linux ] && hostname -i || hostname )
export RUNTIME_VER := $(shell go version)
export BUILT_ON_OS=$(shell uname -a)

export LATEST_COMMIT := $(shell git rev-parse HEAD 2> /dev/null)
export COMMIT_CNT := $(shell git rev-list --all 2> /dev/null | wc -l | sed 's/ //g' )
export BRANCH := $(shell git branch  2> /dev/null |grep -v "no branch"| grep \*|cut -d ' ' -f2)
export GIT_REPO := $(shell git config --get remote.origin.url  2> /dev/null)
export COMMIT_DATE := $(shell git log -1 --format=%cd  2> /dev/null)

export BUILT_BY := $(shell whoami  2> /dev/null)
export VERSION=v0.0.1


ifeq ($(BRANCH),)
BRANCH := master
endif

ifeq ($(COMMIT_CNT),)
COMMIT_CNT := 0
endif

export BUILD_NUMBER := ${BRANCH}-${COMMIT_CNT}


export COMPILE_LDFLAGS=-s -X "main.BuildDate=${DATE}" \
                          -X "main.GitRepo=${GIT_REPO}" \
                          -X "main.BuiltBy=${BUILT_BY}" \
                          -X "main.CommitDate=${COMMIT_DATE}" \
                          -X "main.LatestCommit=${LATEST_COMMIT}" \
                          -X "main.Branch=${BRANCH}" \
						  -X "main.Version=${VERSION}"

create_dir:
	@mkdir -p $(BIN_DIR)

check_prereq: create_dir

build_info: check_prereq ## Build the container
	@echo ''
	@echo '---------------------------------------------------------'
	@echo 'BUILT_ON_IP       $(BUILT_ON_IP)'
	@echo 'BUILT_ON_OS       $(BUILT_ON_OS)'
	@echo 'DATE              $(DATE)'
	@echo 'LATEST_COMMIT     $(LATEST_COMMIT)'
	@echo 'BRANCH            $(BRANCH)'
	@echo 'COMMIT_CNT        $(COMMIT_CNT)'
	@echo 'BUILD_NUMBER      $(BUILD_NUMBER)'
	@echo 'COMPILE_LDFLAGS   $(COMPILE_LDFLAGS)'
	@echo 'PATH              $(PATH)'
	@echo '---------------------------------------------------------'
	@echo ''



####################################################################################################################
##
## help for each task - https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
##
####################################################################################################################
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help



####################################################################################################################
##
## Code vetting tools
##
####################################################################################################################


test: ## run tests
	go test -v $(PROJ_PATH)

fmt: ## run fmt on project
	#go fmt $(PROJ_PATH)/...
	gofmt -s -d -w -l .

doc: ## launch godoc on port 6060
	godoc -http=:6060

deps: ## display deps for project
	go list -f '{{ join .Deps  "\n"}}' . |grep "/" | grep -v $(PROJ_PATH)| grep "\." | sort |uniq

lint: ## run lint on the project
	golint ./...

staticcheck: ## run staticcheck on the project
	staticcheck -ignore "$(shell cat .checkignore)" .

vet: ## run go vet on the project
	go vet .

reportcard: fmt ## run goreportcard-cli
	goreportcard-cli -v

tools: ## install dependent tools for code analysis
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install golang.org/x/lint/golint@latest
	go install github.com/gojp/goreportcard/cmd/goreportcard-cli@latest
	go install github.com/goreleaser/goreleaser@latest




####################################################################################################################
##
## Build
##
####################################################################################################################


build: palclient_linux## build binary


build_gui: ## build gui
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	cd ./frontend && yarn run build


build_vroom: ## build vroom
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	cd ./vroom && nuxt build



clean:
	@test ! -e ./dist || rm -rf ./dist


release: release_palclient

release_palclient: palclient_linux palclient_osx palclient_windows
	curl -ualexj:AP12k8ThDp6hvjw6 -T ./dist/osx/paltalk_client          "https://jfrog.theirweb.net/artifactory/misc/osx/paltalk_client"
	curl -ualexj:AP12k8ThDp6hvjw6 -T ./dist/windows/paltalk_client.exe  "https://jfrog.theirweb.net/artifactory/misc/windows/paltalk_client.exe"
	curl -ualexj:AP12k8ThDp6hvjw6 -T ./dist/linux/paltalk_client        "https://jfrog.theirweb.net/artifactory/misc/linux/paltalk_client"


palclient: palclient_linux ## build distribution to ./dist/linux

palclient_linux: build_info ## build distribution to ./dist/linux
	@rm -rf ./dist/linux || true
	@mkdir -p ./dist/linux
	GOOS=linux GOARCH=amd64 go build -o ./dist/linux/paltalk_client -a -ldflags '$(COMPILE_LDFLAGS)' ./app

palclient_osx: build_info ## build distribution to ./dist/linux
	@rm -rf ./dist/osx || true
	@mkdir -p ./dist/osx
	GOOS=darwin GOARCH=amd64 go build -o ./dist/osx/paltalk_client -a -ldflags '$(COMPILE_LDFLAGS)' ./app


palclient_windows: build_info ## build distribution to ./dist/linux
	@rm -rf ./dist/windows || true
	@mkdir -p ./dist/windows
	GOOS=windows GOARCH=amd64 go build -o ./dist/windows/paltalk_client.exe -a -ldflags '$(COMPILE_LDFLAGS)' ./app


all: build## build binaries


publish:
	git add -A
	git commit -m "latest version: $(VERSION)"
	git tag  "$(VERSION)"
	git push origin "$(VERSION)"
	git push

