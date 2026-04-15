build:
	go build -o bin/cctl ./cmd/cctl

install: build
	cp bin/cctl ~/.local/bin/cctl

clean:
	rm -rf bin/

vet:
	go vet ./...

.PHONY: build install clean vet
