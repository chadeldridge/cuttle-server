# Makefile

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
	@cp ~/cuttle/cuttle.yaml /tmp/cuttle.yaml
	@cp ~/cuttle/certs/certificate.crt /tmp/cuttle_cert.cert
	@cp ~/cuttle/certs/privatekey.key /tmp/cuttle_key.key
.PHONY: setup
