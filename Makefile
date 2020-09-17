

SCRIPTS_PATH ?= build

# Install software dependencies
INSTALL_DEPENDENCIES ?= ${SCRIPTS_PATH}/install-dependencies.sh
# The command to run to execute unit tests
UNIT_TEST_COMMAND ?= ${SCRIPTS_PATH}/run-unit-tests.sh

export GOPACKAGES ?= ./pkg/...

.PHONY: deps
deps:
	$(INSTALL_DEPENDENCIES)

.PHONY: test
## Runs go unit tests
test:
	@if ! which kubebuilder > /dev/null; then \
	  echo "Please install kubebuilder, run 'make deps' and 'export PATH=$PATH:/usr/local/kubebuilder/bin'"; \
	else \
		$(UNIT_TEST_COMMAND); \
	fi

.PHONY: go/gosec-install
## Installs latest release of Gosec
go/gosec-install:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(GOPATH)/bin


.PHONY: go/gosec
## Runs gosec in quiet mode (meaning output only if issues found). Any findings will be printed to stdout.
go/gosec: go/gosec-install
	gosec --quiet ./...
	
SONAR_GO_TEST_ARGS ?= ./...


.PHONY: sonar/go
## Run SonarCloud analysis for Go on Travis CI. This will not be run during local development.
sonar/go: go/gosec-install
	@echo "-> Starting sonar/go"
	@echo "--> Starting go test"
	go test -coverprofile=coverage.out -json ${SONAR_GO_TEST_ARGS} > report.json
	@echo "---> go test report.json"
	@grep -v '"Action":"output"' report.json
	@echo "--> Running gosec"
	gosec -fmt sonarqube -out gosec.json -no-fail ./...
	@echo "---> gosec gosec.json"
	@cat gosec.json
	@echo "--> Running sonar-scanner"
	unset SONARQUBE_SCANNER_PARAMS
	sonar-scanner --debug


# This expects that your code uses Jest to execute tests.
# Add this field to your jest.config.js file to generate the report:
#     testResultProcessor: 'jest-sonar-reporter',
# It must be run before make component/test/unit.
.PHONY: sonar/js/jest-init
## Install npm module to make Sonar test reports in Jest on Travis. This will not be run during local development.
sonar/js/jest-init:
	npm install -D jest-sonar-reporter


# Test reports and code coverage must be generated before running the scanner.
# It must be run after make component/test/unit.
.PHONY: sonar/js
## Runs the SonarCloud analysis for JavaScript on Travis. This will not be run during local development.
sonar/js:
	unset SONARQUBE_SCANNER_PARAMS
	sonar-scanner --debug

.PHONY: go-bindata
go-bindata:
	@if which go-bindata > /dev/null; then \
		echo "##### Updating go-bindata..."; \
		cd $(mktemp -d) && GOSUMDB=off go get -u github.com/go-bindata/go-bindata/...; \
	fi
	@go-bindata --version
	go-bindata -nometadata -pkg bindata -o examples/applier/bindata/bindata_generated.go -prefix examples/applier/resources/yamlfilereader  examples/applier/resources/yamlfilereader/...

.PHONY: examples
examples:
	@mkdir -p examples/bin
	go build -o examples/bin/apply-some-yaml examples/applier/apply-some-yaml/main.go
	go build -o examples/bin/apply-yaml-in-dir examples/applier/apply-yaml-in-dir/main.go
	go build -o examples/bin/render-list-yaml examples/applier/render-list-yaml/main.go
	go build -o examples/bin/render-yaml-in-dir examples/applier/render-yaml-in-dir/main.go
	
.PHONY: build
build:
	go build -o bin/applier cmd/applier/main.go

.PHONY: functional-test
functional-test:
	ginkgo -tags functional -v --slowSpecThreshold=30 test/functional -test.coverprofile test/functional/coverage/functional.out -- -v=1

.PHONY: functional-test-full
functional-test-full: 
	@build/run-functional-tests.sh

.PHONY: kind-cluster-setup
kind-cluster-setup: 
	@echo "No setup to do"
