DOCKER_IMAGE=danielpieper/mrcli-go

.PHONY: install
install:
	docker build -t ${DOCKER_IMAGE}:latest .

.PHONY: sh
sh:
	docker run --rm -it -v mrcli-go_cache:/go/pkg/mod/cache -v $$PWD:/build ${DOCKER_IMAGE}:latest sh

.PHONY: build
build:
	docker run --rm -t -v mrcli-go_cache:/go/pkg/mod/cache -v $$PWD:/build ${DOCKER_IMAGE}:latest

.PHONY: lint
lint:
	docker run --rm -t -v mrcli-go_cache:/go/pkg/mod/cache -v $$PWD:/build ${DOCKER_IMAGE}:latest \
		golangci-lint run --enable-all

.PHONY: test
test:
	docker run --rm -t -v mrcli-go_cache:/go/pkg/mod/cache -v $$PWD:/build ${DOCKER_IMAGE}:latest \
		go test -cover -v

.PHONY: coverage
coverage:
	docker run --rm -t -v mrcli-go_cache:/go/pkg/mod/cache -v $$PWD:/build ${DOCKER_IMAGE}:latest \
		go test -covermode=count -coverprofile=coverage/coverage.out && go tool cover -html=coverage/coverage.out -o coverage/index.html
	URL="file://$$PWD/coverage/index.html"; open $$URL 2>/dev/null || xdg-open $$URL 2>/dev/null

.PHONY: docs
docs:
	URL="http://localhost:8080/pkg/"; open $$URL 2>/dev/null || xdg-open $$URL 2>/dev/null
	godoc -http :8080

.PHONY: clean
clean:
	docker rmi ${DOCKER_IMAGE}:latest

.DEFAULT_GOAL := install
