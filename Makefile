# Makefile

# make run: expects a default named config file in one of the default locations.
# I use '~/.config/cuttle' for testing. Don't forget the certs!

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
CONFIG_DIR=~/.config/cuttle
COMPONENTS_DIR=${ROOT_DIR}/web/components
TARGET_DIR=${ROOT_DIR}/bin
TARGET_NAME=cuttle
TARGET_APP=$(TARGET_DIR)/$(TARGET_NAME)

watch:
	@air
.PHONY: watch

tidy:
	@go mod tidy
.PHONY: tidy

build: tidy gen-templ go-build
.PHONY: build

run: build
	@cd ${TARGET_DIR} && ./${TARGET_NAME}
.PHONY: run

gen-templ:
	@cd ${COMPONENTS_DIR} && templ generate
.PHONY: gen-templ

sync-assets:
	@rsync -ah --delete ./assets ${TARGET_DIR}/
.PHONY: sync-assets

go-build:
	@if [ ! -d ${TARGET_DIR} ]; then mkdir ${TARGET_DIR}; fi && cd cmd/server/ && go build -o ${TARGET_APP}
.PHONY: go-build

test-build:
	@if [ ! -d ${TARGET_DIR} ]; then mkdir ${TARGET_DIR}; fi && cd cmd/server/ && go build -o ${TARGET_APP}
.PHONY: go-build

sshTestServer-start:
	@docker start sshd_test
.PHONY: sshTestServer-start

sshTestServer-update:
	docker stop sshd_test
	docker rm sshd_test
	docker run --name sshd_test -d -p 22:22 cuttle-test-sshd:latest
	@watch "docker ps | grep sshd_test"
.PHONY: sshTestServer-update

sshTestServer-stop:
	@docker stop sshd_test
.PHONY: sshTestServer-stop

setup:
	@cp ./test_helpers/cuttle.yaml ${CONFIG_DIR}/cuttle.yaml
	@cp ./test_helpers/certs/certificate.crt ${CONFIG_DIR}/certs/
	@cp ./test_helpers/certs/privatekey.key ${CONFIG_DIR}/certs/
	@ln -s web/assets assets
.PHONY: setup

clean:
	@go clean
