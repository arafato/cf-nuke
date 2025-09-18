# cf-nuke

`cf-nuke` is a command-line tool for removing all resources from a Cloudflare account. Use it with caution, as it cannot distinguish between production and non-production resources.

## Usage

```bash
$ cf-nuke nuke [flags]

Flags:
  -a, --account-id string   Cloudflare account id (required)
  -c, --config string       Path to configuration file (required)
  -h, --help                help for nuke
  -k, --key string          Key for operation (required)
  -m, --mode string         The mode of operation ('token' or 'account')
      --no-dry-run          Execute without dry run
  -u, --user string         The user identifier (required only for 'account' mode)
```

The primary command is `nuke`, which performs the destructive operations.

#### Account Mode (using an Account Global Token)
```bash
$ ./cf-nuke nuke --mode account -k <your-api-token> -u <your-email> -c config.yaml --no-dry-run
```

#### Token Mode (using an API Token)
```bash
cf-nuke nuke -c config.yaml -a <your-account-id> -k <your-api-token>
```

### Dry Run

By default, `cf-nuke` runs in dry-run mode. It will list all the resources it has found and which of them will be removed based on the configuration.

```
┌──────────────┬─────────────┬──────────┐
│   PRODUCT    │  ID / NAME  │  STATUS  │
├──────────────┼─────────────┼──────────┤
│ KV           │ testkv2     │ Ready    │
│ KV           │ testkv1     │ Ready    │
│ AccountToken │ donotdelete │ Filtered │
└──────────────┴─────────────┴──────────┘

Status: 3 resources in total. Removed 0, In-Progress 0, Filtered 1
```

To execute the removal of resources, you must use the `--no-dry-run` flag.

### Executing the Nuke

When you run with `--no-dry-run`, you will be prompted for a final confirmation before the deletion begins.

```bash
Executing actual nuke operation... do you really want to continue (yes/no)?
```

## Configuration

`cf-nuke` lets you configure two different kinds of filters.
- Resource-Type filters that filter all instances of a particular resource type (e.g. all KV instances).
- Resource-ID filters that filter one particular resource instance based on its ID or name.

```yaml
resource-types:
  excludes:
    - SecretsStore
    - WorkersScripts
    - KV
    - AIGateway
    - D1
    - R2
    - Queue
    - AccountToken

resource-ids:
  excludes:
    - resourceType: KV
      id: 609a604e17d24ad0a1bda78cdcf35733
    - resourceType: KV
      id: idef45a748fad047b1b985d7bf29b5d8f3
    - resourceType: D1
      id: 16d30729-f27e-474a-885c-5d285804a7ac
```

## Caution!

Be aware that *cf-nuke* is a very destructive tool, hence you have to be very
careful while using it. Otherwise you might delete production data.

**We strongly advise you to not run this application on any Cloudflare account, where
you cannot afford to lose all resources.**

To reduce the blast radius of accidents, there are some safety precautions:

1. By default *cf-nuke* only lists all nukeable resources. You need to add
   `--no-dry-run` to actually delete resources.

2. more to come...
