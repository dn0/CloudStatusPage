#!/usr/bin/env bash

set -euo pipefail

BASE_DIR="$(cd "$(dirname "$0")/.." || exit 64 ; pwd -P)"
CLOUDENV="${1:-""}"
CMD="${2:-"vms"}"
MAINTF="${BASE_DIR}/cloud/${CLOUDENV}/main.tf"

if [[ -z "$CLOUDENV" ]] || [[ ! -f "$MAINTF" ]]; then
    echo "Usage: $0 <cloud/environment> [cmd]"
    exit 1
fi

case "$CMD" in
    vms|list)  SSH_CMD="" ;;
    uptime)    SSH_CMD="uptime" ;;
    status)    SSH_CMD="sudo systemctl status mon-agent -n0" ;;
    logs|logs) SSH_CMD="sudo journalctl -u mon-agent -n20 --no-pager" ;;
    errors)    SSH_CMD="sudo journalctl -u mon-agent -n100 --no-pager | grep ERROR | cat" ;;
    init-logs-aws) SSH_CMD="sudo journalctl -u cloud-final --no-pager" ;;
    init-errors-aws) SSH_CMD="sudo journalctl -u cloud-final --no-pager | grep fail | cat" ;;
    reboot)    SSH_CMD="sudo reboot" ;;
    poweroff)  SSH_CMD="sudo poweroff" ;;
    *)
        echo "ERROR: unsupported command: '${CMD}'" >&2
        exit 2
    ;;
esac

_slow_cmd() {
    case "$CMD" in
        reboot|poweroff)
            echo "Waiting for 45 seconds..."
            sleep 45
        ;;
        *)
            :
        ;;
    esac
}

CLOUD="${CLOUDENV%%/*}"
ENV="${CLOUDENV##*/}"
USER_ID=""
ACCOUNT_ID="$(grep -m1 '\s\+account_id\s\+=\s\+".*"$' "$MAINTF" | cut -d '"' -f 2)"
[[ -z "$ACCOUNT_ID" ]] && exit 3

SSH_ARGS="-n \
-c chacha20-poly1305@openssh.com \
-o BatchMode=yes \
-o LogLevel=ERROR \
-o UserKnownHostsFile=/dev/null \
-o StrictHostKeyChecking=no \
-o ControlMaster=auto \
-o ControlPath=~/.ssh/master-%r@%h:%p \
-o ControlPersist=3m \
-o ConnectTimeout=10"
C_NC='\033[0m'
C_RED='\033[1;31m'
C_GREEN='\033[1;32m'
C_YELLOW='\033[1;33m'
C_BLUE='\033[1;34m'
C_MAGENTA='\033[0;35m'
C_CYAN='\033[0;36m'
C_WHITE='\033[1;37m'
C_GRAY='\033[2m'
case "$ENV" in 
    production) C_ENV="$C_RED" ;;
    staging)    C_ENV="$C_YELLOW" ;;
    testing)    C_ENV="$C_GREEN" ;;
esac

### Base functions

_print_delim() {
    echo
}

_print_header() {
    local suffix="${1:-""}"
    echo -e "${C_WHITE}**${C_NC} CLOUD=${C_BLUE}${CLOUD}${C_NC} ENV=${C_ENV}${ENV}${C_NC} ${suffix}"
}

_vm_print_header() {
    local region="$1"
    local vm="$2"
    echo -e "${C_WHITE}*${C_NC} REGION=${C_MAGENTA}${region}${C_NC} VM=${C_CYAN}${vm}${C_NC}"
}

#######
### GCP

whoami_print_gcp() {
    USER_ID="$(gcloud auth list --filter='status:ACTIVE' --format='value(account)' || exit 5)"
    _print_header "USER_ID=${C_GRAY}${USER_ID}${C_NC}"
}
_vm_list_gcp() {
    local format="${1:-"table"}"
    gcloud compute instances list --project="$ACCOUNT_ID" --format="${format}[no-heading](NAME,ZONE,INTERNAL_IP,EXTERNAL_IP,STATUS,labels.app_version)"
}
_vm_ssh_gcp() {
    local instance="$1"
    local zone="$2"
    local cmd="$3"
    gcloud compute ssh --project="$ACCOUNT_ID" --zone="$zone" --tunnel-through-iap --ssh-flag="$SSH_ARGS" --command="$cmd" "$instance"
}
vms_print_gcp() {
    _vm_list_gcp
}
vms_ssh_gcp() {
    local cmd="$1"
    while IFS=$',' read -r instance zone ip_private ip_public _unused_; do
        [[ -z "$instance" ]] && continue
        _vm_print_header "$zone" "$instance ($ip_private, $ip_public)"
        _vm_ssh_gcp "$instance" "$zone" "$cmd"
        _print_delim
        _print_header
        _slow_cmd
    done <<< "$(_vm_list_gcp csv | grep RUNNING)"
}

#######
### AWS

whoami_print_aws() {
    USER_ID="$(aws sts get-caller-identity --profile "$ENV" --region us-east-1 --output text --query '[Arn]' || exit 5)"
    _print_header "USER_ID=${C_GRAY}${USER_ID}${C_NC}"
    echo
}
_region_list_aws() {
    aws ec2 describe-regions --profile "$ENV" --region us-east-1 --output text --query 'Regions[*].{Name:RegionName}'
}
_vm_list_aws() {
    # TODO: this will produce a very messy output with multiple VMs
    while read -r region; do
        # shellcheck disable=SC2016
        aws ec2 describe-instances --profile "$ENV" --region "$region" --output text --query 'Reservations[*].Instances[*].[InstanceId, Placement.AvailabilityZone, PrivateIpAddress, PublicIpAddress, State.Name, Tags[?Key==`app_version`]|[0].Value]' &
    done <<< "$(_region_list_aws)"
    wait
}
_vm_ssh_aws() {
    local instance="$1"
    local zone="$2"
    local cmd="$3"
    # shellcheck disable=SC2086
    AWS_DEFAULT_REGION="${zone%?}" ssh ${SSH_ARGS} -l ec2-user "$instance" -- "$cmd"
}
vms_print_aws() {
    _vm_list_aws
}
vms_ssh_aws() {
    local cmd="$1"
    while IFS=$'\t' read -r instance zone ip_private ip_public _unused_; do
        [[ -z "$instance" ]] && continue
        _vm_print_header "$zone" "$instance ($ip_private, $ip_public)"
        _vm_ssh_aws "$instance" "$zone" "$cmd"
        _print_delim
        _print_header
        _slow_cmd
    done <<< "$(_vm_list_aws | grep running)"
}

#########
### Azure

whoami_print_azure() {
    USER_ID="$(az ad signed-in-user show --output tsv --query userPrincipalName || exit 5)"
    _print_header "USER_ID=${C_GRAY}${USER_ID}${C_NC}"
    echo
}
_vm_list_azure() {
    az vm list --subscription "$ACCOUNT_ID" --show-details --output tsv --query '[].[name, location, publicIps, tags.app_version, powerState]'
}
_vm_ssh_azure() {
    local ip="$1"
    local cmd="$2"
    # shellcheck disable=SC2086
    ssh ${SSH_ARGS} -l azureuser "$ip" -- "$cmd"
}
vms_print_azure() {
    _vm_list_azure
}
vms_ssh_azure() {
    local cmd="$1"
    while IFS=$'\t' read -r instance region ip _unused_; do
        [[ -z "$instance" ]] && continue
        _vm_print_header "$region" "$instance ($ip)"
        _vm_ssh_azure "$ip" "$cmd"
        _print_delim
        _print_header
        _slow_cmd
    done <<< "$(_vm_list_azure | grep running)"
}

########
### Main

case "$CLOUD" in
    gcp)
        whoami_print_gcp
        if [[ -z "$SSH_CMD" ]]; then
            vms_print_gcp
        else
            vms_ssh_gcp "$SSH_CMD"
        fi
    ;;
    aws)
        whoami_print_aws
        if [[ -z "$SSH_CMD" ]]; then
            vms_print_aws
        else
            vms_ssh_aws "$SSH_CMD"
        fi
    ;;
    azure)
        whoami_print_azure
        if [[ -z "$SSH_CMD" ]]; then
            vms_print_azure
        else
            vms_ssh_azure "$SSH_CMD"
        fi
    ;;
esac
exit 0
