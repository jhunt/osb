LDFLAGS := -X main.Version=v$(VERSION)

default: dev
dev:
	go build .
build:
	@echo "Checking that VERSION was defined in the calling environment"
	@test -n "$(VERSION)"
	@echo "OK.  VERSION=$(VERSION)"
	go build -ldflags="$(LDFLAGS)" .

docker:
	@echo "Checking that VERSION was defined in the calling environment"
	@test -n "$(VERSION)"
	@echo "OK.  VERSION=$(VERSION)"
	docker build -t huntprod/osb:$(VERSION) --build-arg VERSION=$(VERSION) .
	docker tag huntprod/osb:$(VERSION) huntprod/osb:latest

push: docker
	docker push huntprod/osb:$(VERSION)
	docker push huntprod/osb:latest

release:
	@echo "Checking that VERSION was defined in the calling environment"
	@test -n "$(VERSION)"
	@echo "OK.  VERSION=$(VERSION)"
	rm -rf artifacts && mkdir artifacts
	GOOS=linux  GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o artifacts/osb-linux-amd64  .
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o artifacts/osb-darwin-amd64 .
