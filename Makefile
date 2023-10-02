.PHONY: all clean

APP=tiny-url
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./...)
IMPORT_PACKAGES=$(shell go list ./... | awk -F"\/" '{print $$4}')

setup:
	@echo "setting up build environment..."
	@echo "installing lint tool ..."
	@go install golang.org/x/lint/golint@latest
	@echo "installing import manager tool ..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "installing swagger gen tool"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "installing all project dependencies ..."
	@go get ./...

init-swagger:
	@echo "Update swagger docs..."
	@swag init

clean:
	@echo "cleaning test data cache..."
	@rm -rf out/
	@rm -f *.out
	@go clean -testcache

imports:
	@echo "Running imports..."
	@goimports -w -local github.com/isomnath/$(APP) $(IMPORT_PACKAGES)

fmt:
	@echo "running fmt..."
	@go fmt ./...

vet:
	@echo "running vet..."
	@go vet ./...

lint:
	@echo "running lint..."
	@for p in $(ALL_PACKAGES); do \
  		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
  	done

compile:
	@echo "Building executable..."
	@mkdir -p out/
	@go build -o $(APP_EXECUTABLE)

build: imports lint vet compile

copy-config:
	@echo "copying sample config file to environment config file..."
	@cp application.sample.yml application.yml

test: clean
	@echo "running tests..."
	@go test $(ALL_PACKAGES)

test-cover: clean
	@echo "running tests..."
	@mkdir -p out/
	@go test $(ALL_PACKAGES) -coverprofile=coverage.out

test-cover-report: test-cover
	@echo 'Total coverage: '`go tool cover -func coverage.out | tail -1 | awk '{print $$3}'`

test-cover-html: test-cover
	@go tool cover -html=coverage.out


start-application:
	@echo "starting application..."
	@$(APP_EXECUTABLE)

setup-local-infra: clean-stale-infra setup-docker-directory
	@echo "starting all infra for tests..."
	@docker-compose up -d
	@echo "waiting for all components to be available ..."
	@sleep 10

clean-stale-infra:
	@echo "cleaning up stale infra ..."
	@docker stop tiny_url_redis || true
	@docker rm tiny_url_redis || true

setup-docker-directory:
	@echo "setting docker volumes directories ..."
	@rm -rf ~/Docker/redis/data
	@mkdir -p ~/Docker/redis/data