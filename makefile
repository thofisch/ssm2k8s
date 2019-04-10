DOCKER_IMAGE_NAME = thofisch/secrets:$(IMAGE_TAG)

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
