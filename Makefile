APP = proxy-hooker
CONTAINER = proxy-hooker
IMAGE = fguillot/proxy-hooker
TAG = latest

build:
	@ GOOS=linux GOARCH=amd64 go build -o $(APP) src/*.go

image:
	@ docker build -t $(IMAGE):$(TAG) .

push:
	@ docker push $(IMAGE)

run:
	@ docker run -d --name $(CONTAINER) \
	-p 80:80 \
	-v /var/lib/boot2docker:/certs:ro \
	-e DOCKER_HOST=$$DOCKER_HOST \
	$(IMAGE):$(TAG)

logs:
	@ docker logs -f $(CONTAINER)

destroy:
	@ docker rm -f $(CONTAINER)
	@ docker rmi $(IMAGE):$(TAG)

clean:
	@ rm -f $(APP)

all:
	clean build destroy image run
