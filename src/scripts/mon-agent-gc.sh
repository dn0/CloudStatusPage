#!/usr/bin/env bash
#shellcheck disable=SC2317

set -euo pipefail

#####

DRY_RUN="${DRY_RUN:-""}"
NAME_FILTER="test"
# brew install coreutils
DATE="$(which gdate || which date)"
OLDER_THAN="30 minutes ago"

#####

C_NC='\033[0m'
C_RED='\033[1;31m'
C_WHITE='\033[1;37m'
C_GRAY='\033[2m'

function log() {
    echo -e "${C_GRAY}[$("$DATE" --rfc-3339=second)] ${*@Q}${C_NC}" >&2
}

function older_than() {
    local epoch="$1"
    local d t line
    while read -r line; do
        [[ -z "$line" ]] && break
        d="$(echo "$line" | awk '{print $2}')"
        t="$("$DATE" --date="$d" +'%s')"
        [[ "$t" -lt "$epoch" ]] && echo "$line"
    done < /dev/stdin
    return 0
}

function purge() {
    if [[ -n "$DRY_RUN" ]]; then
        echo -e "${C_GRAY}[DRY RUN]${C_NC} ${*}"
        return 0
    fi
    "${@}"
}

########
### AWS
export AWS_PROFILE="production"
unset AWS_DEFAULT_REGION
AWS_NAME_FILTER="$NAME_FILTER"
AWS="$(which aws)"
function aws() {
    log aws "${@}"
    "$AWS" --cli-connect-timeout 5 --cli-read-timeout 20 "${@}"
}
function list_aws() {
    local d
    d="$("$DATE" --date="$OLDER_THAN" +'%s')"
    while read -r region; do
      list_aws_region "$region" "$d"
    done <<< "$(aws ec2 describe-regions --region us-east-1 --output text --query 'Regions[*].{Name:RegionName}')"
    aws s3api list-buckets \
        --region us-east-1 \
        --query "Buckets[?(contains(Name,'$AWS_NAME_FILTER'))].{Created:CreationDate,Name:Name}" \
        --output text | sed -e "s/^/BUCKETS\t/" -e "s/$/\tglobal/" | older_than "$d"
}
function list_aws_region() {
    local region="$1"
    local d="$2"
    aws ec2 describe-instances --region "$region" \
        --query "Reservations[*].Instances[*].{Created:LaunchTime,ID:InstanceId,Name:Tags[?(Key==\`Name\`)]|[0].Value,State:State.Name}" \
        --filters "Name=tag:Name,Values=*${AWS_NAME_FILTER}*" "Name=instance-state-name,Values=pending,running,shutting-down,stopping,stopped" \
        --output text | sed -e "s/^/INSTANCES\t/" -e "s/$/\t$region/" | older_than "$d"
    aws ec2 describe-volumes --region "$region" \
        --query "Volumes[*].{Created:CreateTime,ID:VolumeId,Name:Tags[?(Key==\`Name\`)]|[0].Value}" \
        --filters "Name=tag:Name,Values=*${AWS_NAME_FILTER}*" \
        --output text | sed -e "s/^/VOLUMES\t/" -e "s/$/\t$region/" | older_than "$d"
    aws ec2 describe-snapshots --region "$region" \
        --owner-ids self \
        --query "Snapshots[*].{Created:StartTime,ID:SnapshotId,Name:Tags[?(Key==\`Name\`)]|[0].Value}" \
        --filters "Name=tag:Name,Values=*${AWS_NAME_FILTER}*" \
        --output text | sed -e "s/^/SNAPSHOTS\t/" -e "s/$/\t$region/" | older_than "$d"
}
function purge_aws() {
    while read -r kind _created id rest; do
        region=$(echo "$rest" | awk '{ print $NF }')
        case "$kind" in
            "BUCKETS")
                region=$(echo "$id" | sed 's/cloudstatus-probe-[tsp]-\(.*\)-test.*/\1/')
                purge aws s3api delete-bucket --bucket "$id" --region "$region" || true
            ;;
            "INSTANCES")
                purge aws ec2 terminate-instances --instance-ids "$id" --region "$region" || true
            ;;
            "DISKS")
                :
            ;;
            "SNAPSHOTS")
                purge aws ec2 delete-snapshot --snapshot-id "$id" --region "$region" || true
            ;;
        esac
    done
}

########
### Azure
AZURE_SUBSCRIPTION="123"
AZURE_NAME_FILTER="$NAME_FILTER"
AZ="$(which az)"
function az() {
    log az "${@}"
    "$AZ" "${@}"
}
function list_azure() {
    local d storage_accounts
    d="$("$DATE" --date="$OLDER_THAN" +'%s')"
    storage_accounts="$(az resource list \
        --subscription "$AZURE_SUBSCRIPTION" \
        --query "[?(type=='Microsoft.Storage/storageAccounts')].{name:name,group:resourceGroup}" \
        --output tsv)"
    while read -r account group; do
        az storage container list \
            --subscription "$AZURE_SUBSCRIPTION" \
            --auth-mode login \
            --account-name "$account" \
            --query "[?(contains(name,'$AZURE_NAME_FILTER'))].{modified:properties.lastModified,name:name}" \
            --output tsv | sed -e "s/^/Microsoft.Storage\/containers\t/" -e "s/$/\t$account/" | older_than "$d"
    done <<< "$storage_accounts"
    az resource list \
        --subscription "$AZURE_SUBSCRIPTION"  \
        --query "[?(contains(name,'$AZURE_NAME_FILTER')) && ((type=='Microsoft.Compute/virtualMachines') || (type=='Microsoft.Compute/disks') || (type=='Microsoft.Compute/snapshots') || (type=='Microsoft.Network/publicIPAddresses'))].{kind:type,created:createdTime,name:name,group:resourceGroup}" \
        --output tsv | older_than "$d"
}
function purge_azure() {
    while read -r kind _created name group rest; do
        case "$kind" in
            "Microsoft.Storage/containers")
                purge az storage container delete -n "$name" --account-name "$group" --subscription "$AZURE_SUBSCRIPTION" --auth-mode login
            ;;
            "Microsoft.Compute/virtualMachines")
                purge az vm delete -n "$name" -g "$group" --subscription "$AZURE_SUBSCRIPTION" --yes || true
            ;;
            "Microsoft.Compute/disks")
                :
            ;;
            "Microsoft.Compute/snapshots")
                purge az snapshot delete -n "$name" -g "$group" --subscription "$AZURE_SUBSCRIPTION" || true
            ;;
        esac
    done
}

########
### GCP
GCP_PROJECT="cloudstatus-probe-p"
GCP_NAME_FILTER="$NAME_FILTER"
GCLOUD="$(which gcloud)"
function gcloud() {
    log gcloud "${@}"
    "$GCLOUD" "${@}"
}
function list_gcp() {
    local d b f
    d="$("$DATE" --date="$OLDER_THAN" --rfc-3339=second)"
    b='table[no-heading](creation_time.date(tz=UTC):label=CREATED,name,location)'
    f='table[no-heading](kind,creationTimestamp.date(tz=UTC):label=CREATED,name,zone.basename(),status)'
    gcloud storage buckets list \
        --project="$GCP_PROJECT" \
        --filter="creation_time < '${d}' AND name ~ '$GCP_NAME_FILTER'" \
        --format="$b" | sed 's/^/storage#bucket  /'
    gcloud compute instances list \
        --project="$GCP_PROJECT" \
        --filter="creationTimestamp < '${d}' AND name ~ '$GCP_NAME_FILTER'" \
        --format="$f"
    gcloud compute disks list \
        --project="$GCP_PROJECT" \
        --filter="creationTimestamp < '${d}' AND name ~ '$GCP_NAME_FILTER'" \
        --format="$f"
    gcloud compute snapshots list \
        --project="$GCP_PROJECT" \
        --filter="creationTimestamp < '${d}' AND name ~ '$GCP_NAME_FILTER'" \
        --format="$f"
}
function purge_gcp() {
    while read -r kind _created name location rest; do
        case "$kind" in
            "storage#bucket")
                purge gcloud storage buckets delete "gs://$name" --project="$GCP_PROJECT" || true
            ;;
            "compute#instance")
                purge gcloud compute instances delete "$name" --zone="$location" --project="$GCP_PROJECT" || true
            ;;
            "compute#disks")
                :
            ;;
            "compute#snapshot")
                purge gcloud compute snapshots delete "$name" --project="$GCP_PROJECT" || true
            ;;
        esac
    done
}

########
### Main

CLOUD="${1:-""}"
PURGE="${2:-""}"
SAVE_FILE="/tmp/mon-agent-gc-${CLOUD}.txt"

function purge_mode() {
    if [[ "$PURGE" != "purge" ]]; then
        return 1
    fi

    if ! test -r "$SAVE_FILE"; then
        echo "ERROR: Save file is not available" >&2
        exit 2
    fi
    if [[ -z "$(tr -d '[:space:]' < "$SAVE_FILE")" ]]; then
        echo "NOOP: Nothing to purge in $CLOUD"
        exit 0
    fi

    echo
    echo -e "${C_RED}WARNING:${C_NC} ${C_WHITE}Going to delete the following resources in ${CLOUD}:${C_NC}"
    echo
    cat "$SAVE_FILE"
    echo
    [[ -z "$DRY_RUN" ]] && sleep 5

    return 0
}


case "$CLOUD" in
    gcp|aws|azure)
        if purge_mode; then
            eval "purge_$CLOUD" < "$SAVE_FILE"
        else
            eval "list_$CLOUD" | tee "$SAVE_FILE"
        fi
    ;;
    *)
    echo "Usage: $0 {aws|azure|gcp}" >&2
    exit 1
esac
exit 0
