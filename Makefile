DOCKER_IMAGE=danielpieper/mrcli-go

install:
	docker build -t ${DOCKER_IMAGE}:latest .

sh:
	docker run --rm -it -v mrcli-go_cache:/go/pkg/mod/cache -v $$PWD:/build ${DOCKER_IMAGE}:latest sh

build:
	docker run --rm -t -v mrcli-go_cache:/go/pkg/mod/cache -v $$PWD:/build ${DOCKER_IMAGE}:latest

lint:
	docker run --rm -t -v mrcli-go_cache:/go/pkg/mod/cache -v $$PWD:/build ${DOCKER_IMAGE}:latest \
		golangci-lint run --enable-all

test:
	docker run --rm -t -v mrcli-go_cache:/go/pkg/mod/cache -v $$PWD:/build ${DOCKER_IMAGE}:latest \
		go test -cover -v

coverage:
	docker run --rm -t -v mrcli-go_cache:/go/pkg/mod/cache -v $$PWD:/build ${DOCKER_IMAGE}:latest \
		go test -covermode=count -coverprofile=coverage/coverage.out && go tool cover -html=coverage/coverage.out -o coverage/index.html
	URL="file://$$PWD/coverage/index.html"; open $$URL 2>/dev/null || xdg-open $$URL 2>/dev/null

docs:
	URL="http://localhost:8080/pkg/"; open $$URL 2>/dev/null || xdg-open $$URL 2>/dev/null
	godoc -http :8080

clean:
	docker rmi ${DOCKER_IMAGE}:latest

.PHONY: install sh build lint test coverage docs clean
.DEFAULT_GOAL := install
