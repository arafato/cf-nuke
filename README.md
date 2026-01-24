<p align="center">
  <img src="https://img.shields.io/badge/cf--nuke-F38020?style=for-the-badge&logo=cloudflare&logoColor=white" alt="cf-nuke" />
</p>

<h1 align="center">cf-nuke</h1>
<p align="center"><strong>Safely delete all resources from a Cloudflare account</strong></p>

<p align="center">
  <a href="https://github.com/arafato/cf-nuke/blob/main/LICENSE"><img src="https://img.shields.io/github/license/arafato/cf-nuke.svg" alt="license" /></a>
  <a href="https://github.com/arafato/cf-nuke/releases"><img src="https://img.shields.io/github/release/arafato/cf-nuke.svg" alt="release" /></a>
  <a href="https://goreportcard.com/report/github.com/arafato/cf-nuke"><img src="https://goreportcard.com/badge/github.com/arafato/cf-nuke" alt="Go Report Card" /></a>
  <img src="https://img.shields.io/github/downloads/arafato/cf-nuke/total" alt="downloads" />
  <a href="https://github.com/arafato/homebrew-tap"><img src="https://img.shields.io/badge/brew-arafato/tap/cf--nuke-FBB040?logo=homebrew&logoColor=white" alt="homebrew" /></a>
</p>

---

> **WARNING:** This tool is highly destructive. It is designed to **delete all resources** from a Cloudflare account and cannot distinguish between production and non-production resources. **Use with extreme caution.** Always back up critical data and configurations before execution.

---

## Table of Contents

*   [Installation](#installation)
*   [Overview](#overview)
*   [Supported Resources](#supported-resources)
*   [Usage](#usage)
    *   [List Command](#list-command)
    *   [Dry Run Mode](#dry-run-mode)
    *   [Actual Nuke Execution](#actual-nuke-execution)
    *   [Account Mode](#account-mode)
    *   [Token Mode](#token-mode)
*   [Configuration](#configuration)
*   [Important Considerations](#important-considerations)

---

## Installation

### Homebrew (macOS and Linux)

```bash
brew tap arafato/tap
brew install cf-nuke
```

Or install directly:

```bash
brew install arafato/tap/cf-nuke
```

### Download Binary

Download the latest release for your platform from the [Releases page](https://github.com/arafato/cf-nuke/releases).

Available binaries:
- `cf-nuke_*_darwin_arm64.tar.gz` - macOS (Apple Silicon)
- `cf-nuke_*_darwin_amd64.tar.gz` - macOS (Intel)
- `cf-nuke_*_linux_arm64.tar.gz` - Linux (ARM64)
- `cf-nuke_*_linux_amd64.tar.gz` - Linux (x86_64)
- `cf-nuke_*_windows_amd64.zip` - Windows (x86_64)

### Build from Source

Requires Go 1.24 or later:

```bash
git clone https://github.com/arafato/cf-nuke.git
cd cf-nuke
go build -o cf-nuke .
```

---

## Overview

`cf-nuke` provides a way to clean up all resources within a Cloudflare account. It is intended for use cases where a complete reset of an account's resources is necessary. Due to its destructive nature, careful planning and execution are paramount.

---

## Supported Resources

`cf-nuke` can delete the following Cloudflare resource types:

### Account-Scoped Resources

These resources are managed at the account level.

| Resource Type | Config Name | Description |
|---------------|-------------|-------------|
| ZT Access Applications | `ZTAccessApplication` | Zero Trust Access applications (SaaS, self-hosted, etc.) |
| Account Tokens | `AccountToken` | API tokens created for the account |
| AI Gateway | `AIGateway` | AI Gateway configurations |
| Calls SFU Apps | `CallsApp` | Cloudflare Calls SFU (video/audio) applications |
| Calls TURN Keys | `CallsTurnKey` | Cloudflare Calls TURN server keys |
| D1 Databases | `D1` | Serverless SQL databases |
| Hyperdrive | `Hyperdrive` | Database connection pooling configurations |
| Images | `Image` | Cloudflare Images stored in the account |
| KV Namespaces | `KV` | Key-Value storage namespaces |
| Load Balancer Monitors | `LoadBalancerMonitor` | Health check monitors for load balancers |
| Load Balancer Pools | `LoadBalancerPool` | Origin server pools for load balancers |
| Logpush Jobs | `LogpushJob` | Log delivery job configurations |
| Pages Projects | `PagesProject` | Cloudflare Pages deployment projects |
| Pipelines | `Pipeline` | Data pipelines (Pipelines product) |
| Queues | `Queue` | Message queues |
| R2 Buckets | `R2` | Object storage buckets |
| Rulesets | `Ruleset` | Custom rulesets (excludes Cloudflare-managed) |
| RUM Sites | `RUMSite` | Real User Monitoring (Web Analytics) sites |
| Secrets Stores | `SecretsStore` | Secrets Store configurations |
| Stream Live Inputs | `StreamLiveInput` | Live streaming input configurations |
| Stream Videos | `StreamVideo` | Uploaded and encoded videos |
| Turnstile Widgets | `Turnstile` | CAPTCHA alternative widgets |
| Vectorize Indexes | `Vectorize` | Vector database indexes for AI applications |
| Workers for Platforms | `DispatchNamespace` | Dispatch namespaces for Workers for Platforms |
| Workers Scripts | `WorkersScripts` | Worker scripts and their configurations |
| Workflows | `Workflow` | Workflows definitions |

### Security & Certificates (Account-Scoped)

| Resource Type | Config Name | Description |
|---------------|-------------|-------------|
| MTLS Certificates | `MTLSCertificate` | Mutual TLS client certificates for API Shield |
| Firewall Access Rules | `FirewallAccessRule` | IP-based access rules (allow/block/challenge) |

### Zero Trust & Networking (Account-Scoped)

| Resource Type | Config Name | Description |
|---------------|-------------|-------------|
| ZT Access Groups | `ZTAccessGroup` | Zero Trust Access policy groups |
| ZT Service Tokens | `ZTServiceToken` | Zero Trust service-to-service authentication tokens |
| ZT Bookmarks | `ZTBookmark` | Zero Trust application bookmarks |
| Cloudflare Tunnels | `Tunnel` | Cloudflare Tunnel connections |
| Virtual Networks | `VirtualNetwork` | Virtual networks for tunnel routing |
| Tunnel Routes | `TunnelRoute` | Private network routes through tunnels |
| Gateway Rules | `GatewayRule` | Zero Trust Gateway filtering rules |
| DNS Firewall | `DNSFirewall` | DNS Firewall clusters |

### Zone-Scoped Resources

These resources are associated with specific DNS zones.

| Resource Type | Config Name | Description |
|---------------|-------------|-------------|
| Custom Hostnames | `CustomHostname` | SSL for SaaS custom hostnames |
| DNS Records | `DNSRecord` | DNS records (A, AAAA, CNAME, MX, TXT, etc.) |
| Healthchecks | `Healthcheck` | Standalone health check configurations |
| Load Balancers | `LoadBalancer` | Load balancer configurations |
| Snippets | `Snippet` | Cloudflare Snippets (edge code snippets) |
| Spectrum Apps | `SpectrumApp` | Spectrum TCP/UDP proxy applications |
| Waiting Rooms | `WaitingRoom` | Virtual waiting room configurations |
| Web3 Hostnames | `Web3Hostname` | Web3 gateway hostnames (IPFS, ENS) |
| Zones | `Zone` | DNS zones (domains) |

### Security & Certificates (Zone-Scoped)

| Resource Type | Config Name | Description |
|---------------|-------------|-------------|
| Custom Certificates | `CustomCertificate` | Custom SSL/TLS certificates |
| Origin CA Certificates | `OriginCACertificate` | Cloudflare-signed origin certificates |
| Client Certificates | `ClientCertificate` | mTLS client certificates for API Shield |
| Keyless SSL Certificates | `KeylessCertificate` | Keyless SSL configurations |
| API Gateway Operations | `APIGatewayOperation` | API Gateway endpoint operations |
| Page Shield Policies | `PageShieldPolicy` | Page Shield CSP policies |

> **Note:** Use the **Config Name** values when specifying resource types in your `config.yaml` file for filtering (see [Configuration](#configuration)).

---

## Usage

The primary command for deleting resources is `nuke`. It accepts several flags to configure the operation.

```bash
$ cf-nuke nuke [flags]

Flags:
  -a, --account-id string   Cloudflare account ID (required for Token mode)
  -c, --config string       Path to configuration file (required)
  -h, --help                help for nuke
  -k, --key string          API Key or Token (required)
  -m, --mode string         The mode of operation ('token' or 'account')
      --no-dry-run          Execute without dry run (perform actual deletion)
  -u, --user string         User identifier (required only for 'account' mode)
```

### List Command

To see all supported resource types, use the `list` command:

```bash
$ cf-nuke list
```

This outputs an alphabetically sorted list of all resource collector names:

```
account-token
ai-gateway
api-gateway-operation
calls-app
calls-turn-key
client-certificate
custom-certificate
custom-hostname
d1
...
```

This is useful for understanding what resources `cf-nuke` can manage and for reference when configuring exclusions.

### Dry Run Mode

By default, `cf-nuke` operates in dry-run mode. This is a safety measure that lists all resources it identifies and indicates which would be deleted based on your configuration. **Always perform a dry run first.**

**Example:**
```bash
# Using Token Mode (Recommended for most use cases)
cf-nuke nuke \
  --account-id YOUR_ACCOUNT_ID \
  --key YOUR_API_TOKEN \
  -c path/to/your/config.yaml

# Using Account Mode (requires Account Global Token and user email)
./cf-nuke nuke \
  --mode account \
  -k YOUR_ACCOUNT_GLOBAL_TOKEN \
  -u YOUR_EMAIL \
  -c path/to/your/config.yaml
```

**Dry Run Output Example:**
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
In this output:
*   **Ready:** Resources that `cf-nuke` has identified and will attempt to delete (if `--no-dry-run` is used).
*   **Filtered:** Resources that have been explicitly excluded by your configuration.

### Permission Warnings

If `cf-nuke` encounters permission errors while scanning certain resources or zones, it will continue scanning and display a summary of warnings after collection completes:

```
[WARNINGS] 3 issue(s) encountered during collection:
  - DNSRecord (example.com): insufficient permissions
  - SpectrumApp (test.io): insufficient permissions
  - Image: insufficient permissions or feature not available
```

This ensures that missing permissions for specific resources or zones don't stop the entire operation. Resources that couldn't be scanned due to permissions will simply be skipped.

### Actual Nuke Execution

To proceed with the deletion of resources, you must use the `--no-dry-run` flag. This will trigger a final confirmation prompt before any actions are taken.

**Example:**
```bash
cf-nuke nuke \
  --account-id YOUR_ACCOUNT_ID \
  --key YOUR_API_TOKEN \
  -c path/to/your/config.yaml \
  --no-dry-run
```

Upon running this command, you will be prompted:
```
Executing actual nuke operation... do you really want to continue (yes/no)?
```
You must type `yes` and press Enter to confirm the deletion. **There is no undo.**

### Deletion Errors

If any resources fail to delete during the nuke operation, `cf-nuke` will display a summary of failed resources with their error messages at the end:

```
[ERRORS] 2 resource(s) failed to delete:
  - KV (my-namespace): resource is in use by a Worker
  - CustomCertificate (example.com): certificate is currently active
```

This helps you identify resources that require manual intervention or have dependencies that need to be resolved first.

---

## Configuration

`cf-nuke` uses a YAML configuration file to control which resources are targeted for deletion. This is crucial for preventing accidental deletion of important resources.

### Filter Types

You can configure three types of filters:

| Filter Type | Purpose |
|-------------|---------|
| `zones.excludes` | Exclude **entire DNS zones** (and their zone-scoped resources) by domain name |
| `resource-types.excludes` | Exclude **all instances** of specific resource types |
| `resource-ids.excludes` | Exclude **specific resources** by their ID or name |

### Example Configuration

```yaml
# config.yaml

# Exclude entire zones from deletion
# This also protects all zone-scoped resources within these zones
# Note: Zones using Cloudflare Registrar cannot be deleted
zones:
  excludes:
    - mydomain.com        # Preserve this zone and all its resources
    - production.io       # Preserve production domain

# Exclude entire resource types from deletion
resource-types:
  excludes:
    - AccountToken      # Preserve all API tokens
    - Zone              # Preserve all DNS zones
    - WorkersScripts    # Preserve all Worker scripts
    - D1                # Preserve all D1 databases
    - R2                # Preserve all R2 buckets
    - KV                # Preserve all KV namespaces

# Exclude specific resources by ID or name
resource-ids:
  excludes:
    # Preserve a specific KV namespace
    - resourceType: KV
      id: 609a604e17d24ad0a1bda78cdcf35733

    # Preserve a specific D1 database
    - resourceType: D1
      id: 16d30729-f27e-474a-885c-5d285804a7ac

    # Preserve a specific Worker script (by name)
    - resourceType: WorkersScripts
      id: my-production-worker

    # Preserve a specific load balancer
    - resourceType: LoadBalancer
      id: abc123def456
```

### Configuration Reference

| Field | Type | Description |
|-------|------|-------------|
| `zones.excludes` | `string[]` | List of zone domain names to exclude from deletion |
| `resource-types.excludes` | `string[]` | List of resource type names to exclude (see [Supported Resources](#supported-resources) for valid names) |
| `resource-ids.excludes` | `object[]` | List of specific resources to exclude |
| `resource-ids.excludes[].resourceType` | `string` | The resource type (must match a valid Config Name) |
| `resource-ids.excludes[].id` | `string` | The unique identifier or name of the resource |

> **Tip:** Run `cf-nuke` in dry-run mode first (without `--no-dry-run`) to see all discovered resources and their IDs before configuring exclusions.

---

## Important Considerations

*   **Destructive Operation:** `cf-nuke` is designed to delete resources. Once deleted, they cannot be recovered. **Verify your actions meticulously.**
*   **Account Security:** Protect your Cloudflare API tokens and account IDs. Avoid committing them directly into version control. Use environment variables or secure secret management practices.
*   **Backup:** Before running `cf-nuke`, ensure you have backups of any critical data or configurations managed within Cloudflare (e.g., DNS records, Worker scripts, KV data if not explicitly excluded).
*   **Targeted Use:** This tool is best suited for development accounts, testing environments, or situations where a complete account reset is intentionally desired. **Do not use it on production accounts unless you fully understand and accept the consequences.**
