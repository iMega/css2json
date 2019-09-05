REPO = github.com/imega/css2json
CWD = /go/src/$(REPO)

test:
	docker run --rm -v $(CURDIR):$(CWD) -w $(CWD) golang:alpine \
		sh -c "go list ./... | xargs go test"
