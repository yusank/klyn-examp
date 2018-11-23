APP?=klyn-examp
RELEASE?=0.0.4
GOOS?=linux
GOARCH?=amd64
PORT?=8081
CONTAINER_IMAGE?=docker.io/yusank/${APP}

clean:
	rm -f ${APP}

build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${APP}

container: build
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

push: container
	docker push $(CONTAINER_IMAGE):$(RELEASE)

run: container
	docker stop $(APP):$(RELEASE) || true && docker rm $(APP):$(RELEASE) || true
	docker run --name $(APP) -p $(PORT):$(PORT) --rm \
	$(APP):$(RELEASE)