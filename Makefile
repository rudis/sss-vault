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

.PHONY: all
