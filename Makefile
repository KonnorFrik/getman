.PHONY: build
build:
	go build -o getman getman_cli/main.go

.PHONY: install
install: build
	sudo install -m 755 getman /usr/local/bin/getman
