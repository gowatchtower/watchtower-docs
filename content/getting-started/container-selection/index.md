---
title: Container Selection
weight: 4
---
## Overview

By default, Watchtower monitors and updates all running containers on the connected Docker daemon.
This behavior can be customized through a combination of configuration options and container labels to control exactly which containers are managed and how.

Container selection works through a **filter chain**: a series of criteria applied in sequence.
A container is monitored only if it passes every filter in the chain. The filters are evaluated in the following order:

| #  | Filter                                                       | Description                                                                                                                                           |
|----|--------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1  | Old Watchtower container exclusion                           | Watchtower containers renamed with a `watchtower-old-` prefix during self-updates are always excluded.                                                |
| 2  | [Disabled label check](#enabledisable-labels)                | Containers with the Docker label `com.centurylinklabs.watchtower.enable=false` are excluded.                                                          |
| 3  | [Scope filter](#monitoring-scopes)                           | Only containers matching the configured scope are included (default: `"none"`).                                                                       |
| 4  | [Enable label filter](#enabledisable-labels)                 | If [label enable](../../configuration/container-selection/#enable-label-filter) is set, only containers with the enable label present are included.   |
| 5  | [Disabled containers by label](#disable-containers-by-label) | Containers matching any [disabled label pair](../../configuration/container-selection/#disable-containers-by-label) are excluded.                     |
| 6  | [Enabled containers by label](#enable-containers-by-label)   | If set, only containers matching at least one [enabled label pair](../../configuration/container-selection/#enable-containers-by-label) are included. |
| 7  | [Image skip patterns](#exclude-specific-images)              | Containers whose image matches a [skip pattern](../../configuration/container-selection/#skip-specific-images) are excluded.                          |
| 8  | [Monitored image name patterns](#monitor-specific-images)    | If set, only containers whose image matches a [monitored pattern](../../configuration/container-selection/#monitor-specific-images) are included.     |
| 9  | [Disabled container names](#exclude-specific-containers)     | Containers whose name matches a [disable pattern](../../configuration/container-selection/#disable-specific-containers) are excluded.                 |
| 10 | [Container name arguments](#container-name-filtering)        | If positional name arguments are provided, only containers matching at least one are included.                                                        |

{{< callout type="info" >}}

- If an image name pattern is configured, then Watchtower will exclusively manage only the respective container(s).
- If a container name is provided as an argument, then Watchtower will exclusively manage only the specified container(s).

{{< /callout >}}

{{< callout type="info" >}}
**All criteria must be satisfied**

A container must pass **every** filter in the chain to be monitored. If any single filter rejects it, the container is excluded.

{{< /callout >}}

## Management Modes

Watchtower supports two management modes for containers:

- **Full Management** (default): Watchtower checks for updates, pulls new images, recreates containers, and runs lifecycle hooks.
- [**Monitor Only Mode**](../../configuration/update-behavior/#monitor-only): Watchtower checks for updates, sends notifications, and runs lifecycle hooks, but does **not** recreate containers.

## Container State Filtering

By default, Watchtower only processes containers in the `running` state.

Two configuration options extend this to include additional states:

| Option                                                                                                  | Environment Variable            | Effect                                        |
|---------------------------------------------------------------------------------------------------------|---------------------------------|-----------------------------------------------|
| [Include Stopped Containers](../../configuration/container-selection/#include-stopped-containers)       | `WATCHTOWER_INCLUDE_STOPPED`    | Include `created` and `exited` containers     |
| [Include Restarting Containers](../../configuration/container-selection/#include-restarting-containers) | `WATCHTOWER_INCLUDE_RESTARTING` | Include `restarting` containers (Docker only) |

{{< callout type="info" >}}
**Podman compatibility**

The `restarting` state is not available on Podman and is automatically excluded when Podman is detected.

{{< /callout >}}

## Enable/Disable Labels

The `com.centurylinklabs.watchtower.enable` label controls whether Watchtower manages a container.

{{< callout type="info" >}}
This label is set on the container you want to manage, not on the Watchtower instance.
{{< /callout >}}

### Default Behavior

When [label enable](../../configuration/container-selection/#enable-label-filter) is **not** set:

- Containers **without** the label are **monitored**
- Containers with `enable=true` are **monitored**
- Containers with `enable=false` are **excluded**

### With Label Enable

When [label enable](../../configuration/container-selection/#enable-label-filter) is set:

- Containers with `enable=true` are **monitored**
- Containers with `enable=false` are **excluded**
- Containers **without** the label are **excluded**

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    someimage:
        labels:
            - "com.centurylinklabs.watchtower.enable=true"
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}

```bash
docker run -d \
    --label=com.centurylinklabs.watchtower.enable=true someimage
```

{{< /tab >}}

{{< tab name="Dockerfile" >}}

```dockerfile
LABEL com.centurylinklabs.watchtower.enable="true"
```

{{< /tab >}}

{{< /tabs >}}

To exclude a container, set the label to `false`:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    someimage:
        labels:
            - "com.centurylinklabs.watchtower.enable=false"
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}

```bash
docker run -d \
    --label=com.centurylinklabs.watchtower.enable=false someimage
```

{{< /tab >}}

{{< tab name="Dockerfile" >}}

```dockerfile
LABEL com.centurylinklabs.watchtower.enable="false"
```

{{< /tab >}}

{{< /tabs >}}

## Monitor-Only Mode

Individual containers can be set to [monitor-only mode](../../configuration/update-behavior/#monitor-only), where Watchtower checks for updates and sends notifications but does not recreate the container.

Set the `com.centurylinklabs.watchtower.monitor-only` label to `true` on the container:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    someimage:
        labels:
            - "com.centurylinklabs.watchtower.monitor-only=true"
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}

```bash
docker run -d \
    --label=com.centurylinklabs.watchtower.monitor-only=true someimage
```

{{< /tab >}}

{{< tab name="Dockerfile" >}}

```dockerfile
LABEL com.centurylinklabs.watchtower.monitor-only="true"
```

{{< /tab >}}

{{< /tabs >}}

{{< callout type="info" >}}
The per-container label has the same effect as the global [monitor-only](../../configuration/update-behavior/#monitor-only) option, but applies only to that specific container.

When combined with [label precedence](../../configuration/container-selection/#label-precedence), the container label overrides the global option. Without label precedence, the container is monitor-only if **either** the label or the global option is set.

{{< /callout >}}

## Container Name Filtering

Watchtower can filter containers based on their container name using [Go regex pattern matching](#regex-pattern-matching).

### Include Specific Containers

Pass container names as positional arguments to Watchtower.
When provided, only containers matching at least one name are monitored.

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        command: nginx redis
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}

```bash
docker run -d \
    --name watchtower \
    -v /var/run/docker.sock:/var/run/docker.sock \
    nickfedor/watchtower \
    nginx redis
```

{{< /tab >}}

{{< /tabs >}}

### Exclude Specific Containers

Use the [disable containers](../../configuration/container-selection/#disable-specific-containers) option to exclude containers by name.
This supports comma- or space-separated values and regex patterns.

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_DISABLE_CONTAINERS=container1,container2
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_DISABLE_CONTAINERS="container1,container2" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --disable-containers container1,container2
```

{{< /tab >}}

{{< /tabs >}}

## Label-Based Container Filtering

Watchtower can include or exclude containers based on arbitrary Docker label key-value pairs.

### Enable Containers by Label

Use the [enable containers by label](../../configuration/container-selection/#enable-containers-by-label) option to restrict monitoring to containers that have at least one of the specified label pairs.

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_ENABLE_CONTAINERS_BY_LABEL=env=prod,team=platform
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_ENABLE_CONTAINERS_BY_LABEL="env=prod,team=platform" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --enable-containers-by-label "env=prod,team=platform"
```

{{< /tab >}}

{{< /tabs >}}

An empty value (`key=`) matches any container where the label key exists, regardless of its value.

For example, `WATCHTOWER_ENABLE_CONTAINERS_BY_LABEL=env=` enables any container with the `env` label set to any value.

### Disable Containers by Label

Use the [disable containers by label](../../configuration/container-selection/#disable-containers-by-label) option to exclude containers that have any of the specified label pairs.

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_DISABLE_CONTAINERS_BY_LABEL=managed-by=external
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_DISABLE_CONTAINERS_BY_LABEL="managed-by=external" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --disable-containers-by-label "managed-by=external"
```

{{< /tab >}}

{{< /tabs >}}

{{< callout type="info" >}}
Values containing commas are not supported. Use individual `key=value` pairs separated by commas.

{{< /callout >}}

{{< callout type="info" >}}
A label entry with an empty value (`key=`) performs a presence check: the label must exist on the container with any value. A non-empty value requires an exact match.

{{< /callout >}}

## Image Name Filtering

Watchtower can filter containers based on their image name using [Go regex pattern matching](#regex-pattern-matching).

Image name patterns match against the **full image name including its tag** (e.g., `nginx:latest`, `docker.io/library/nginx:1.25`).

{{< callout type="info" >}}
If no tag is specified in the image reference, `:latest` is assumed.
{{< /callout >}}

### Monitor Specific Images

Use the [monitor image names](../../configuration/container-selection/#monitor-specific-images) configuration option to restrict monitoring to containers whose image matches at least one pattern.

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_MONITOR_IMAGE_NAMES=nginx:.*,redis:7.*
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_MONITOR_IMAGE_NAMES="nginx:.*,redis:7.*" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --monitor-image-names "nginx:.*,redis:7.*"
```

{{< /tab >}}

{{< /tabs >}}

### Exclude Specific Images

Use the [skip image names](../../configuration/container-selection/#skip-specific-images) option to exclude containers whose image matches at least one pattern.

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_SKIP_IMAGE_NAMES=postgres:.*,mcr.microsoft.com/*
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_SKIP_IMAGE_NAMES="postgres:.*,mcr.microsoft.com/*" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --skip-image-names "postgres:.*,mcr.microsoft.com/*"
```

{{< /tab >}}

{{< /tabs >}}

## Monitoring Scopes

Scopes allow multiple Watchtower instances to run on the same Docker host without interfering with each other.
Each instance manages only the containers within its scope.

1. Use the [scope filter](../../configuration/container-selection/#scope-filter) option to define a scope.
2. Then, use the `com.centurylinklabs.watchtower.scope` label on containers to assign them to that scope.

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    # Scoped Application
    app-production:
        image: myapp:latest
        labels:
            - "com.centurylinklabs.watchtower.scope=production"

  # Scoped Watchtower watching "production" scope
    watchtower-production:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_SCOPE=production
        labels:
            - "com.centurylinklabs.watchtower.scope=production"
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}

```bash
docker run -d \
    --name watchtower-production \
    -v /var/run/docker.sock:/var/run/docker.sock \
    nickfedor/watchtower \
    --scope production
```

{{< /tab >}}

{{< /tabs >}}

{{< callout type="info" >}}

- Without a [scope filter](../../configuration/container-selection/#scope-filter), Watchtower defaults to scope `"none"` and manages only **unscoped** containers (those without a `com.centurylinklabs.watchtower.scope` label, or with it set to `"none"` or `""`).
- Containers with a non-empty scope label (e.g., `scope=production`) are **not** monitored by an unscoped Watchtower instance.
- Set the scope filter to `none` to explicitly manage only unscoped containers (same as the default behavior).
- Two instances cannot share the same scope.
- An unscoped instance coexists with scoped instances.

See [Running Multiple Instances](../../advanced-features/running-multiple-instances/) for a complete guide.

{{< /callout >}}

## Regex Pattern Matching

Container name and image name filters support regular expressions using [Go regex syntax](https://pkg.go.dev/regexp/syntax).

{{< callout type="warning" >}}
**Patterns are anchored to match the **full** name**

Patterns are automatically anchored with `^...$`, meaning they must match the entire container or image name. Use `.*` for wildcard matching instead of bare `*`.

{{< /callout >}}

### Container Name Patterns

Container names are normalized before matching (the leading `/` stripped). Both positional arguments and the [disable containers](../../configuration/container-selection/#disable-specific-containers) option support regex.

| Pattern           | Matches                                  |
|-------------------|------------------------------------------|
| `container.*`     | "container1", "container-abc"            |
| `.*-dev`          | "web-dev", "api-dev", "db-dev"           |
| `.*`              | Any container name                       |
| `nginx\|redis`    | Either "nginx" or "redis"                |
| `db-.*\|cache-.*` | Any name starting with "db-" or "cache-" |

### Image Name Patterns

Image name patterns match against the full `name:tag` string. Use `.*` to match any tag.

| Pattern                 | Matches                            |
|-------------------------|------------------------------------|
| `nginx:latest`          | Only `nginx:latest`                |
| `nginx:.*`              | `nginx` with any tag               |
| `docker\.io/library/.*` | Any official Docker Hub image      |
| `.*\.azurecr\.io/.*`    | Any Azure Container Registry image |

### Examples

Exclude all containers starting with a prefix:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_DISABLE_CONTAINERS=web-.*
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_DISABLE_CONTAINERS="web-.*" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --disable-containers "web-.*"
```

{{< /tab >}}

{{< /tabs >}}

Include only containers matching specific patterns:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        command: ["nickfedor/watchtower", "db-.*", "cache-.*"]
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    "db-.*" "cache-.*"
```

{{< /tab >}}

{{< /tabs >}}

Monitor only containers using nginx or redis images with any tag:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_MONITOR_IMAGE_NAMES=nginx:.*,redis:.*
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_MONITOR_IMAGE_NAMES="nginx:.*,redis:.*" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --monitor-image-names "nginx:.*,redis:.*"
```

{{< /tab >}}

{{< /tabs >}}

## Label Precedence

By default, when a container-level label (e.g., `com.centurylinklabs.watchtower.monitor-only`) and a global option (e.g., [monitor-only](../../configuration/update-behavior/#monitor-only)) are both set, the container uses the **combined** effect (either triggers the behavior).

With [label precedence](../../configuration/container-selection/#label-precedence), container labels **override** the global options. This allows per-container control even when global options are set.

| [Label Precedence](../../configuration/container-selection/#label-precedence) | Container Label | Global Option | Result            |
|-------------------------------------------------------------------------------|-----------------|---------------|-------------------|
| false (default)                                                               | not set         | false         | false             |
| false (default)                                                               | not set         | true          | true              |
| false (default)                                                               | true            | false         | true              |
| false (default)                                                               | true            | true          | true              |
| true                                                                          | not set         | any           | global flag value |
| true                                                                          | true            | false         | true              |
| true                                                                          | false           | true          | false             |

This applies to the [`monitor-only`](../../configuration/update-behavior/#monitor-only) and [`no-pull`](../../configuration/update-behavior/#disable-image-pulling) configuration options.

## Complete Configuration Reference

### CLI Flags and Environment Variables

| Flag                            | Environment Variable                     | Type     | Default | Description                           |
|---------------------------------|------------------------------------------|----------|---------|---------------------------------------|
| *(positional args)*             | N/A                                      | []string | []      | Container names/patterns to include   |
| `--disable-containers` / `-x`   | `WATCHTOWER_DISABLE_CONTAINERS`          | []string | []      | Container names/patterns to exclude   |
| `--enable-containers-by-label`  | `WATCHTOWER_ENABLE_CONTAINERS_BY_LABEL`  | []string | []      | Label key=value pairs to include      |
| `--disable-containers-by-label` | `WATCHTOWER_DISABLE_CONTAINERS_BY_LABEL` | []string | []      | Label key=value pairs to exclude      |
| `--monitor-image-names`         | `WATCHTOWER_MONITOR_IMAGE_NAMES`         | []string | []      | Image name patterns to monitor        |
| `--skip-image-names`            | `WATCHTOWER_SKIP_IMAGE_NAMES`            | []string | []      | Image name patterns to exclude        |
| `--label-enable` / `-e`         | `WATCHTOWER_LABEL_ENABLE`                | bool     | false   | Require enable label on containers    |
| `--scope`                       | `WATCHTOWER_SCOPE`                       | string   | ""      | Monitoring scope                      |
| `--include-stopped` / `-S`      | `WATCHTOWER_INCLUDE_STOPPED`             | bool     | false   | Include created and exited containers |
| `--include-restarting`          | `WATCHTOWER_INCLUDE_RESTARTING`          | bool     | false   | Include restarting containers         |
| `--label-take-precedence`       | `WATCHTOWER_LABEL_TAKE_PRECEDENCE`       | bool     | false   | Labels override global flags          |

### Container Labels

| Label                                                                             | Values                | Effect                                          |
|-----------------------------------------------------------------------------------|-----------------------|-------------------------------------------------|
| `com.centurylinklabs.watchtower.enable`                                           | true / false          | Enable or disable management                    |
| `com.centurylinklabs.watchtower.monitor-only`                                     | true / false          | Monitor without updating                        |
| `com.centurylinklabs.watchtower.no-pull`                                          | true / false          | Skip image pulls                                |
| `com.centurylinklabs.watchtower.scope`                                            | any string            | Assign to a monitoring scope                    |
| `com.centurylinklabs.watchtower.depends-on`                                       | comma-separated names | Declare container dependencies                  |
| `com.centurylinklabs.watchtower.cooldown-delay`                                   | duration string       | Minimum image age before updating               |
| [`com.centurylinklabs.watchtower.copy-file`](../../advanced-features/copy-files/) | comma-separated paths | Copy named files onto the replacement container |

## Common Patterns

### Run Multiple Watchtower Instances

Run one instance for production containers and another for development:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower-prod:
        image: nickfedor/watchtower:latest
        command: --scope production --interval 300
        labels:
            - "com.centurylinklabs.watchtower.scope=production"
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock

    watchtower-dev:
        image: nickfedor/watchtower:latest
        command: --scope development --interval 30
        labels:
            - "com.centurylinklabs.watchtower.scope=development"
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}

```bash
docker run -d \
    --name watchtower-production \
    -v /var/run/docker.sock:/var/run/docker.sock \
    nickfedor/watchtower \
    --scope production --interval 300

docker run -d \
    --name watchtower-development \
    -v /var/run/docker.sock:/var/run/docker.sock \
    nickfedor/watchtower \
    --scope development --interval 30
```

{{< /tab >}}

{{< /tabs >}}

### Exclude System Containers

Exclude Watchtower itself and other infrastructure containers:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_DISABLE_CONTAINERS=watchtower,traefik,portainer
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_DISABLE_CONTAINERS="watchtower,traefik,portainer" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --disable-containers "watchtower,traefik,portainer"
```

{{< /tab >}}

{{< /tabs >}}

### Exclude Containers by Label

Exclude containers managed by third-party orchestrators using arbitrary labels:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_DISABLE_CONTAINERS_BY_LABEL=managed-by=external
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_DISABLE_CONTAINERS_BY_LABEL="managed-by=external" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --disable-containers-by-label "managed-by=external"
```

{{< /tab >}}

{{< /tabs >}}

### Monitor Containers by Label

Only monitor containers with specific labels:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_ENABLE_CONTAINERS_BY_LABEL=env=prod,team=platform
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_ENABLE_CONTAINERS_BY_LABEL="env=prod,team=platform" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --enable-containers-by-label "env=prod,team=platform"
```

{{< /tab >}}

{{< /tabs >}}

### Monitor Only Specific Image Registries

Monitor only images from your private registry:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_MONITOR_IMAGE_NAMES=registry.example.com/.*
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
    -e WATCHTOWER_MONITOR_IMAGE_NAMES="registry.example.com/.*" \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
    nickfedor/watchtower \
    --monitor-image-names "registry.example.com/.*"
```

{{< /tab >}}

{{< /tabs >}}

### Selective Monitoring with Enable Label

Use [label enable](../../configuration/container-selection/#enable-label-filter) to explicitly opt containers into monitoring:

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        environment:
            - WATCHTOWER_LABEL_ENABLE=true
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock

    myapp:
        image: myapp:latest
        labels:
            - "com.centurylinklabs.watchtower.enable=true"
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}

```bash
# Start Watchtower with label filtering
docker run -d \
    -e WATCHTOWER_LABEL_ENABLE=true \
    nickfedor/watchtower

# Only this container will be monitored
docker run -d \
    --label com.centurylinklabs.watchtower.enable=true \
    myapp:latest
```

{{< /tab >}}

{{< /tabs >}}
