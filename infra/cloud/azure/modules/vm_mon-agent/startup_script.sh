#!/usr/bin/env bash

set -xeo pipefail
/bin/test -f /etc/systemd/system/mon-agent.service && exit 0

case "$(uname -m)" in
  aarch64*|armv8*|arm64) GOARCH="arm64" ;;
  *)                     GOARCH="amd64" ;;
esac

which jq > /dev/null || yum install -q -y jq

cat << EOF >> /home/azureuser/.ssh/authorized_keys
${file("../../../etc/mon-agent.authorized_keys")}
EOF

mkdir -p /usr/local/bin
version="$(curl -fsS -m 5 -X GET --noproxy '*' -H 'Metadata:true' 'http://169.254.169.254/metadata/instance/compute/tagsList?api-version=2019-06-04' | jq -r '.[] | select(.name | contains("app_version")) | .value')"
storage_token="$(curl -fsS -m 5 -X GET --noproxy '*' -H 'Metadata:true' 'http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fstorage.azure.com%2F' | jq -r .access_token)"
url="${ARTIFACTS_CONTAINER_URL}/$version/mon-agent-azure.$GOARCH"
curl -fsS -m 30 -X GET -H "Authorization: Bearer $storage_token" -H "x-ms-version: 2023-11-03" -o /usr/local/bin/mon-agent "$url"
chmod +x /usr/local/bin/mon-agent
restorecon -v /usr/local/bin/mon-agent || true

vm_id="$(curl -fsS -m 5 -X GET --noproxy '*' -H 'Metadata:true' 'http://169.254.169.254/metadata/instance/compute/vmId?api-version=2017-08-01&format=text')"
vm_name="$(curl -fsS -m 5 -X GET --noproxy '*' -H 'Metadata:true' 'http://169.254.169.254/metadata/instance/compute/name?api-version=2017-08-01&format=text')"
vm_type="$(curl -fsS -m 5 -X GET --noproxy '*' -H 'Metadata:true' 'http://169.254.169.254/metadata/instance/compute/vmSize?api-version=2017-08-01&format=text')"
azone="$(curl -fsS -m 5 -X GET --noproxy '*' -H 'Metadata:true' 'http://169.254.169.254/metadata/instance/compute/physicalZone?api-version=2017-08-01&format=text' || echo '-')"
cat << EOF > /etc/mon-agent.conf
CLOUD_VM_ID=$vm_id
CLOUD_VM_NAME=$vm_name
CLOUD_VM_TYPE=$vm_type
CLOUD_ZONE=$azone
%{ for key, val in ENV ~}
${key}=${val}
%{ endfor ~}
EOF
chmod 600 /etc/mon-agent.conf

vault_token="$(curl -fsS -m 5 -X GET --noproxy '*' -H 'Metadata:true' 'http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fvault.azure.net' | jq -r .access_token)"
curl -fsS -m 5 -X GET -H "Authorization: Bearer $vault_token" "${VAULT_URL}/secrets/mon-agent-env?api-version=7.4" | jq -r .value | jq -r 'to_entries|map("\(.key)=\(.value|tostring|@sh)")|.[]' >> /etc/mon-agent.conf

${file("../../../etc/mon-agent-systemd.sh")}
