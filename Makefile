NAME := terraform-provider-unofficial-travis

default: test

test: fmtcheck
	go test ./...

testacc: fmtcheck
	TF_ACC=1 go test -v ./...

build:
	go build -o terraform-provider-utravis_v0.0.0

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w .

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: test testacc build release vet fmt fmtcheck
