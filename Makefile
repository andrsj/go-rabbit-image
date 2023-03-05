BINARY_NAME=server

rabbit:
	docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management

rabbit-stop:
	docker stop some-rabbit
	docker rm some-rabbit

.PHONY: run
run:
	go run cmd/main.go

.PHONY: build
build:
	go build -o build/${BINARY_NAME} cmd/main.go

.PHONY: brun
brun: build
	./build/${BINARY_NAME}

.PHONY: clean
clean:
	go clean
	rm build/${BINARY_NAME}