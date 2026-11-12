.PHONY: build
build:
	go build -o getman ./cmd/getman

# .PHONY: install
# install: build
# 	sudo install -m 755 getman /usr/local/bin/getman
