SHELL := /bin/bash

all: friends metrics tracer

friends:
	cd "$$GOPATH/src/github.com/bobintornado/spfriends"
	docker build -t friends-amd64 -f dockerfile.friends .
	docker system prune -f

metrics:
	cd "$$GOPATH/src/github.com/bobintornado/spfriends"
	docker build -t metrics-amd64 -f dockerfile.metrics .
	docker system prune -f

tracer:
	cd "$$GOPATH/src/github.com/bobintornado/spfriends"
	docker build -t tracer-amd64 -f dockerfile.tracer .
	docker system prune -f

up:
	docker-compose up

down:
	docker-compose down

test:  
	cd "$$GOPATH/src/github.com/bobintornado/spfriends"
	go test ./... -v
