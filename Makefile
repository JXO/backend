precommit: fmt test

test:
	@go test -race $$(go list ./... | grep -v vendor)
fmt:
	@go fmt $$(go list ./... | grep -v vendor)
fast_test:
	@go test $$(go list ./... | grep -v vendor)

check_fmt:
ifneq ($(shell gofmt -l ./ | grep -v vendor | grep -v testdata),)
	$(error code not fmted, run make fmt. $(shell gofmt -l ./ | grep -v vendor | grep -v testdata))
endif
