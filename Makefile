all: build

# "make build" is intended for building directly on the raspberry pi.
# See "make release" below if you wish to cross-compile for the pi.

.PHONY: build
build:
	go build -o build/_output/freezermon ./cmd

# In case you're like me and your IDE doesn't automatically format your go code,
# then "make go-format" is for you.

.PHONY: go-format
go-format:
	go fmt $(shell sh -c "go list ./...")

# There aren't any tests, but if there were...

.PHONY: test
test:
	go test ./...

# "make lint" runs golangci-lint on the codebase. It requires golangci-lint to be
# installed.

.PHONY: lint
lint:
	golangci-lint run ./...

# Use "make install" only when building on the pi itself, as this uses the
# copy that is built in build/_output/

.PHONY: install
install: build/_output/freezermon
	mkdir -p /usr/local/freezermon
	cp build/_output/freezermon /usr/local/freezermon

# Use "make release" to cross-compile for the raspberry pi. Then you can scp
# it over to your pi. It builds to release/linux/arm/

.PHONY: release
release:
	GOOS=linux GOARCH=arm go build -o release/linux/arm/freezermon ./cmd

.PHONY: clean
clean:
	rm -f release/linux/arm/freezermon build/_output/freezermon

