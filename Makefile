-include .env
export $(shell [ -f ".env" ] && sed 's/=.*//' .env)

export COMMIT_CNT := $(shell git rev-list --count HEAD 2> /dev/null)
export DATE := $(shell date +%Y.%m.%d-%H%M)
export LATEST_COMMIT := $(shell git rev-parse HEAD 2> /dev/null)
export BRANCH := $(shell git branch  2> /dev/null |grep -v "no branch"| grep \*|cut -d ' ' -f2)
export GIT_REPO := $(shell git config --get remote.origin.url  2> /dev/null)
export COMMIT_DATE := $(shell git log -1 --format=%cd  2> /dev/null)
export VERSION_FILE   := version.txt
export TAG     := $(shell [ -f "$(VERSION_FILE)" ] && cat "$(VERSION_FILE)" || echo '0.0.4')
export VERMAJMIN      := $(subst ., ,$(TAG))
export VERSION        := $(word 1,$(VERMAJMIN))
export MAJOR          := $(word 2,$(VERMAJMIN))
export MINOR          := $(word 3,$(VERMAJMIN))
export NEW_MINOR      := $(shell expr "$(MINOR)" + 1)
export NEW_TAG := $(VERSION).$(MAJOR).$(NEW_MINOR)

ifeq ($(BRANCH),)
BRANCH := master
endif


export COMPILE_LDFLAGS=-s -X "main.BuildDate=${DATE}" \
                          -X "main.LatestCommit=${LATEST_COMMIT}" \
						  -X "main.Version=${NEW_TAG}"\
						  -X "main.GitRepo=${GIT_REPO}" \
                          -X "main.GitBranch=${BRANCH}"


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


build: ## build binary
	go build ./...

clean:
	@test ! -e ./dist || rm -rf ./dist

all: build## build binaries


publish:
	git add -A
	git commit -m "latest version: $(VERSION)"
	git tag  "$(VERSION)"
	git push origin "$(VERSION)"
	git push





upgrade:
	go get -u ./...
	go mod tidy