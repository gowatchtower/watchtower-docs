# Quickstart

## Prerequisites

- Docker installed and running: <https://docs.docker.com/engine/install/>

## Overview

Watchtower is designed to run as a Docker container. It only requires access to the Docker host's Docker socket in order to interact with the Docker engine. By default, Watchtower will review the monitored containers for image updates. If a container has an updated image, then Watchtower will perform the update task by pulling the new image, stopping the old container, rebuilding the container with the new image, and starting the new container.

## Running Watchtower

Watchtower can be run using either a Docker Compose file or via the Docker CLI.

{{< tabs >}}

{{< tab name="Docker Compose" >}}

<ol>
<li>

Obtain the example Docker Compose file:

{{< tabs >}}

{{< tab name="Copy" >}}

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        restart: unless-stopped
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Download via PowerShell (Windows)" >}}

```powershell
iwr -Uri https://raw.githubusercontent.com/nicholas-fedor/watchtower/refs/heads/main/examples/default/docker-compose.yaml -OutFile docker-compose.yaml
```

{{< /tab >}}

{{< tab name="Download via Bash (Linux)" >}}

```bash
curl -L https://raw.githubusercontent.com/nicholas-fedor/watchtower/refs/heads/main/examples/default/docker-compose.yaml -o docker-compose.yaml
```

{{< /tab >}}

{{< /tabs >}}

</li>
<li>

Run the Compose file:

```bash
docker compose up -d
```

[Docker Compose File Reference](https://docs.docker.com/reference/compose-file/)

</li>
</ol>

{{< /tab >}}

{{< tab name="Docker CLI" >}}

Run the following Docker CLI command:

```bash
docker run -d \
    --name watchtower \
    --restart unless-stopped \
    -v /var/run/docker.sock:/var/run/docker.sock \
    nickfedor/watchtower
```

[Docker Run CLI Reference](https://docs.docker.com/reference/cli/docker/container/run/)

{{< /tab >}}

{{< /tabs >}}

## Expected Behavior

When running Watchtower with default settings:

- It will monitor all running containers on the host
- Every 24 hours, it will poll if the monitored containers have updated image digests

If an updated image digest is detected, then Watchtower will:

- Pull the updated container image
- Perform a graceful shutdown of the target container and its dependencies
- Start a new container with the updated image while maintaining the previous container's configuration

