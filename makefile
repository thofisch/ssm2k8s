DOCKER_IMAGE_NAME	:= thofisch/secrets:$(IMAGE_TAG)
BIN					:= myapp
ALL_PLATFORMS		:= linux/amd64 darwin/amd64 windows/amd64
CURRENT_OS			= $(shell go env GOOS)
CURRENT_ARCH		= $(shell go env GOARCH)
OS					:= $(if $(GOOS),$(GOOS),$(CURRENT_OS))
ARCH				:= $(if $(GOARCH),$(GOARCH),$(CURRENT_ARCH))
BIN_NAME			:= $(BIN)-$(OS)-$(ARCH)$(EXT)
OUTBIN				:= bin/$(BIN_NAME)

all: build

build-windows-%: EXT = .exe
build-%:
	@$(MAKE) build							\
		--no-print-directory				\
		GOOS=$(firstword $(subst -, , $*))	\
		GOARCH=$(lastword $(subst -, , $*))	\
		EXT=$(EXT)

all-build: $(addprefix build-,$(subst /,-,$(ALL_PLATFORMS)))

build:
	@echo "Building: OS=$(OS), ARCH=$(ARCH), OUTBIN=$(OUTBIN)"
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -a -v -installsuffix cgo -o $(OUTBIN) ./cmd/mystico

release-windows-%: EXT = .exe
release-%:
	@$(MAKE) release						\
		--no-print-directory				\
		GOOS=$(firstword $(subst -, , $*))	\
		GOARCH=$(lastword $(subst -, , $*))	\
		EXT=$(EXT)

all-release: $(addprefix release-,$(subst /,-,$(ALL_PLATFORMS)))

#     -t, --tag            Git tag to upload to (*)
#     -n, --name           Name of the file (*)
#     -l, --label          Label (description) of the file
#     -f, --file           File to upload (use - for stdin) (*)
#     -R, --replace        Replace asset with same name if it already exists (WARNING: not atomic, failure to upload will remove the original asset too)

release:
	@echo "Uploading binary: $(OUTBIN)"
	@github-release -v upload	\
		-s eed63942ff8087e955829bd94a50d15e6623073e \
		--user thofisch 	\
		--repo ssm2k8s		\
		--tag v0.1.0		\
		--name $(BIN_NAME)	\
		--file $(OUTBIN)

github-release-create: github-release
	@github-release -v release	\
		-s eed63942ff8087e955829bd94a50d15e6623073e \
		--tag v0.1.0		\
		--user thofisch		\
		--repo ssm2k8s		\

github-release:
	GOOS=$(CURRENT_OS) GOARCH=$(CURRENT_ARCH) go get -u github.com/aktau/github-release

docker: docker.build docker.push

docker.build: required
	docker build -t $(DOCKER_IMAGE_NAME) .

docker.push: required
	docker push $(DOCKER_IMAGE_NAME)

.PHONY: required
required:
ifndef IMAGE_TAG
	$(error IMAGE_TAG not set)
endif
