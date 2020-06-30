# DynamicFlare
## Description
DynamicFlare is a DNS utility, designed to programmatically update DNS records found in [Cloudflare](https://www.cloudflare.com/). Its primary purpose is to dynamically update DNS records.

## Installation
### Download a release (recommended)

The latest version of the utility [is available in the "Tags" section of this repository](https://github.com/arekmano/dynamicflare/tags)

### Building from Source

To build from source, Go 1.11 is required.

Clone the repository:

```bash
git clone git@github.com:arekmano/dynamicflare.git
```

Build the utility:

```bash
go build
```

### Using Docker

A docker image is available through Dockerhub. Execute:

```bash
docker pull arekmano/dynamicflare
```

to pull the latest version of the image.

Then run the

## Configuration

### Example

See the [sample configuration file](https://github.com/arekmano/dynamicflare/blob/master/sample.config.yml) included in the project.

### Obtain a Cloudflare API key

[See the Cloudflare documentation](https://support.cloudflare.com/hc/en-us/articles/200167836-Managing-API-Tokens-and-Keys) for details on how to obtain an API key.

The keys are defined in the configuration file in the following way:

```yaml
cloudflare:
  key: ABC123ABC123ABC123 # Cloudflare API key, obtained from Cloudflare
  email: mail@example.com # The email associated with your Cloudflare account
```

### Records

The records to be updated should be added to the configuration. Records must exist previously and are defined in the following format:

```yaml
- id: ABC123ABC123ABC123 # Cloudflare Record ID
  type: A # The record type
  zoneid: ABC123ABC123ABC123 # Cloudflare Zone ID of the record
  name: domain.example.com # The FQDN of the record
```

## Running Commands

### Updating DNS records

The following command will update the DNS entries with the detected public IP:

```bash
dynamicflare -c config.yml update
```


### Updating DNS records

The following command will update the DNS entries with the detected public IP:

```bash
dynamicflare -c config.yml update
```


### Printing all DNS records

The following command will print all of the records associated with the Cloudflare account:

```bash
dynamicflare -c config.yml records
```
