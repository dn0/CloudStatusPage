# infra/cloud

> [!NOTE]
> Some infra parts are missing and/or were obfuscated.

## Prerequisites

### GCP

```sh
gcloud auth application-default login
gcloud auth application-default set-quota-project cloudstatus-central
```

### AWS

```sh
aws configure sso
# afterwards use `aws sso login`
```

### Azure

```sh
az login
az vm image terms accept --urn resf:rockylinux-x86_64:9-base:latest --subscription ... # run for every environment
```


## New Cloud Account

```sh
tf init
tf providers lock -platform=linux_amd64 -platform=darwin_amd64 -platform=linux_arm64 -platform=darwin_arm64
```

### New Google Cloud Project

1. Run tf apply; if it fails run it again.
2. Disable or delete the `PROJECT_NUMBER-compute@developer.gserviceaccount.com` service account.
3. Delete the `default` (or other existing networks) in GCP.
4. Check `gcloud services list --project PROJECT_ID` and compare it with google_project_service resources.
   You may need to disable some services with `gcloud services disable --project PROJECT_ID SERVICE`


### New AWS account

1. Root IAM user is using a specific email pattern.
2. Keep the default VPC although we will move away from it.

### New Azure account

1. Create a subscription manually
2. You may need to delete the automatically created alias: `az account alias delete -n <some-uid>`
3. Run tf apply in the org folder + `tf import 'azurerm_management_group_subscription_association.root["<name>"]' /managementGroup/<uuid>/subscription/<id>`
4. Remove the Owner role assignment in the subscription's IAM policy
   (keep the inherited one from the root management group)


## SSH

You can use some `make` commands (`bin/cloud-ssh.sh`) to run SSH commands across all VMs.

### GCP

`gcloud compute ssh` should work out-of-the-box.
SSH is open only through the Identity-Aware Proxy.

```sh
gcloud compute instances list --project=<project-id>
gcloud compute ssh <instance-name> --project=<project-id> --zone=<instance-zone> --tunnel-through-iap
```

### AWS

Setup SSM:

1. https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html
2. https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-getting-started-enable-ssh-connections.html

```sh
aws ec2 describe-regions --output table
aws ec2 describe-instances --region "<region>" --output table --query 'Reservations[*].Instances[*].[InstanceId, Placement.AvailabilityZone, PrivateIpAddress, PublicIpAddress, State.Name, Tags[?Key==`app_version`]|[0].Value]'
aws ssm start-session --target <instance-id>
```

### Azure

Azure does not have a SSH proxy/IAM service :(
Port 22 is restricted via security groups in `cloud/azure/modules/config/vpc.tf`.
There can be only one SSH key so we will use a shared key and instruct the startup script to add more authorized keys.

```sh
az vm list --subscription "<account-id>"
ssh -l azureuser <vm-public-ip>'
```


## Secrets

mon-agent needs a GCP service account to connect to Pub/Sub.
To obtain a service account key:
```sh
gcloud iam service-accounts keys create key.json --iam-account=mon-agent-gcp@cloudstatus-<env>.iam.gserviceaccount.com --project=cloudstatus-<env>
```

### GCP

No secrets are needed for mon-agent.

### AWS

One secret - `mon-agent/env` - stored in one AWS region is read during deployment of mon-agents in all regions.
Secret value can contain multiple env vars. To store the GCP service account:
```sh
SECRET=$(jq -cnrM '{GOOGLE_APPLICATION_CREDENTIALS_JSON: $key}' --arg key "$(jq -rcM . key.json)")
aws secretsmanager create-secret --region $AWS_REGION --name mon-agent/env --secret-string "$SECRET"
```

### Azure

One secret stored in eu-central-1 key vault - `https://mon-probe-<env>/vault.azure.net/secrets/mon-agent-env` - is read during deployment of mon-agents in all regions. Secret value can contain multiple env vars. To store the GCP service account:
```sh
SECRET=$(jq -cnrM '{GOOGLE_APPLICATION_CREDENTIALS_JSON: $key}' --arg key "$(jq -rcM . key.json)")
az keyvault secret set --subscription <subscription-id> --name mon-agent-env --vault-name mon-probe-<env> --value "$SECRET"
```

# infra/core

## Kubernetes

### PostgreSQL

1. Install CloudNativePG operator:
    ```sh
    kubectl apply --server-side -f \
      https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.24/releases/cnpg-1.24.1.yaml
    kubectl get deployment,pod,svc -n cnpg-system
    ```
    Note: If you hit a GKE firewall issue please read: https://cloudnative-pg.io/documentation/1.24/installation_upgrade/

2. Set passwords for DB roles and create application secrets with DB connections strings:
    ```bash
    for app in mon-scribe mon-analyst mon-web; do
        password="$(LC_ALL=C tr -dc 'A-Za-z0-9' < /dev/urandom | head -c 32)"
        db_write_url="postgresql://${app}:${password}@core-rw.postgres.svc.cluster.local.:5432/app?pool_max_conns=10"
        if [[ "$app" == "mon-scribe" ]]; then
            db_read_url="-"
        else
            db_read_url="postgresql://${app}:${password}@core-r.postgres.svc.cluster.local.:5432/app?pool_max_conns=10"
        fi
        NAME="${app}" \
        DB_USERNAME="$(echo -ne "$app" | base64)" \
        DB_PASSWORD="$(echo -ne "$password" | base64)" \
        DB_WRITE_URL="$(echo -ne "$db_write_url" | base64)" \
        DB_READ_URL="$(echo -ne "$db_read_url" | base64)" \
        envsubst '${DB_USERNAME} ${DB_PASSWORD} ${DB_WRITE_URL} ${DB_READ_URL} ${NAME}' < k8s/postgres/secret.tmpl.yaml | kubectl apply -f -
    done
    ```

3. Deploy a PostgreSQL cluster with the TimescaleDB extension:
    ```sh
    kubectl apply -k k8s/postgres/<env>
    kubectl get clusters,poolers,pod,svc -n postgres
    ```
    Note: Make sure that the PostgreSQL image is accessible from the cluster.
    Note: Consider installing the [cnpg plugin](https://cloudnative-pg.io/documentation/current/kubectl-plugin/) so you can do things like:
        ```sh
        kubectl cnpg status core -n postgres
        ```

4. Create the necessary Timescale DB extensions in the `app` DB. Although, they could be created by the application's `make db/sql` initialization scripts, the `timescaledb_toolkit` requires superuser permissions so it's easier to do it here:
    ```sql
    $ kubectl cnpg psql core -n postgres

    \c app
    CREATE EXTENSION timescaledb;
    CREATE EXTENSION timescaledb_toolkit;
    ```

5. Initialize the app DB
    To access the DB from localhost:
    ```sh
    kubectl port-forward -n postgres services/core-rw 5432:5432
    ```
    Create app tables and indexes:
    ```sh
    export DATABASE_URL=$(kubectl get secret -n postgres core-app -o json | jq -r .data.uri | base64 -D | sed 's/core-rw\.postgres/localhost/')
    cd src
    make db/init
    ```

#### Upgrade

1. Upgrade CloudNativePG operator - apply installation step#1:
    ```sh
    kubectl apply --server-side -f ...cnpg....yaml
    ```
    Note: This triggers a rolling restart of the cluster.

2. Upgrade PostgreSQL cluster - apply installation step#3
    ```sh
    kubectl apply -k k8s/postgres/<env>
    ```

3. Upgrade Timescale DB extensions:
    ```sql
    $ kubectl cnpg psql core -n postgres -- --dbname=app -P pager=off

    \c app
    ALTER EXTENSION timescaledb UPDATE;
    ALTER EXTENSION timescaledb_toolkit UPDATE;
    ```


