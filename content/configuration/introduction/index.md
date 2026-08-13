---
title: Introduction
weight: 1
---
## Deprecation Notice

{{< callout type="warning" >}}
**Watchtower v2 Legacy Notification Deprecation**

**Watchtower has a number of legacy notification configuration options that will be removed with the release of Watchtower v2:**

- [Email Notifications](../../notifications/deprecations/email/)
- [Gotify Notifications](../../notifications/deprecations/gotify/)
- [Microsoft Teams Notifications](../../notifications/deprecations/teams/)
- [Slack Notifications](../../notifications/deprecations/slack/)

Migration to the [`NOTIFICATION URL`](../notifications/#notification_url) with the appropriate Shoutrrr URL scheme is strongly recommended.

Use Watchtower's CLI [`migration tool`](../../notifications/deprecations/migration-tool/) to help convert legacy email configurations to Shoutrrr URLs or use the [Shoutrrr Playground](https://shoutrrr.nickfedor.com/latest/playground/) to help convert configurations for other services to Shoutrrr URLs.

{{< /callout >}}

## Overview

By default, Watchtower monitors all containers running on the Docker daemon it connects to (typically the local daemon, configurable via the `--host` flag).

To limit monitoring to specific containers, provide their names as arguments when starting Watchtower.

{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
    watchtower:
        command: nginx redis
        image: nickfedor/watchtower:latest
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}
```bash
docker run -d \
    --name watchtower \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --restart unless-stopped \
    nickfedor/watchtower \
    nginx redis
```

{{< /tab >}}

{{< /tabs >}}

- In the above example, Watchtower monitors only the `nginx` and `redis` containers, ignoring others.

To run a single update attempt and exit, use the [`run-once`](../scheduling/#run_once) configuration option.

{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
    watchtower:
        command: nginx redis
        environment:
            WATCHTOWER_RUN_ONCE: true
        image: nickfedor/watchtower:latest
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}
```bash
docker run --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    nickfedor/watchtower \
    --run-once \
    nginx redis
```

{{< /tab >}}

{{< /tabs >}}

- This command triggers an update attempt for "nginx" and "redis" containers, displays debug output, and removes the Watchtower container upon completion.

{{< callout type="info" >}}
Regex patterns are supported. See [Regex Pattern Matching](../../getting-started/container-selection/#regex_pattern_matching) for details.

{{< /callout >}}

## General Options

### Help

Displays documentation for supported flags.

```text
            Argument: --help
Environment Variable: N/A
                Type: N/A
             Default: N/A
```

### Health Check

Returns a success exit code for Docker `HEALTHCHECK`, verifying another process is running in the container.

```text
            Argument: --health-check
Environment Variable: None
                Type: N/A
             Default: N/A
```

{{< callout type="info" >}}
Intended solely for Docker `HEALTHCHECK`.
Do not use on the main command line.

{{< /callout >}}
