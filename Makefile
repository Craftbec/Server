
.PHONY: build
build:
	@docker build -t serv . 


.PHONY: run
run: build
	@docker run --name=serv1 -it -p 8080:8080 serv


.PHONY: clean
clean:
	@docker rm serv1
	@docker rmi serv

.PHONY: test
test:
	go test -v ./internal/storage
