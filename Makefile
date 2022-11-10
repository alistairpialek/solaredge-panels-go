.PHONY: run
run:
	go run .

.PHONY: test
test:
	docker run --rm -v ${PWD}:/go/src -w /go/src \
	golang:1.19 \
	go test -v -cover ./solaredge
