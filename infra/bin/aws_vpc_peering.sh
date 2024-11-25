#!/usr/bin/env bash

set -euo pipefail

DATA_FILE="/tmp/aws_vpc_peering_data.$$.txt"
( for f in vm_mon-agent_*.tf; do
    region="$(basename "$f" .tf | cut -c 14-)"
    vpc_tf="$(tofu state show "module.config_${region}.module.vpc.aws_vpc.this[0]")"
    public_rtb_tf="$(tofu state show "module.config_${region}.module.vpc.aws_route_table.public[0]")"
    vpc_id="$(echo "$vpc_tf" | grep '^\ \ \ \ id\s*=' | cut -d '"' -f 2)"
    vpc_cidr="$(echo "$vpc_tf" | grep '^\ \ \ \ cidr_block\s*=' | cut -d '"' -f 2)"
    vpc_rtb="$(echo "$public_rtb_tf" | grep '^\ \ \ \ id\s*=' | cut -d '"' -f 2)"
    echo "$region $vpc_id $vpc_rtb $vpc_cidr"
done ) | sort > "$DATA_FILE" 

while read -r src_region src_vpc_id src_rtb_id src_cidr; do
    go=""
    while read -r dst_region dst_vpc_id dst_rtb_id dst_cidr; do
        if [[ -z "$go" ]]; then
            if [[ "$src_region" == "$dst_region" ]]; then
                go="true"
            fi
            continue
        fi

        name="peering_${src_region}_${dst_region}"
        vpc_peering_id=$(aws ec2 describe-vpc-peering-connections \
            --region "$src_region" \
            --filters "Name=tag:Name,Values=${name}" \
            --query 'VpcPeeringConnections[*].VpcPeeringConnectionId' \
            --output text)
        if [[ -n "$vpc_peering_id" ]]; then
            echo "${name}: already exists (id:${vpc_peering_id})"
            continue
        fi

        echo "${name}: creating..."
        set -x
        vpc_peering_id=$(aws ec2 create-vpc-peering-connection \
            --region "$src_region" \
            --vpc-id "$src_vpc_id" \
            --peer-region "$dst_region" \
            --peer-vpc-id "$dst_vpc_id" \
            --tag-specifications "ResourceType=vpc-peering-connection,Tags=[{Key=Name,Value=${name}}]" \
            --query VpcPeeringConnection.VpcPeeringConnectionId --output text)
        sleep 10
        aws ec2 accept-vpc-peering-connection \
            --region "$dst_region" \
            --vpc-peering-connection-id "$vpc_peering_id"
        sleep 10
        aws ec2 create-tags \
            --region "$dst_region" \
            --resources "$vpc_peering_id" \
            --tags "Key=Name,Value=${name}"

        aws ec2 modify-vpc-peering-connection-options \
            --region "$src_region" \
            --vpc-peering-connection-id "$vpc_peering_id" \
            --requester-peering-connection-options "AllowDnsResolutionFromRemoteVpc=true"
        aws ec2 modify-vpc-peering-connection-options \
            --region "$dst_region" \
            --vpc-peering-connection-id "$vpc_peering_id" \
            --accepter-peering-connection "AllowDnsResolutionFromRemoteVpc=true"

        aws ec2 create-route \
            --region "$src_region" \
            --route-table-id "$src_rtb_id" \
            --destination-cidr-block "$dst_cidr" \
            --vpc-peering-connection-id "$vpc_peering_id"
        aws ec2 create-route \
            --region "$dst_region" \
            --route-table-id "$dst_rtb_id" \
            --destination-cidr-block "$src_cidr" \
            --vpc-peering-connection-id "$vpc_peering_id"
        set +x

    done < "$DATA_FILE"
done < "$DATA_FILE"

rm -f "$DATA_FILE"
