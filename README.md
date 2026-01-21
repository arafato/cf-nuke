# cf-nuke: Cloudflare Resource Deleter

`cf-nuke` is a powerful command-line tool designed to **delete all resources** from a Cloudflare account. **Use this tool with extreme caution**, as it is irreversible and cannot distinguish between production and non-production resources.

**WARNING:** This tool is highly destructive. It is strongly advised **not** to run this application on any Cloudflare account where you cannot afford to lose all resources. Always back up critical data and configurations before execution.

---

## Table of Contents

*   [Overview](#overview)
*   [Supported Resources](#supported-resources)
*   [Usage](#usage)
    *   [Dry Run Mode](#dry-run-mode)
    *   [Actual Nuke Execution](#actual-nuke-execution)
    *   [Account Mode](#account-mode)
    *   [Token Mode](#token-mode)
*   [Configuration](#configuration)
*   [Important Considerations](#important-considerations)

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
| Account Tokens | `AccountToken` | API tokens created for the account |
| AI Gateway | `AIGateway` | AI Gateway configurations |
| D1 Databases | `D1` | Serverless SQL databases |
| Hyperdrive | `Hyperdrive` | Database connection pooling configurations |
| KV Namespaces | `KV` | Key-Value storage namespaces |
| Load Balancer Monitors | `LoadBalancerMonitor` | Health check monitors for load balancers |
| Load Balancer Pools | `LoadBalancerPool` | Origin server pools for load balancers |
| Pages Projects | `PagesProject` | Cloudflare Pages deployment projects |
| Pipelines | `Pipeline` | Data pipelines (Pipelines product) |
| Queues | `Queue` | Message queues |
| R2 Buckets | `R2` | Object storage buckets |
| Secrets Stores | `SecretsStore` | Secrets Store configurations |
| Stream Live Inputs | `StreamLiveInput` | Live streaming input configurations |
| Stream Videos | `StreamVideo` | Uploaded and encoded videos |
| Turnstile Widgets | `Turnstile` | CAPTCHA alternative widgets |
| Vectorize Indexes | `Vectorize` | Vector database indexes for AI applications |
| Workers for Platforms | `DispatchNamespace` | Dispatch namespaces for Workers for Platforms |
| Workers Scripts | `WorkersScripts` | Worker scripts and their configurations |
| Workflows | `Workflow` | Workflows definitions |

### Zone-Scoped Resources

These resources are associated with specific DNS zones.

| Resource Type | Config Name | Description |
|---------------|-------------|-------------|
| Custom Hostnames | `CustomHostname` | SSL for SaaS custom hostnames |
| Load Balancers | `LoadBalancer` | Load balancer configurations |
| Waiting Rooms | `WaitingRoom` | Virtual waiting room configurations |
| Zones | `Zone` | DNS zones (domains) |

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

---

## Configuration

`cf-nuke` uses a YAML configuration file to control which resources are targeted for deletion. This is crucial for preventing accidental deletion of important resources.

### Filter Types

You can configure two types of filters:

| Filter Type | Purpose |
|-------------|---------|
| `resource-types.excludes` | Exclude **all instances** of specific resource types |
| `resource-ids.excludes` | Exclude **specific resources** by their ID |

### Example Configuration

```yaml
# config.yaml

# Exclude entire resource types from deletion
resource-types:
  excludes:
    - AccountToken      # Preserve all API tokens
    - Zone              # Preserve all DNS zones
    - WorkersScripts    # Preserve all Worker scripts
    - D1                # Preserve all D1 databases
    - R2                # Preserve all R2 buckets
    - KV                # Preserve all KV namespaces

# Exclude specific resources by ID
resource-ids:
  excludes:
    # Preserve a specific KV namespace
    - resourceType: KV
      id: 609a604e17d24ad0a1bda78cdcf35733

    # Preserve a specific D1 database
    - resourceType: D1
      id: 16d30729-f27e-474a-885c-5d285804a7ac

    # Preserve a specific Worker script
    - resourceType: WorkersScripts
      id: my-production-worker

    # Preserve a specific load balancer
    - resourceType: LoadBalancer
      id: abc123def456
```

### Configuration Reference

| Field | Type | Description |
|-------|------|-------------|
| `resource-types.excludes` | `string[]` | List of resource type names to exclude (see [Supported Resources](#supported-resources) for valid names) |
| `resource-ids.excludes` | `object[]` | List of specific resources to exclude |
| `resource-ids.excludes[].resourceType` | `string` | The resource type (must match a valid Config Name) |
| `resource-ids.excludes[].id` | `string` | The unique identifier of the resource |

> **Tip:** Run `cf-nuke` in dry-run mode first (without `--no-dry-run`) to see all discovered resources and their IDs before configuring exclusions.

---

## Important Considerations

*   **Destructive Operation:** `cf-nuke` is designed to delete resources. Once deleted, they cannot be recovered. **Verify your actions meticulously.**
*   **Account Security:** Protect your Cloudflare API tokens and account IDs. Avoid committing them directly into version control. Use environment variables or secure secret management practices.
*   **Backup:** Before running `cf-nuke`, ensure you have backups of any critical data or configurations managed within Cloudflare (e.g., DNS records, Worker scripts, KV data if not explicitly excluded).
*   **Targeted Use:** This tool is best suited for development accounts, testing environments, or situations where a complete account reset is intentionally desired. **Do not use it on production accounts unless you fully understand and accept the consequences.**
