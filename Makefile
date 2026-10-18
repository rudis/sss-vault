all:
	go build
	go test ./...
	go fmt
	go vet

# Code without comments for archival
code.txt: main.go shamir/shamir.go
	for x in $^; do \
	    echo "$$x:"; \
	    awk '/^[[:space:]]*\/\/ (Copyright|SPDX-License-Identifier:) / \
	         || !/^[[:space:]]*\/\// && !/^$$/' "$$x"; \
	    echo; \
	done > $@

static:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 GOAMD64=v1 go build -trimpath -buildvcs=false -o static/sss-vault_linux_amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GOAMD64=v1 go build -trimpath -buildvcs=false -o static/sss-vault_windows_amd64.exe

.PHONY: all static
