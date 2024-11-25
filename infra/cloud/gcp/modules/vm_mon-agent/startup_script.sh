#!/usr/bin/env bash

set -xeo pipefail
/bin/test -f /etc/systemd/system/mon-agent.service && exit 0

case "$(uname -m)" in
  aarch64*|armv8*|arm64) GOARCH="arm64" ;;
  *)                     GOARCH="amd64" ;;
esac

which jq > /dev/null || yum install -q -y jq

mkdir -p /usr/local/bin
token="$(curl -fsS -m 5 -X GET -H 'Metadata-Flavor: Google' 'http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token' | jq -r .access_token)"
version="$(curl -fsS -m 5 -X GET -H 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/attributes/app_version)"
url="https://artifactregistry.googleapis.com/v1/projects/cloudstatus-central/locations/${GAR_LOCATION}/repositories/mon-agent/files/mon-agent-gcp:$version:mon-agent-gcp.$GOARCH:download?alt=media"
curl -fsS -m 30 -X GET -H "Authorization: Bearer $token" -o /usr/local/bin/mon-agent "$url"
chmod +x /usr/local/bin/mon-agent
restorecon -v /usr/local/bin/mon-agent || true

vm_id="$(curl -fsS -m 5 -X GET -H 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/id | cut -d '/' -f 4)"
vm_name="$(curl -fsS -m 5 -X GET -H 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/name | cut -d '/' -f 4)"
vm_type="$(curl -fsS -m 5 -X GET -H 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/machine-type | cut -d '/' -f 4)"
azone="$(curl -fsS -m 5 -X GET -H 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/zone | cut -d '/' -f 4)"
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

${file("../../../etc/mon-agent-systemd.sh")}
