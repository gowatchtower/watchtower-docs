---
title: Linked Containers
weight: 4
---
## Overview

Watchtower, by default, ensures that interdependent containers are updated in the correct order to maintain application stability.
It automatically detects container dependencies through various mechanisms and uses topological sorting to determine the optimal update sequence.

When containers depend-upon each other (such as a web application depending on a database), updating them in the wrong order can cause service disruptions.

Watchtower addresses this by:

- **Detecting dependencies** through Docker links, labels, and network configurations
- **Performing topological sorting** to determine the correct update order
- **Stopping containers in reverse dependency order** (dependents first - i.e., web app before database)
- **Restarting containers in dependency order** (dependencies first - i.e., database before web app)

This ensures that dependent services are stopped before their dependencies are updated, and dependencies are available when dependents restart.

## How It Works

### Dependency Detection

Watchtower detects container dependencies from labels and from each container's Docker configuration.

Labels are checked in the following priority order, and only the first one present is used:

1. **Watchtower depends-on label** (`com.centurylinklabs.watchtower.depends-on`)
2. **Docker Compose depends_on label** (`com.docker.compose.depends_on`)

An empty Watchtower depends-on label counts as absent, so the Docker Compose depends_on label is used instead when Compose detection is enabled.

The Docker configuration always counts as well, whichever label is present:

- **Docker links and network mode** (legacy Docker linking, `network_mode: service:SERVICE` for a Compose service, and `network_mode: container:CONTAINER` for a container name or ID)
- **Volumes from another container** (`volumes_from` in Compose, `--volumes-from` in the Docker CLI)

A container named by both a label and the Docker configuration is a single dependency.

#### Watchtower Depends-On Label

The `com.centurylinklabs.watchtower.depends-on` label allows explicit declaration of dependencies:

```dockerfile
LABEL com.centurylinklabs.watchtower.depends-on="database,redis"
```

This label accepts a comma-separated list of container names that must be available before this container starts.
This supports referencing containers from other Docker Compose projects or stacks, enabling cross-project dependency management.

An entry may also name a Docker Compose service, such as `db` or `my-database`, or a project-qualified service, such as `myproject-db`.
A service name resolves to every replica of that service.

A service name used in more than one Compose project, such as `db`, resolves to the service in the labeled container's own project when that project has one.
Otherwise it resolves to the service in every project that has it, so the container restarts when any of them is updated.
When the service in its own project is unmonitored, the entry is left alone instead (see [Dependencies on Unmonitored Containers](#dependencies-on-unmonitored-containers)).
Use the container name or a project-qualified service name, such as `myproject-db`, to depend on only one of them.

{{< callout type="warning" >}}
An entry that names no container is ignored and logged as a warning.

{{< /callout >}}

#### Docker Compose Depends-On Label

Watchtower automatically recognizes Docker Compose's `depends_on` relationships:

```yaml
services:
  web:
    image: nginx
    depends_on:
      - database
  database:
    image: postgres
```

The `com.docker.compose.depends_on` label is automatically set by Docker Compose and parsed by Watchtower to extract service names and support implicitly restarting linked services.

<!-- markdownlint-disable MD046 -->
{{< callout type="warning" >}}
**Docker Compose depends_on Restart Behavior**

Docker Compose's [`depends_on`](https://docs.docker.com/reference/compose-file/services/#depends_on) has a `restart` attribute in the long-form syntax:

```yaml
services:
  web:
    depends_on:
      database:
        condition: service_healthy
        restart: true
```

Watchtower does not support using this to control implicit restarts, because this is an explicit, opt-in feature that defaults to `false` when omitted, such as when using the short-form syntax.

{{< /callout >}}

<!-- markdownlint-enable MD046 -->

#### Docker Links and Network Mode

For legacy Docker links, and for containers that use another container's network (`network_mode: service:SERVICE` or `network_mode: container:CONTAINER` in Compose, `--network container:CONTAINER` in the Docker CLI), Watchtower treats these as implicit dependencies:

```dockerfile
# Container with explicit link
docker run --link database:db nginx

# Container using another container's network
docker run --network container:CONTAINER nginx
```

Replace *`CONTAINER`* with the name or ID of the container whose network to use.

These count as dependencies even when a depends-on or Compose `depends_on` label names other containers, because the container cannot run without the container whose network it uses.
The same applies to `volumes_from` (`--volumes-from`): the container is always ordered after the container whose volumes it uses.

### Topological Sorting

Watchtower uses topological sorting to determine the correct update order. This algorithm:

- Builds a dependency graph from all detected relationships
- Detects cycles and sets aside every container in one
- Produces a linear ordering where dependencies precede dependents

{{< callout type="warning" >}}
Containers in a circular dependency are skipped and reported as skipped, while the rest of the update continues.
Ensure your dependency graph is acyclic so those containers are updated.

{{< /callout >}}

### Update Order

When updates are needed, Watchtower follows this sequence:

1. **Identify all containers requiring updates**
2. **Expand the set** to include all containers in the dependency chain
3. **Sort containers** using topological order (dependencies first)
4. **Stop containers** in reverse topological order (dependents first)
5. **Update and restart containers** in topological order (dependencies first)

This ensures that:

- Dependent services are stopped before their dependencies change
- Dependencies are fully restarted before dependents attempt to connect

### Restarting Dependent Containers

When a container is updated, Watchtower also restarts the containers that depend on it, directly or through other restarted containers, without updating their images.
Dependents are found with the same dependency graph that orders the update, so the containers restarted and the update order always agree.

- A Watchtower container restarts when a container it depends on restarts, such as a VPN container it shares a network with.
- A container that names a Watchtower container exactly, by container name or ID, restarts when that Watchtower container is updated. The name without its replica number (such as `infra-watchtower` for `infra-watchtower-1`) also works when exactly one Watchtower container has it, so it names nothing when the service runs several replicas.
- Containers in a circular dependency are neither updated nor restarted.

### Dependencies on Unmonitored Containers

A container that Watchtower does not monitor is left alone, including when a monitored container depends on it.
Containers are excluded from monitoring by the [container selection](../../configuration/container-selection/) options, such as the enable label, the disable lists, image filters, and scopes.

- A dependency on an unmonitored container is ignored when ordering updates and restarting dependents.
- A dependency that names an unmonitored container exactly never resolves to a different, monitored container with a similar name.
- A service name whose service in the container's own Compose project is unmonitored never resolves to the same service in another project.

## Configuration

### Automatic Detection

In most cases, no additional configuration is required. Watchtower automatically detects dependencies from:

- Docker Compose `depends_on` declarations
- Existing Docker links
- `network_mode: service:SERVICE` and `network_mode: container:CONTAINER` configurations
- `volumes_from` configurations

{{< callout type="info" >}}
**Docker Compose Considerations**

When using Docker Compose, Watchtower leverages `depends_on` declarations for dependency detection. Dependencies are resolved using service names, not container names. Ensure your `depends_on` references service names correctly.

{{< /callout >}}

{{< callout type="warning" >}}
**Rolling restart is currently not supported when used in combination with linked-containers.**

This limitation exists because linked-containers require coordinated updates across dependency chains, which conflicts with the incremental nature of rolling restarts.

{{< /callout >}}

### Disable Docker Compose Depends-On

If you want to disable automatic dependency detection from the Docker Compose `depends_on` configuration while preserving other dependency sources, use the following configuration:

<!-- markdownlint-disable MD046 -->
{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  watchtower:
    image: nickfedor/watchtower
    environment:
      - WATCHTOWER_USE_COMPOSE_DEPENDS_ON=false
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}
```bash
docker run -d \
    --name watchtower \
    -v /var/run/docker.sock:/var/run/docker.sock \
    nickfedor/watchtower \
    --use-compose-depends-on=false
```

{{< /tab >}}

{{< /tabs >}}

<!-- markdownlint-enable MD046 -->

This disables parsing of the `com.docker.compose.depends_on` label while still honoring:

- Watchtower's explicit `com.centurylinklabs.watchtower.depends-on` label
- Legacy Docker links
- Network mode dependencies (`network_mode: service:SERVICE` and `network_mode: container:CONTAINER`)
- Volumes from another container (`volumes_from`)

### Explicit Dependencies

For cases where automatic detection is insufficient, use the Watchtower depends-on label:

<!-- markdownlint-disable MD046 -->
{{< tabs >}}

{{< tab name="Dockerfile" >}}
```dockerfile
FROM nginx:latest

# Declare dependencies on database and cache services
LABEL com.centurylinklabs.watchtower.depends-on="postgres,redis"
```

{{< /tab >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  web:
    image: nginx
    labels:
      - com.centurylinklabs.watchtower.depends-on=postgres,redis
  postgres:
    image: postgres
  redis:
    image: redis
```

{{< /tab >}}

{{< /tabs >}}

<!-- markdownlint-enable MD046 -->

### Advanced Scenarios

#### Multiple Dependencies

```dockerfile
LABEL com.centurylinklabs.watchtower.depends-on="database,cache,queue"
```

#### Complex Dependency Chains

```yaml
services:
  web:
    image: nginx
    depends_on:
      - api
  api:
    image: myapi
    depends_on:
      - database
      - cache
  database:
    image: postgres
  cache:
    image: redis
```

In this scenario, Watchtower will update containers in the order: `cache`, `database`, `api`, `web`.

## Examples

### MySQL-WordPress Scenario

Consider a classic WordPress setup with MySQL database:

<!-- markdownlint-disable MD046 -->
{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  wordpress:
    image: wordpress:latest
    depends_on:
      - mysql
    environment:
      WORDPRESS_DB_HOST: mysql
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: 'DB_ROOT_PASSWORD'
```

Replace *`DB_ROOT_PASSWORD`* with a root password for the database, inside the single quotes so that YAML reads every character as part of it. Inside the quotes, write each `'` as `''`. Write each `$` as `$$`, because Compose reads a `$` as the start of a variable even inside quotes.

{{< /tab >}}

{{< tab name="Update Process" >}}
When Watchtower detects a MySQL update:

1. **Dependency Detection**: Identifies that `wordpress` depends on `mysql`
2. **Stop Order**: Stops `wordpress` first, then `mysql`
3. **Update Order**: Updates and restarts `mysql` first, then `wordpress`

This prevents WordPress from losing database connectivity during the update.

{{< /tab >}}

{{< /tabs >}}

<!-- markdownlint-enable MD046 -->

### Microservices Architecture

For complex applications with multiple services:

<!-- markdownlint-disable MD046 -->
{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  api-gateway:
    image: nginx
    depends_on:
      - auth-service
      - user-service
  auth-service:
    image: auth-service
    depends_on:
      - redis
  user-service:
    image: user-service
    depends_on:
      - postgres
  redis:
    image: redis
  postgres:
    image: postgres
```

{{< /tab >}}

{{< tab name="Dependency Graph" >}}
```
postgres → user-service
redis → auth-service
auth-service → api-gateway
user-service → api-gateway
```

Update order: `postgres`, `redis`, `user-service`, `auth-service`, `api-gateway`

{{< /tab >}}

{{< /tabs >}}

<!-- markdownlint-enable MD046 -->

### Legacy Docker Links

For applications using traditional Docker linking:

<!-- markdownlint-disable MD046 -->
{{< tabs >}}

{{< tab name="Docker Run Commands" >}}
```bash
# Start database
docker run -d --name mysql mysql:8.0

# Start web app with link
docker run -d --name webapp --link mysql:db nginx
```

{{< /tab >}}

{{< tab name="Watchtower Behavior" >}}
Watchtower automatically detects the link and ensures `webapp` is stopped before `mysql` updates, and `mysql` restarts before `webapp`.

{{< /tab >}}

{{< /tabs >}}

### Network Mode Dependencies

Containers using `network_mode: service:SERVICE`:

{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  sidecar:
    image: sidecar
    network_mode: service:main-app
  main-app:
    image: main-app
```

{{< /tab >}}

{{< tab name="Watchtower Behavior" >}}
The `sidecar` container is treated as dependent on `main-app`, ensuring proper update sequencing.

{{< /tab >}}

{{< /tabs >}}

<!-- markdownlint-enable MD046 -->

## Troubleshooting

### Common Issues

Watchtower reports each of the following dependency problems once per container.
The first scan that finds a problem logs and notifies it as a warning.
Later scans log it at debug level.
Every scan still skips a container that depends on itself or sits in a circular dependency, and lists it as skipped in the session report.
A container depends on itself when one of its dependencies names it by its container name or by its own Compose service.
A skipped container is not restarted when another of its dependencies is updated.
A depends-on entry that names no container is only ignored, and its container is still checked and updated.
A recreated container or a restart of Watchtower reports the problem again.

#### Containers Skipped Due to Circular Dependencies

{{< callout type="warning" >}}
If you see "Skipping container update (circular dependency)" warnings, check your dependency declarations for cycles.

{{< /callout >}}

**Solution**: Review and remove circular dependencies. For example, if A depends on B and B depends on A, remove one of the dependencies or restructure your services.

#### Ignored Depends-On Entries

{{< callout type="warning" >}}
If you see "Ignoring depends-on entry that names no container" warnings, the named `com.centurylinklabs.watchtower.depends-on` entry matched no container.

{{< /callout >}}

**Solution**: Use the exact container name, a Compose service name, or a project-qualified service name such as `myproject-db`.

#### Containers Not Updating in Expected Order

**Check**:

- Verify dependency labels are correctly formatted
- Ensure container names match exactly
- Check Docker Compose service names vs container names

#### Missing Dependencies

**Symptoms**: Containers update out of order or fail to connect after updates.

**Debug**: Enable debug logging to see detected dependencies:

```bash
watchtower --debug
```

Look for log messages like:

- "Found container dependencies", which lists each container's dependencies once per scan

### Debugging Commands

Enable verbose logging to inspect dependency detection:

```bash
docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  nickfedor/watchtower \
  --debug \
  --run-once
```

Check container labels:

```bash
docker inspect container_name | jq '.[0].Config.Labels'
```
