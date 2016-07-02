APP = proxy-hooker
CONTAINER = proxy-hooker
IMAGE = fguillot/proxy-hooker
TAG = latest

build-linux:
	@ GOOS=linux GOARCH=amd64 go build -o $(APP) src/*.go

build-darwin:
	@ GOOS=darwin GOARCH=amd64 go build -o $(APP) src/*.go

image: build-linux
	@ docker build -t $(IMAGE):$(TAG) .

pull:
	@ docker pull $(IMAGE):$(TAG)

push:
	@ docker push $(IMAGE)

run:
	@ docker run --name $(CONTAINER) \
	-p 10000:80 \
	-v /var/run/docker.sock:/var/run/docker.sock:ro \
	$(IMAGE):$(TAG)

logs:
	@ docker logs -f $(CONTAINER)

destroy:
	@ docker rm -f $(CONTAINER)
	@ docker rmi $(IMAGE):$(TAG)

clean:
	@ rm -f $(APP)

all:
	clean build-linux destroy image run
