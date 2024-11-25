BASE := $(CURDIR)
VERSION ?= $(shell git describe --tags --always --dirty 2> /dev/null)
GOOS ?= linux
GOARCH ?= amd64
GOPROJECT = cspage
LDFLAGS_BASE := -X=$(GOPROJECT)/pkg/config.Version=$(VERSION)
LDFLAGS_PROD := $(LDFLAGS_BASE) -s -w
LDFLAGS_DEV  := $(LDFLAGS_BASE)
BLDFLAGS_PROD := -ldflags "$(LDFLAGS_PROD)" -tags prod
BLDFLAGS_DEV  := -ldflags "$(LDFLAGS_DEV)" -tags dev
TMPDIR := /tmp

header = "\e[1;34m%-30s\e[m \n\n"
row = "\e[1mmake %-32s\e[m %-50s \n"

ifeq ($(shell command -v podman 2> /dev/null),)
    DOCKER=docker
else
    DOCKER=podman
endif

ifdef TAG_LATEST
    UPLOAD_VERSION=latest
else
    UPLOAD_VERSION=$(VERSION)
endif

.PHONY: all
all:
	@printf $(header) "*** CloudStatus.Page ***"
	@printf $(row) "version" "Print current version"
	@printf $(row) "proto" "Generate .pb.go code from .proto files"
	@printf $(row) "templ" "Generate templ.go code from .templ files"
	@printf $(row) "tailwind" "Build style.min.css"
	@printf $(row) "js" "Build .min.js and .js.min.html files"
	@printf $(row) "lint" "Run linters"
	@printf $(row) "test" "Run tests"
	@printf $(row) "test-report" "Generate a test coverage report"
	@printf $(row) "mod/tidy" "Add missing and remove unused Go modules"
	@printf $(row) "mod/update" "Update all Go modules"
	@printf $(row) "bin/<service>" "Go build service binary in bin/"
	@printf $(row) "bin/<service>/<action>" "Run other action on a specific binary"
	@printf $(row) "docker/<service>" "Build docker image for a specific service"
	@printf $(row) "docker/<service>/<action>" "Run other action for a specific docker image"
	@printf $(row) "docker" "Build and upload all docker images"
	@printf $(row) "clean/<service>" "Remove service binary from bin/"
	@printf $(row) "clean" "Remove all service binaries from bin/"
	@printf $(row) "build/<goarch>" "Build all service binaries for specific GOARCH"
	@printf $(row) "build" "Build all service binaries for all GOARCHs"
	@printf $(row) "upload/list" "List all uploaded artifacts"
	@printf $(row) "upload/sync" "Sync artifacts across clouds"
	@printf $(row) "upload" "Upload all service binaries for all GOARCHs"

.PHONY: version
version:
	@echo "$(VERSION)"

.PHONY: clean/%
clean/%:
	rm -fv bin/$(@:clean/%=%).$(GOARCH)

.PHONY: clean
clean:
	rm -fv bin/* tmp/*

.PHONY: proto
proto:
	protoc --go_out=./ --go_opt=paths=source_relative pkg/pb/*.proto

.PHONY: templ
templ:
	TEMPL_EXPERIMENT=rawgo templ fmt pkg
	TEMPL_EXPERIMENT=rawgo templ generate -path pkg

.PHONY: tailwind
tailwind:
	tailwindcss -c .tailwind.config.js -i ./etc/mon-web/input.css -o ./srv/mon-web/css/style.min.css --minify

.PHONY: js
js:
	minify -o pkg/mon/web/templates/main.js.min.html pkg/mon/web/templates/main.js.html
	minify -o srv/mon-web/js/echarts.themes.min.js srv/mon-web/js/echarts.themes.js

.PHONY: lint
lint:
	@which golangci-lint > /dev/null || \
      curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
      sh -s -- -b `go env GOPATH`/bin v1.61.0
	golangci-lint run ./...

.PHONY: test
test:
	go test ./... -v --cover $(BLDFLAGS_PROD)

.PHONY: test-report
test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: mod/tidy
mod/tidy:
	go mod tidy

.PHONY: mod/update
mod/update:
	go get -u all
	make mod/tidy

.PHONY: check/gcp
check/gcp:
	@echo -n GCP:
	@gcloud auth list --filter=status:ACTIVE --format="value(account)"

.PHONY: check/aws
check/aws:
	@echo -n AWS:
	@aws sts get-caller-identity --profile org --output text

.PHONY: check/azure
check/azure:
	@echo -n Azure:
	@az ad signed-in-user show --output tsv

.PHONY: bin/mon-agent-gcp/upload
bin/mon-agent-gcp/upload: check/gcp bin/mon-agent-gcp.amd64 bin/mon-agent-gcp.arm64
	for goarch in amd64 arm64; do \
     for loc in asia europe us; do \
      gcloud artifacts generic upload --quiet \
        --project=cloudstatus-central \
        --repository=mon-agent \
        --location=$$loc \
        --package=$(@:bin/%/upload=%) \
        --version=$(UPLOAD_VERSION) \
        --source=$(BASE)/bin/$(@:bin/%/upload=%).$$goarch ; \
     done \
    done

.PHONY: bin/mon-agent-aws/upload
bin/mon-agent-aws/upload: check/aws bin/mon-agent-aws.amd64 bin/mon-agent-aws.arm64
	for goarch in amd64 arm64; do \
     aws s3 cp --profile org --no-progress $(BASE)/bin/$(@:bin/%/upload=%).$$goarch \
      s3://cloudstatus-artifacts-eu-central-1/mon-agent/$(UPLOAD_VERSION)/$(@:bin/%/upload=%).$$goarch && \
      gcloud artifacts generic upload --quiet \
        --project=cloudstatus-central \
        --repository=mon-agent \
        --location=europe \
        --package=$(@:bin/%/upload=%) \
        --version=$(UPLOAD_VERSION) \
        --source=$(BASE)/bin/$(@:bin/%/upload=%).$$goarch ; \
    done

.PHONY: bin/mon-agent-azure/upload
bin/mon-agent-azure/upload: check/azure bin/mon-agent-azure.amd64 bin/mon-agent-azure.arm64
	for goarch in amd64 arm64; do \
     az storage blob upload --no-progress --overwrite true \
       --auth-mode login \
       --subscription 12345 \
       --account-name cloudstatusartifacts \
       --container-name mon-agent \
       --file $(BASE)/bin/$(@:bin/%/upload=%).$$goarch \
       --name $(UPLOAD_VERSION)/$(@:bin/%/upload=%).$$goarch && \
      gcloud artifacts generic upload --quiet \
        --project=cloudstatus-central \
        --repository=mon-agent \
        --location=europe \
        --package=$(@:bin/%/upload=%) \
        --version=$(UPLOAD_VERSION) \
        --source=$(BASE)/bin/$(@:bin/%/upload=%).$$goarch ; \
    done

bin/%.amd64: GOARCH=amd64
bin/%.arm64: GOARCH=arm64
bin/%:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
    go build $(BLDFLAGS_PROD) \
      -o $(BASE)/bin/$(basename $(@:bin/%=%)).$(GOARCH) \
      "$(GOPROJECT)/cmd/$(basename $(@:bin/%=%))"

.PHONY: docker/login
docker/login:
	gcloud auth print-access-token --quiet | \
      podman login https://europe-docker.pkg.dev --username=oauth2accesstoken --password-stdin

.PHONY: docker/mon-scribe/push
docker/mon-scribe/push: docker/login docker/mon-scribe
	$(DOCKER) push europe-docker.pkg.dev/cloudstatus-central/mon/$(@:docker/mon-%/push=%):$(VERSION)

.PHONY: docker/mon-analyst/push
docker/mon-analyst/push: docker/login docker/mon-analyst
	$(DOCKER) push europe-docker.pkg.dev/cloudstatus-central/mon/$(@:docker/mon-%/push=%):$(VERSION)

.PHONY: docker/mon-web/push
docker/mon-web/push: docker/login docker/mon-web
	$(DOCKER) push europe-docker.pkg.dev/cloudstatus-central/mon/$(@:docker/mon-%/push=%):$(VERSION)

.PHONY: docker/mon-%
docker/mon-%:
	TARGET=$(@:docker/%=%) envsubst '$$TARGET' < Dockerfile | $(DOCKER) build . -f - \
    --platform linux/amd64 \
    --build-arg="VERSION=$(VERSION)" \
    -t europe-docker.pkg.dev/cloudstatus-central/mon/$(@:docker/mon-%=%):$(VERSION) \
    -t mon/$(@:docker/mon-%=%):latest \
    -t mon/$(@:docker/mon-%=%):$(VERSION)
	@echo "*****\n* version=$(VERSION) build complete\n***"

.PHONY: docker
docker: docker/mon-web/push docker/mon-analyst/push docker/mon-scribe/push
	@echo "*****\n* version=$(VERSION) uploaded\n***"

.PHONY: build/amd64
build/amd64: bin/mon-agent-gcp.amd64 bin/mon-agent-aws.amd64 bin/mon-agent-azure.amd64
.PHONY: build/arm64
build/arm64: bin/mon-agent-gcp.arm64 bin/mon-agent-aws.arm64 bin/mon-agent-azure.arm64
.PHONY: build
build: build/arm64 build/amd64
	@echo "*****\n* version=$(VERSION) build complete\n***"

.PHONY: upload/list
upload/list:
	gcloud artifacts files list \
      --project=cloudstatus-central \
      --repository=mon-agent \
      --location=europe \
      --format 'table(name,updateTime,sizeBytes)'

.PHONY: upload/sync
upload/sync: check/gcp check/aws check/azure
	gcloud artifacts files list --project=cloudstatus-central --repository=mon-agent --location=europe --format 'value(name)' | cut -d ':' -f 2- | sed -e 's#:#/#' > "$(TMPDIR)/artifacts"
	aws s3api list-objects --profile org --bucket 'cloudstatus-artifacts-eu-central-1' --prefix 'mon-agent/' --query 'Contents[].{Key: Key}' --output text | cut -d '/' -f 2- | while read -r obj_path; do \
      grep "^$$obj_path$$" "$(TMPDIR)/artifacts" || aws s3 rm --profile org "s3://cloudstatus-artifacts-eu-central-1/mon-agent/$$obj_path" ; \
    done
	az storage blob list --auth-mode login --subscription 12345 --account-name cloudstatusartifacts --container-name mon-agent --output tsv --query '[].[name]' | while read -r obj_path; do \
      grep "^$$obj_path$$" "$(TMPDIR)/artifacts" || az storage blob delete --auth-mode login --subscription 12345 --account-name cloudstatusartifacts --container-name mon-agent --name "$$obj_path" ; \
    done
	@rm -f "$(TMPDIR)/artifacts"

.PHONY: upload
upload: check/gcp check/aws check/azure bin/mon-agent-gcp/upload bin/mon-agent-aws/upload bin/mon-agent-azure/upload
	@echo "*****\n* version=$(VERSION) uploaded\n***"

.PHONY: db/data
db/data:
	@for cloud in gcp aws azure; do \
      echo "COPY $$cloud.mon_config_probe (id, name, description, enabled, config) FROM stdin (FORMAT CSV);" ;\
      cat sql/mon/90-initial_data.$$cloud.config_probe.csv ;\
      echo '\.' ;\
      echo ;\
      echo "COPY $$cloud.mon_config_region (id, name, location, enabled, lat, lon) FROM stdin (FORMAT CSV);" ;\
      cat sql/mon/90-initial_data.$$cloud.config_region.csv ;\
      echo '\.' ;\
      echo ;\
    done
	@echo "SELECT 'All good - Initial DB data were loaded' AS result;"

.PHONY: db/sql
db/sql:
	@echo '\set ON_ERROR_STOP on'
	@echo '\set ECHO queries'
	@cat sql/mon/00-init.sql
	@for cloud in dummy gcp aws azure; do \
      export CLOUD=$$cloud ;\
      envsubst < sql/mon/01-cloud.config_probe.sql ;echo ;\
      envsubst < sql/mon/01-cloud.config_region.sql ;echo ;\
      envsubst < sql/mon/02-cloud.agent.sql ;echo ;\
      envsubst < sql/mon/03-cloud.job.sql ;echo ;\
      envsubst < sql/mon/04-cloud.ping.sql ;echo ;\
      sed -ne 's/^.*ProbeName\s\+=\ "\(\w\+\)".*doc="\(.\+\)".*$$/\1 \2/p' pkg/mon/agent/cloud/$$cloud/*/*_probe.go 2> /dev/null | while read -r probe_name probe_desc; do \
        PROBE_NAME=$$probe_name envsubst < sql/mon/05-cloud.probe.sql ;echo ;\
      done ;\
      envsubst < sql/mon/06-cloud.alert.sql ;echo ;\
      envsubst < sql/mon/07-cloud.incident.sql ;echo ;\
      envsubst < sql/mon/08-cloud.checkpoint.sql ;echo ;\
    done
	@echo
	@echo "SELECT 'All good - DB schema was initialized' AS result;"

.PHONY: db/init
db/init:
	make db/sql | psql --single-transaction "$(DATABASE_URL)"
	make db/data | psql --single-transaction "$(DATABASE_URL)"

####
# Local dev targets
.PHONY: run/%
.PHONY: build/dev
build/dev:
	go build -v $(BLDFLAGS_DEV) -o tmp/ -race ./...
	@echo

.PHONY: templ/watch
templ/watch:
	TEMPL_EXPERIMENT=rawgo templ generate -path pkg -watch -v -proxy http://localhost:8080

.PHONY: tailwind/watch
tailwind/watch:
	tailwindcss -c .tailwind.config.js -i ./etc/mon-web/input.css -o ./srv/mon-web/css/style.min.css --watch

build/mon-web:
	go build -v $(BLDFLAGS_DEV) -o tmp/mon-web cmd/mon-web/main.go
	@echo NEW: tmp/mon-web

run/mon-web: templ tailwind build/mon-web
	LOG_FORMAT=dev air \
      -build.bin tmp/mon-web \
      -build.args_bin '-log-level debug -http-listen-addr 127.0.0.1:8080' \
      -build.pre_cmd '/usr/bin/make templ && /usr/bin/make tailwind && /usr/bin/make js' \
      -build.cmd '/usr/bin/make build/mon-web'

run/mon-analyst: build/dev
	LOG_FORMAT=dev CLOUD=dummy \
    tmp/mon-analyst -log-level debug -http-listen-addr 127.0.0.1:8002 -probe-zscore-threshold 1.0

run/mon-analyst3: build/dev
	LOG_FORMAT=dev CLOUD=gcp DRY_RUN=true LOG_FORMAT=dev \
    tmp/mon-analyst -log-level debug -http-listen-addr 127.0.0.1:8002 \
      -incident-check-interval 0 \
      -ping-check-interval 0 \
      -probe-check-interval 0

run/mon-scribe: build/dev
	LOG_FORMAT=dev \
    tmp/mon-scribe -log-level debug -http-listen-addr 127.0.0.1:8001

run/mon-scribe2: build/dev
	LOG_FORMAT=dev PUBSUB_PROJECT_ID=cloudstatus-t PUBSUB_EMULATOR_HOST=localhost:8085 PUBSUB_PROBE_SUBSCRIPTION="mon-pings-dummy" \
    tmp/mon-scribe -log-level debug -http-listen-addr 127.0.0.1:8001

run/mon-scribe3: build/dev
	LOG_FORMAT=dev PUBSUB_PROJECT_ID=cloudstatus-t PUBSUB_EMULATOR_HOST=localhost:8085 PUBSUB_PROBE_SUBSCRIPTION="mon-pings-dummy,mon-probes-dummy" \
    tmp/mon-scribe -log-level info -http-listen-addr 127.0.0.1:8001

run/mon-agent: build/dev
	CLOUD_REGION=end-of-world \
    SITE_URL="http://localhost:8080/api" \
    tmp/mon-agent-dummy -log-level debug -http-listen-addr 127.0.0.1:8000 -ping-interval 300s \
        -probe-interval-default 0 \
        -probe-long-interval-default 0 \
        -dummy-agent-ping-probe-interval 0 \
        -dummy-internet-ping-probe-interval 0


run/mon-agent2: build/dev
	LOG_FORMAT=dev PUBSUB_PROJECT_ID=cloudstatus-t PUBSUB_EMULATOR_HOST=localhost:8085 CLOUD_REGION=end-of-world \
    tmp/mon-agent-dummy -log-level debug -http-listen-addr 127.0.0.1:8000 -ping-interval 5000ms

run/mon-agent-gcp: build/dev
	LOG_FORMAT=dev CLOUD_REGION=europe-west3 \
    tmp/mon-agent-gcp -log-level debug -http-listen-addr 127.0.0.1:8000 -ping-interval 600s \
                      -probe-interval-default 0 \
                      -probe-long-interval-default 0 \
                      -common-vpc-inter-ping-count 0 \
                      -common-vpc-intra-ping-count 0 \
                      -gcp-storage-bucket-probe-interval 0 \
                      -gcp-storage-object-probe-interval 0 \
                      -gcp-compute-vm-probe-interval 0 \
                      -gcp-compute-vm-spot-probe-interval 0 \
                      -gcp-compute-vm-metadata-probe-interval 0 \
                      -gcp-compute-disk-snapshot-probe-interval 0 \
                      -gcp-pubsub-message-interval 30s \
                      -gcp-vpc-inter-ping-interval 0

run/mon-agent-aws: build/dev
	LOG_FORMAT=dev CLOUD_REGION=eu-central-1 \
    tmp/mon-agent-aws -log-level debug -http-listen-addr 127.0.0.1:8000 -ping-interval 600s \
                      -probe-interval-default 0 \
                      -probe-long-interval-default 0 \
                      -common-vpc-inter-ping-count 0 \
                      -common-vpc-intra-ping-count 0 \
                      -aws-s3-bucket-probe-interval 0 \
                      -aws-s3-object-probe-interval 0 \
                      -aws-ec2-vm-probe-interval 0 \
                      -aws-ec2-vm-spot-probe-interval 0 \
                      -aws-ec2-vm-metadata-probe-interval 0 \
                      -aws-ec2-ebs-snapshot-probe-interval 0 \
                      -aws-sqs-message-interval 0 \
                      -aws-vpc-inter-ping-interval 0

run/mon-agent-azure: build/dev
	LOG_FORMAT=dev CLOUD_REGION=westeurope \
    tmp/mon-agent-azure -log-level debug -http-listen-addr 127.0.0.1:8000 -ping-interval 600s \
                        -probe-interval-default 0 \
                        -probe-long-interval-default 0 \
                        -common-vpc-inter-ping-count 0 \
                        -common-vpc-intra-ping-count 0 \
                        -azure-storage-container-probe-interval 0 \
                        -azure-storage-blob-probe-interval 0 \
                        -azure-compute-vm-probe-interval 0 \
                        -azure-compute-vm-spot-probe-interval 0 \
                        -azure-compute-vm-metadata-probe-interval 0 \
                        -azure-compute-vhd-snapshot-probe-interval 0 \
                        -azure-servicebus-queue-message-interval 0 \
                        -azure-vpc-inter-ping-interval 0
