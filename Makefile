.PHONY: build
build:
	go build -o getman main.go

# .PHONY: install
# install: build
# 	sudo install -m 755 getman /usr/local/bin/getman
