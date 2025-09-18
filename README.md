# cf-nuke
Removes all resources from a Cloudflare account.

```
$ ./cf-nuke nuke --mode account -a 4e6... -k d2c... -u a***@***.com -c config.yaml --no-dry-run```
```
┌──────────────┬─────────────┬──────────┐
│   PRODUCT    │  ID / NAME  │  STATUS  │
├──────────────┼─────────────┼──────────┤
│ KV           │ testkv2     │ Ready    │
│ KV           │ testkv1     │ Ready    │
│ AccountToken │ donotdelete │ Filtered │
└──────────────┴─────────────┴──────────┘

Status: 3 resources in total. Removed 0, In-Progress 0, Filtered 1
Executing actual nuke operation... do you really want to continue (yes/no)?
```
> **Development Status** *cf-nuke* is not stable and currently under heavy development. It is also likely that not all Cloudflare
resources are covered by it. Be encouraged to add missing resources and create
a Pull Request or to create an [Issue](https://github.com/arafato/cf-nuke/issues/new).

## Configuration
cf-nuke lets you configure two different kinds of filters.
- Resource-Type filters that filter all instances of a particular resource type (e.g. all KV instances).
- Resource-ID filters that filter one particular resource instance based on its ID or name.

```
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
