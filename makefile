include .env
export
# ==================================================================================== #
# HELPERS
# ==================================================================================== #

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## help: ask for help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# DEV
# ==================================================================================== #

test:
	go test ./cmd/googleapiclient

test_coverage:
	go test ./... -coverprofile=coverage.out