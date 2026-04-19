# GCLOUD Provider for Devsy

[![Open in Devsy!](https://img.shields.io/badge/open_in_devsy-8A2BE2?style=for-the-badge)](https://devsy.sh/open#https://github.com/devsy-org/devsy-provider-gcloud)

## Getting started

The provider is available for auto-installation using:

```sh
devsy provider add gcloud -o PROJECT=<project id to use> -o ZONE=<Google Cloud zone to create the VMs in>
devsy provider use gcloud
```

Option `PROJECT` must be set when adding the provider
(unless the project to be used is set as the current project in `gcloud`).

Option `ZONE` should be set when adding the provider.

Options can be set using `devsy provider set-options`, for example:

```sh
devsy provider set-options gcloud -o DISK_IMAGE=my-custom-vm-image
```

Be aware that authentication is obtained using `gcloud` CLI tool, take a look
[here](https://developers.google.com/accounts/docs/application-default-credentials)
for more information.

### Creating your first workspace with gcloud

After the initial setup, just use:

```sh
devsy up .
```

You'll need to wait for the machine and workspace setup.

### Customize the VM Instance

This provider has the following options:

| NAME           | REQUIRED | DESCRIPTION                                                    | DEFAULT                                              |
|----------------|----------|----------------------------------------------------------------|------------------------------------------------------|
| DISK_IMAGE     | false    | The disk image to use.                                         | projects/cos-cloud/global/images/cos-101-17162-127-5 |
| DISK_SIZE      | false    | The disk size to use (GB).                                     | 40                                                   |
| MACHINE_TYPE   | false    | The machine type to use.                                       | c2-standard-4                                        |
| PROJECT        | true     | The project id to use.                                         |                                                      |
| ZONE           | true     | The google cloud zone to create the VM in. E.g. europe-west1-d |                                                      |
| NETWORK        | false    | The network id to use.                                         |                                                      |
| SUBNETWORK     | false    | The subnetwork id to use.                                      |                                                      |
| TAG            | false    | A tag to attach to the instance.                               | devsy                                                |
| SERVICE_ACCOUNT| false    | A service account to attach to instance                        |                                                      |
