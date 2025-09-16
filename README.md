# cf-nuke
Removes all resources from a Cloudflare account.

> **Development Status** *cf-nuke* is not stable and currently under heavy development. It is also likely that not all Cloudflare
resources are covered by it. Be encouraged to add missing resources and create
a Pull Request or to create an [Issue](https://github.com/arafato/cf-nuke/issues/new).

## Caution!

Be aware that *cf-nuke* is a very destructive tool, hence you have to be very
careful while using it. Otherwise you might delete production data.

**We strongly advise you to not run this application on any Cloudflare account, where
you cannot afford to lose all resources.**

To reduce the blast radius of accidents, there are some safety precautions:

1. By default *cf-nuke* only lists all nukeable resources. You need to add
   `--no-dry-run` to actually delete resources.

2. more to come...
