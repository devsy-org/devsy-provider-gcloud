# GCLOUD Provider for Devsy

This provider runs workspaces on [Google Cloud](https://cloud.google.com/) virtual machines.

## Getting started

You need a Google Cloud account with billing, a project, and the Compute Engine API.

```sh
gcloud projects create <your-devsy-project-id>
gcloud services enable compute.googleapis.com --project <your-devsy-project-id>
```

You can manage the project in the
[Google Cloud Console](https://console.cloud.google.com/),
in [Cloud Shell](https://docs.cloud.google.com/shell/docs/how-cloud-shell-works),
or with the [Google Cloud CLI](https://docs.cloud.google.com/sdk/gcloud).

Install the CLI on the machine that runs Devsy.
The desktop app authenticates with it.
See [Application Default Credentials](#application-default-credentials).

Add the provider:

```sh
devsy provider add gcloud \
  -o PROJECT=<your-devsy-project-id> \
  -o ZONE=<Google Cloud zone to create the VMs in>
devsy provider use gcloud
```

`PROJECT` must be set, unless `gcloud` already has a current project.
`ZONE` should be set when you add the provider.

A value given with `-o` wins over an environment variable of the same option.
Devsy fills an option that was not passed on the command line from
`DEVSY_PROVIDER_<provider-name>_<OPTION>`.
For this provider that is `DEVSY_PROVIDER_GCLOUD_PROJECT` and
`DEVSY_PROVIDER_GCLOUD_ZONE`.
If you add the provider under another name, the middle segment is that name
in upper case. Hyphens are kept.
`devsy provider add gcloud --name mygcloud` reads
`DEVSY_PROVIDER_MYGCLOUD_PROJECT`.
A hyphen in the provider name cannot be used in a shell variable, so pick a
name without one.

```sh
export DEVSY_PROVIDER_GCLOUD_PROJECT=<your-devsy-project-id>
export DEVSY_PROVIDER_GCLOUD_ZONE=<Google Cloud zone to create the VMs in>
devsy provider add gcloud
```

Change an option later with `devsy provider set`:

```sh
devsy provider set gcloud -o DISK_IMAGE=my-custom-vm-image
```

Start a workspace:

```sh
devsy workspace up .
devsy workspace up github.com/<user>/<repository>
```

## Authentication

Use one of the two approaches below.

### Application Default Credentials

`gcloud auth application-default login` writes a credential file where the
Google Cloud libraries look for it:

- Linux and macOS: `$HOME/.config/gcloud/application_default_credentials.json`
- Windows: `%APPDATA%\gcloud\application_default_credentials.json`

```sh
gcloud auth application-default login
```

Details:
[Application Default Credentials](https://developers.google.com/accounts/docs/application-default-credentials)
and
[ADC for a user account](https://docs.cloud.google.com/docs/authentication/application-default-credentials#personal).

This is the approach that works from the Devsy desktop app.

### Service-account key

Pass the key as a path. `KEY_FILE` is a provider option:

```sh
export DEVSY_PROVIDER_GCLOUD_KEY_FILE=/path/to/service/account/key.json
devsy provider add gcloud
```

`provider add` and `provider set` copy that variable into the saved option,
because `KEY_FILE` is declared. The saved value is the path.

The JSON contents are not an option. The provider reads
`DEVSY_PROVIDER_<name>_KEY`, and `GCLOUD_JSON_AUTH` when that is empty, and
does not store either value. Prefer `KEY_FILE`. A declared `KEY` option would
be copied into the Devsy config on `provider add`.

Only a service account has this kind of key.
The desktop app does not load these variables for you, so this approach is
for the CLI.
The benefit is more than one Google Cloud organization at a time: each
directory can export a different key, project, and `DEVSY_CONTEXT`.

## Service account

A service account keeps the virtual-machine credentials separate from your
user. See
[Using IAM securely](https://cloud.google.com/iam/docs/using-iam-securely)
and the
[service account overview](https://docs.cloud.google.com/iam/docs/service-account-overview).

```sh
gcloud iam service-accounts create <your-devsy-service-account> \
  --project <your-devsy-project-id>
```

The email is
`<your-devsy-service-account>@<your-devsy-project-id>.iam.gserviceaccount.com`.

```sh
gcloud iam service-accounts keys create /path/to/service/account/key.json \
  --iam-account=<your-devsy-service-account-email> \
  --project <your-devsy-project-id>
```

Grant these roles on the project:

```sh
# Compute Engine billing
gcloud projects add-iam-policy-binding <your-devsy-project-id> \
  --member="serviceAccount:<your-devsy-service-account-email>" \
  --role=roles/serviceusage.serviceUsageConsumer
# Compute Engine instance operations
gcloud projects add-iam-policy-binding <your-devsy-project-id> \
  --member="serviceAccount:<your-devsy-service-account-email>" \
  --role=roles/compute.instanceAdmin.v1
```

If the VM itself should call Google Cloud, create a second service account,
attach it with the `SERVICE_ACCOUNT` option, and also grant
`roles/iam.serviceAccountUser` to the account that creates the VM.
See
[service accounts for instances](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances).

## Options

| NAME | REQUIRED | DESCRIPTION | DEFAULT |
| --- | --- | --- | --- |
| PROJECT | true | The project id to use. | |
| ZONE | true | The Google Cloud zone to create the VM in, e.g. europe-west1-d. | |
| KEY_FILE | false | Path to a service-account JSON key. | |
| DISK_IMAGE | false | The disk image to use. | projects/cos-cloud/global/images/cos-101-17162-127-5 |
| DISK_SIZE | false | The disk size to use (GB). | 40 |
| MACHINE_TYPE | false | The machine type to use. | c2-standard-4 |
| NETWORK | false | The network id to use. | |
| SUBNETWORK | false | The subnetwork id to use. | |
| TAG | false | A tag to attach to the instance. | devsy |
| SERVICE_ACCOUNT | false | A service account to attach to the instance. | |
| PUBLIC_IP_ENABLED | false | Use a public IP to access the instance. | true |
| AGENT_PATH | false | The path where to inject the Devsy agent. | /var/lib/toolbox/devsy |
| INACTIVITY_TIMEOUT | false | Stop the VM after this idle period. | 5m |
| INJECT_DOCKER_CREDENTIALS | false | Inject Docker credentials into the remote host. | true |
| INJECT_GIT_CREDENTIALS | false | Inject Git credentials into the remote host. | true |

## CLI workflow

[direnv](https://direnv.net/) is a convenient way to keep one Google Cloud
organization per directory. Keep the `.envrc` files in the dotfiles repo.

```sh
export DEVSY_CONTEXT="org1"
export DEVSY_PROVIDER="gcloud"
export DEVSY_PROVIDER_GCLOUD_KEY_FILE=/path/to/service/account/key.json
export DEVSY_PROVIDER_GCLOUD_PROJECT="<your-devsy-project-id>"
export DEVSY_PROVIDER_GCLOUD_ZONE="us-east4-b"
export DEVSY_PROVIDER_GCLOUD_MACHINE_TYPE="c2-standard-8"
export DEVSY_PROVIDER_GCLOUD_DISK_SIZE=20
```

`DEVSY_CONTEXT` and `DEVSY_PROVIDER` are wired to the global flags.
`DEVSY_IDE` and `DEVSY_WORKSPACE_ENV` are not.
Set the IDE with `devsy ide use intellij`.
Pass workspace environment on the command that starts the workspace.

With that file loaded:

```sh
devsy provider add gcloud
devsy ide use intellij
# Skip the JetBrains trust dialog the first time a workspace opens.
devsy workspace up . --workspace-env REMOTE_DEV_TRUST_PROJECTS=1
```

A second organization is a second directory, a second `.envrc`, and a
different `DEVSY_CONTEXT`.
Contexts scope the providers added to them.
