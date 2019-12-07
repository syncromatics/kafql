
build:
	docker build -t testing --target test .
	docker build -t final --target final .

test: build
	mkdir -p artifacts
	docker run -v $(PWD)/artifacts:/artifacts -v /var/run/docker.sock:/var/run/docker.sock testing
	cd artifacts && curl -s https://codecov.io/bash | bash

generate-ql:
	go run github.com/99designs/gqlgen -v

ship: build
	docker login --username ${DOCKER_USERNAME} --password ${DOCKER_PASSWORD}
	docker tag final syncromatics/kafql:${VERSION}
	docker push syncromatics/kafql:${VERSION}

