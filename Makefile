# set name to tag docker image with
DOCKER_IMAGE := keydrop-agent:latest

# find list of source files with .go extension
SOURCES := $(shell find . -iname "*.go")

# compile keydrop-agent binary for current OS
keydrop-agent: $(SOURCES)
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o $@ .

# build docker image using binary for 64-bit linux
docker-build: docker-keydrop-agent
	docker build --tag $(DOCKER_IMAGE) .

# cross-compile keydrop-agent binary for 64-bit linux docker containers
docker-keydrop-agent: $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o $@ .

# run keydrop locally as a docker container
docker-run: docker-build
	docker run --rm -it --volume /var/run/docker.sock:/var/run/docker.sock $(DOCKER_IMAGE)
