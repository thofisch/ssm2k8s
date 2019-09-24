DOCKER_IMAGE_NAME	= thofisch/secrets:$(VERSION)
PACKAGE				:= mystico
OWNER				:= thofisch
REPO				:= ssm2k8s
PROJECT				:= github.com/$(OWNER)/$(REPO)
DATE				?= $(shell date +%FT%T%z)
VERSION				?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(CURDIR)/.version 2> /dev/null || echo v0)
COMMIT				?= $(shell git rev-parse HEAD 2> /dev/null)
CONFIG				= $(shell go list ./internal/config)
LDFLAGS 			= -X $(CONFIG).Version=$(VERSION) -X $(CONFIG).BuildDate=$(DATE) -X $(CONFIG).Commit=$(COMMIT)
ALL_PLATFORMS		:= linux/amd64 darwin/amd64 windows/amd64
CURRENT_OS			= $(shell go env GOOS)
CURRENT_ARCH		= $(shell go env GOARCH)
OS					:= $(if $(GOOS),$(GOOS),$(CURRENT_OS))
ARCH				:= $(if $(GOARCH),$(GOARCH),$(CURRENT_ARCH))
BIN					:= $(CURDIR)/bin
EXECUTABLE			:= $(PACKAGE)-$(OS)-$(ARCH)$(EXT)
OUTBIN				:= $(BIN)/$(EXECUTABLE)
M					= $(shell printf "\033[34;1m▶\033[0m")

export GO111MODULE=on

.PHONY: all
all: build

build-windows-%: EXT = .exe
build-%:
	@$(MAKE) build							\
		--no-print-directory				\
		GOOS=$(firstword $(subst -, , $*))	\
		GOARCH=$(lastword $(subst -, , $*))	\
		EXT=$(EXT)

all-build: $(addprefix build-,$(subst /,-,$(ALL_PLATFORMS))) ## build all defined OS architectures

build: ; $(info $(M) Building binary $(OUTBIN)) @ ## build mystico for current OS architecture
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags '$(LDFLAGS)' -a -installsuffix cgo -o $(OUTBIN) ./cmd/mystico

build-server-local: ; $(info $(M) Building server binary...) @ ## build server executable for local usage
	@go build -ldflags '$(LDFLAGS)' -o $(BIN)/mysticod ./cmd/mysticod

release-windows-%: EXT = .exe
release-%:
	@$(MAKE) release						\
		--no-print-directory				\
		GOOS=$(firstword $(subst -, , $*))	\
		GOARCH=$(lastword $(subst -, , $*))	\
		EXT=$(EXT)

all-release: $(addprefix release-,$(subst /,-,$(ALL_PLATFORMS))) ## publish all defined OS architecture release artifacts

release: ; $(info $(M) Uploading binary $(OUTBIN)) @ ## publish release artifact for current OS architecture
	@github-release -v upload									\
		--security-token $$(git config --global github.token)	\
		--user $(OWNER) 										\
		--repo $(REPO)											\
		--tag $(VERSION)										\
		--name $(EXECUTABLE)									\
		--file $(OUTBIN)

github-release-create: github-release
	@github-release -v release									\
		--security-token $$(git config --global github.token)	\
		--tag $(VERSION)										\
		--user $(OWNER) 										\
		--repo $(REPO)											\

github-release:
	GOOS=$(CURRENT_OS) GOARCH=$(CURRENT_ARCH) go get -u github.com/aktau/github-release

.PHONY: docker
docker: docker-build docker-push ## build and push docker container

.PHONY: docker-build
docker-build: ; $(info $(M) Building docker container $(DOCKER_IMAGE_NAME)) @ ## build docker image
	@docker build --build-arg LDFLAGS="$(LDFLAGS)" -t $(DOCKER_IMAGE_NAME) .

.PHONY: docker-push
docker-push: ; $(info $(M) Pushing docker container $(DOCKER_IMAGE_NAME)) @ ## push docker image
	@docker push $(DOCKER_IMAGE_NAME)

.PHONY: clean
clean: ; $(info $(M) Cleaning...) @ ## clean the build artifacts
	@rm -rf $(BIN)

.PHONY: version
version: ## prints the version (from either environment VERSION, git describe, or .version. default: v0)
	@echo $(VERSION)

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-19s\033[0m %s\n", $$1, $$2}'
