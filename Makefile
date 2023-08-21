# Set defaults
export SPANNER_INSTANCE=acquiring-instance
export SPANNER_DATABASE=authorizations
export GCLOUD_PROJECT=cc-acquiring-development
export SPANNER_EMULATOR_HOST=localhost:10010

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST}

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'create migrations files for ${name}'
	migrate create -seq -ext=.sql -dir=./database/migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	@echo "spanner://projects/${GCLOUD_PROJECT}/instances/${SPANNER_INSTANCE}/databases/${SPANNER_DATABASE}?x-clean-statements=true"
	migrate -path ./database/migrations -database spanner://projects/${GCLOUD_PROJECT}/instances/${SPANNER_INSTANCE}/databases/${SPANNER_DATABASE}?x-clean-statements=true up

## db/migrations/down version=$1: apply down database migrations to version specified
.PHONY: db/migrations/down
db/migrations/down:
	@echo 'Running up migrations to version ${version}...'
	@echo "spanner://projects/${GCLOUD_PROJECT}/instances/${SPANNER_INSTANCE}/databases/${SPANNER_DATABASE}?x-clean-statements=true down ${version}"
	migrate -path ./database/migrations -database spanner://projects/${GCLOUD_PROJECT}/instances/${SPANNER_INSTANCE}/databases/${SPANNER_DATABASE}?x-clean-statements=true down ${version}

## development/db/migrations/force version=$1: force DB into specified version
.PHONY: development/db/migrations/force
development/db/migrations/force:
	@echo 'Forcing DB into version ${version}...'
	@echo "spanner://projects/cc-acquiring-development/instances/acquiring-instance/databases/authorizations?x-clean-statements=true force ${version}"
	migrate -path ./database/migrations -database spanner://projects/cc-acquiring-development/instances/acquiring-instance/databases/authorizations?x-clean-statements=true force ${version}

## acceptance/db/migrations/force version=$1: force DB into specified version
.PHONY: acceptance/db/migrations/force
acceptance/db/migrations/force:
	@echo 'Forcing DB into version ${version}...'
	@echo "spanner://projects/cc-acquiring-acceptance/instances/acquiring-instance/databases/authorizations\?x-clean-statements=true force ${version}"
	migrate -path ./database/migrations -database spanner://projects/cc-acquiring-acceptance/instances/acquiring-instance/databases/authorizations?x-clean-statements=true force ${version}

## db/seed/up: apply test data from seed.sql file
.PHONY: db/seed/up
db/seed/up:
	@echo 'Seeding the database...'
	go run database/seed/main.go seed

## db/seed/down: dump test data
.PHONY: db/seed/down
db/seed/down:
	@echo 'Dumping test data...'
	go run database/fixture/seed.go dump

# ==================================================================================== #
# GCLOUD
# ==================================================================================== #

## gcloud/config/spanner/emulator/create: create config for Spanner emulator
.PHONY: gcloud/config/spanner/emulator/create
gcloud/config/spanner/emulator/create:
	@echo 'Create Spanner emulator config...'
	gcloud config configurations create ${GCLOUD_PROJECT}
	gcloud config set auth/disable_credentials true
	gcloud config set project ${GCLOUD_PROJECT}
	gcloud config set api_endpoint_overrides/spanner http://localhost:10020/

## gcloud/config/activate name=$1: activate a gcloud config [default|emulator]
.PHONY: gcloud/config/activate
gcloud/config/activate:
	@echo 'activate gcloud config ${name}...'
	gcloud config configurations activate name ${NAME}

# ==================================================================================== #
# SPANNER
# ==================================================================================== #

## spanner/emulator/start: start the Spanner emulator
.PHONY: spanner/emulator/start
spanner/emulator/start:
	@echo 'Start Spanner emulator...'
	gcloud emulators spanner start &

## spanner/emulator/instance/create: create instance for Spanner emulator
.PHONY: spanner/emulator/instance/create
spanner/emulator/instance/create:
	@echo 'Create spanner instance'
	gcloud spanner instances create ${SPANNER_INSTANCE} --config=emulator-config --description="Creditcard acquiring" --nodes=1

## spanner/emulator/database/create: create database for Spanner emulator
.PHONY: spanner/emulator/database/create
spanner/emulator/database/create:
	@echo 'Create spanner database'
	gcloud spanner databases create ${SPANNER_DATABASE} --instance=${SPANNER_INSTANCE}

## spanner/emulator/run: run the emulator, create an instance and a DB
.PHONY: spanner/emulator/run
spanner/emulator/run: spanner/emulator/start spanner/emulator/instance/create spanner/emulator/database/create

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## Tidy: tidy dependencies
.PHONY: tidy
tidy:
	@echo 'Tidy module dependencies...'
	go mod tidy
	go mod verify

## generate: generate all mocks
.PHONY: generate
generate:
	@echo 'Generate mocks'
	go generate ./...

## fmt: format code with go fmt
.PHONY: fmt
fmt:
	@echo 'Formatting code...'
	go fmt ./...

## GCI: set import order of golang packages
.PHONY: gci
gci:
	@echo 'Setting order golang packages...'
	gci -w .

## lint: lint code with custom linters
.PHONY: lint
lint:
	@echo 'Linting code...'
	golangci-lint run

## vuln: execute package vulnerabilities check
.PHONY: vuln
vuln:
	@echo 'Check code for known vulnerabilities...'
	govulncheck ./...

## test: run all tests with go test
.PHONY: test
test:
	@echo 'Running tests'
	go test -race -vet=off ./...

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: tidy generate fmt lint vuln test

# ==================================================================================== #
# BOOTSTRAP SERVICE LOCALLY
# ==================================================================================== #

## setup/api: Bootstrap Spanner emulator and start the authorization API.
.PHONY: setup/api
setup/api: spanner/emulator/start spanner/emulator/instance/create spanner/emulator/database/create db/migrations/up db/seed/up run/api