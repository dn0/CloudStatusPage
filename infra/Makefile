BASE := $(CURDIR)
COMPOSE := $(BASE)/docker/compose.yaml
TF_LOCAL := $(BASE)/core/local
header = "\e[1;34m%-30s\e[m \n\n"
row = "\e[1mmake %-32s\e[m %-50s \n"

ifeq ($(shell command -v podman 2> /dev/null),)
    DOCKER=docker
else
    DOCKER=podman
endif

ifdef APP_VERSION
	VAR_APP_VERSION=-var 'app_version=$(APP_VERSION)'
else
	VAR_APP_VERSION=
endif

.PHONY: all docker-up docker-down upgrade/% cloud/% deploy/% clean/% check/%

all:
	@printf $(header) "*** CloudStatus.Page ***"
	@printf $(row) "upgrade/<tf-folder>" "Run tf init and lock providers"
	@printf "\n"
	@printf $(row) "cloud/gcp/<env>/<cmd>" "Run a command on all mon-agent VMs in a GCP environment"
	@printf $(row) "cloud/aws/<env>/<cmd>"  "Run a command on all mon-agent VMs in an AWS environment"
	@printf $(row) "cloud/azure/<env>/<cmd>" "Run a command on all mon-agent VMs in an Azure environment"
	@printf $(row) "deploy/mon-agent-gcp/<env>" "Deploy APP_VERSION in a GCP environment"
	@printf $(row) "deploy/mon-agent-aws/<env>" "Deploy APP_VERSION in an AWS environment"
	@printf $(row) "deploy/mon-agent-azure/<env>" "Deploy APP_VERSION in an Azure environment"
	@printf $(row) "clean/mon-agent-gcp/<env>" "Destroy all mon-agent VMs in a GCP environment"
	@printf $(row) "clean/mon-agent-aws/<env>" "Destroy all mon-agent VMs in an AWS environment"
	@printf $(row) "clean/mon-agent-azure/<env>" "Destroy all mon-agent VMs in an Azure environment"
	@printf "\n"
	@printf $(row) "deploy/mon-scribe/<env>" "Deploy mon-scribe"
	@printf $(row) "deploy/mon-analyst/<env>" "Deploy mon-analyst"
	@printf $(row) "deploy/mon-web/<env>" "Deploy mon-web"

docker-up:
	$(DOCKER) compose -f $(COMPOSE) up --detach --wait
	sleep 5
	tofu -chdir=$(TF_LOCAL)  apply -auto-approve
	$(DOCKER) compose -f $(COMPOSE) ps
	@echo
	$(DOCKER) compose -f $(COMPOSE) logs --follow

docker-down:
	$(DOCKER) compose -f $(COMPOSE) down

upgrade/%:
	tofu -chdir=$(BASE)/$(@:upgrade/%=%) init --upgrade
	tofu -chdir=$(BASE)/$(@:upgrade/%=%) providers lock -platform=linux_amd64 -platform=darwin_amd64 -platform=linux_arm64 -platform=darwin_arm64

cloud/%:
	bin/cloud-ssh.sh "$(*D)" "$(*F)"

clean/mon-agent-gcp/% clean/mon-agent-azure/%:
	tofu -chdir=$(BASE)/cloud/$(@:clean/mon-agent-%=%) apply -var 'app_version=none' -target=module.vm_mon-agent -destroy -auto-approve

clean/mon-agent-aws/%: TARGETS=$(shell grep '^module.*vm_mon-agent_.*' $(BASE)/cloud/$(@:clean/mon-agent-%=%)/*.tf | cut -d '"' -f 2 | sed 's/^/-target=module./' | tr '\n' ' ')
clean/mon-agent-aws/%:
	tofu -chdir=$(BASE)/cloud/$(@:clean/mon-agent-%=%) apply -var 'app_version=none' $(TARGETS) -destroy -auto-approve

deploy/mon-agent-gcp/% deploy/mon-agent-aws/% deploy/mon-agent-azure/%:
	tofu -chdir=$(BASE)/cloud/$(@:deploy/mon-agent-%=%) apply $(VAR_APP_VERSION) -auto-approve -parallelism=1

deploy/mon-scribe/%: check/app_version
	bin/helm-deploy.sh mon-scribe $(@:deploy/mon-scribe/%=%)

deploy/mon-analyst/%: check/app_version
	bin/helm-deploy.sh mon-analyst $(@:deploy/mon-analyst/%=%)

deploy/mon-web/%: check/app_version
	bin/helm-deploy.sh mon-web $(@:deploy/mon-web/%=%)

check/app_version:
ifndef APP_VERSION
	$(error APP_VERSION is undefined)
endif
