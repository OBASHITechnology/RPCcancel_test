EXE=client
IMAGE=client-image

all:
	make local

local:
	go build -o $(EXE) client.go

check:
	-go fmt
	-go vet

static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $(EXE) 

docker-image:
	make static
	docker build -t $(IMAGE) .
	docker save -o $(IMAGE) $(IMAGE)

docker-local:
	docker run -d $(IMAGE)

docker:
	make docker-image
	make docker-local
