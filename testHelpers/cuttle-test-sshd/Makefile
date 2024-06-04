# Makefile

build:
	@cat ../testServer_ed25519_no_pass.pub > keys
	@cat ../testServer_ed25519_pass.pub >> keys
	@chmod 600 keys
	docker build -t cuttle-test-sshd .
.PHONY: build

run:
	docker stop sshd_test
	docker rm sshd_test
	docker run --name sshd_test -d -p 22:22 cuttle-test-sshd:latest
	watch "docker ps | grep sshd_test"
.PHONY: run

start:
	docker start sshd_test
.PHONY: start

stop:
	docker stop sshd_test
.PHONY: stop
