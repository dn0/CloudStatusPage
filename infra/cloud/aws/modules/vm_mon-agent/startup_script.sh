#!/usr/bin/env bash

set -xeo pipefail
/bin/test -f /etc/systemd/system/mon-agent.service && exit 0

which jq > /dev/null || yum install -q -y jq
yum install -q -y https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/linux_${GOARCH}/amazon-ssm-agent.rpm

cat << EOF >> /home/ec2-user/.ssh/authorized_keys
${file("../../../etc/mon-agent.authorized_keys")}
EOF

mkdir -p /usr/local/bin
token="$(curl -fsS -m 5 -X PUT -H "X-aws-ec2-metadata-token-ttl-seconds: 60" http://169.254.169.254/latest/api/token)"
version="$(curl -fsS -m 5 -X GET -H "X-aws-ec2-metadata-token: $token" http://169.254.169.254/latest/meta-data/tags/instance/app_version)"
uri="s3://${ARTIFACTS_BUCKET}/mon-agent/$version/mon-agent-aws.${GOARCH}"
aws s3 cp --region "${SECRET_ENV_REGION}" "$uri" /usr/local/bin/mon-agent
chmod +x /usr/local/bin/mon-agent
restorecon -v /usr/local/bin/mon-agent || true

vm_id="$(curl -fsS -m 5 -X GET -H "X-aws-ec2-metadata-token: $token" http://169.254.169.254/latest/meta-data/instance-id)"
vm_name="$(curl -fsS -m 5 -X GET -H "X-aws-ec2-metadata-token: $token" http://169.254.169.254/latest/meta-data/tags/instance/Name)"
vm_type="$(curl -fsS -m 5 -X GET -H "X-aws-ec2-metadata-token: $token" http://169.254.169.254/latest/meta-data/instance-type)"
azone="$(curl -fsS -m 5 -X GET -H "X-aws-ec2-metadata-token: $token" http://169.254.169.254/latest/meta-data/placement/availability-zone)"
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
aws secretsmanager get-secret-value --query SecretString --output text --region "${SECRET_ENV_REGION}" --secret-id "${SECRET_ENV_ARN}" | jq -r 'to_entries|map("\(.key)=\(.value|tostring|@sh)")|.[]' >> /etc/mon-agent.conf

${file("../../../etc/mon-agent-systemd.sh")}
