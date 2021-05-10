build:
	CGO_ENABLED=0 go build -o bin/app main.go

docker: build
ifndef TAG
	$(error TAG is not set)
endif
	docker build . -t clazz/users-svc:$(TAG)
	docker build . -t clazz/users-svc:latest

push: docker
ifndef TAG
	$(error TAG is not set)
endif
	docker push clazz/users-svc:$(TAG)
	docker push clazz/users-svc:latest