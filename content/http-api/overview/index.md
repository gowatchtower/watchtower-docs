---
title: Overview
weight: 1
---

Watchtower has an optional HTTP API server that is enabled by configuring the  [`http-api-endpoints`](../../configuration/http-api/#http_api_endpoints) and [`http-api-token`](../../configuration/http-api/#http_api_token) configuration options and [publishing](https://docs.docker.com/engine/network/port-publishing/) the HTTP API port (default: 8080).

This HTTP API server extends Watchtower's functionality by providing HTTP endpoints that allows for programmatic access to functionality, such as checking the current status of monitored containers and requesting Watchtower to perform container updates.
Parameters, such as container and image name filters, can be used to further fine-tune HTTP requests.

{{< callout type="warning" >}}
**Deprecation of HTTP API Endpoint Configuration Options**

The following options have been deprecated in favor of the unified [`http-api-endpoints`](../../configuration/http-api/#http_api_endpoints) configuration option and will be removed with the v2.x.x release of Watchtower:

- [`http-api-update`](../../configuration/http-api/#http_api_update)
- [`http-api-containers`](../../configuration/http-api/#http_api_containers)
- [`http-api-metrics`](../../configuration/http-api/#http_api_metrics)

{{< /callout >}}

## Endpoints

The following endpoints can be enabled by using the[`http-api-endpoints`](../../configuration/http-api/#http_api_endpoints) configuration option and the respective configuration value.

|                           **Name**                           | **Configuration Value** | **Method** |       **Endpoint**       |                                  **Auth**                                   |                                                                                **Parameters**                                                                                |                                              **Description**                                              |
|:------------------------------------------------------------:|:-----------------------:|:----------:|:------------------------:|:---------------------------------------------------------------------------:|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:---------------------------------------------------------------------------------------------------------:|
|            [Update](../endpoints/update/)            |        `update`         |   `POST`   |       `/v1/update`       |      [API token](../../configuration/http-api/#http_api_token)      | [`image`](../endpoints/update/#image_name), [`container`](../endpoints/update/#container_name), [`async`](../endpoints/update/#asynchronous_updates) |                   Triggers container updates and returns JSON results of the operation                    |
|             [Check](../endpoints/check/)             |         `check`         |   `POST`   |       `/v1/check`        |      [API token](../../configuration/http-api/#http_api_token)      |                                 [`image`](../endpoints/check/#image_name), [`container`](../endpoints/check/#container_name)                                 |                     Checks containers for available updates via registry digest query                     |
|        [Containers](../endpoints/containers/)        |      `containers`       |   `GET`    |     `/v1/containers`     |      [API token](../../configuration/http-api/#http_api_token)      |                              [`name`](../endpoints/containers/#container_name), [`image`](../endpoints/containers/#image_name)                               |                     Lists watched containers and their current running image digests                      |
| [Container Details](../endpoints/container-details/) |      `containers`       |   `GET`    | `/v1/containers/details` |      [API token](../../configuration/http-api/#http_api_token)      |                       [`name`](../endpoints/container-details/#container_name), [`image`](../endpoints/container-details/#image_name)                        | Returns detailed information about each watched container including running state and configuration flags |
|           [History](../endpoints/history/)           |        `history`        |   `GET`    |      `/v1/history`       |      [API token](../../configuration/http-api/#http_api_token)      |                [`since`](../endpoints/history/#since), [`until`](../endpoints/history/#until), [`limit`](../endpoints/history/#limit)                |            Returns historical scan results from the in-memory ring buffer (up to 500 entries)             |
|            [Images](../endpoints/images/)            |        `images`         |   `GET`    |       `/v1/images`       |      [API token](../../configuration/http-api/#http_api_token)      |                                       [`name`](../endpoints/images/#image_name), [`id`](../endpoints/images/#image_id)                                       |                   Lists tracked images with their current digests and container counts                    |
|            [Config](../endpoints/config/)            |        `config`         |   `GET`    |       `/v1/config`       |      [API token](../../configuration/http-api/#http_api_token)      |                                                                                                                                                                              |                           Returns the active Watchtower configuration settings                            |
|            [Events](../endpoints/events/)            |        `events`         |   `GET`    |       `/v1/events`       | [Events token](../../configuration/http-api/#http_api_events_token) |                                                                                                                                                                              |                        Streams real-time operational events via Server-Sent Events                        |
|            [Status](../endpoints/status/)            |        `metrics`        |   `GET`    |       `/v1/status`       |      [API token](../../configuration/http-api/#http_api_token)      |                                                                                                                                                                              |                                Returns the summary of the most recent scan                                |
|           [Metrics](../endpoints/metrics/)           |        `metrics`        |   `GET`    |      `/v1/metrics`       |      [API token](../../configuration/http-api/#http_api_token)      |                                                                                                                                                                              |                     Exposes Prometheus-compatible metrics for monitoring and alerting                     |
|           [Swagger](../endpoints/swagger/)           |        `swagger`        |   `GET`    |       `/swagger/*`       |                                    None                                     |                                                                                                                                                                              |                               Interactive API documentation via Swagger UI                                |
|           [Liveness](../endpoints/health/)           |        `health`         |   `GET`    |         `/livez`         |                                    None                                     |                                                                                                                                                                              |                                Returns `200 OK` when the server is running                                |
|          [Readiness](../endpoints/health/)           |        `health`         |   `GET`    |        `/readyz`         |                                    None                                     |                                                                                                                                                                              |                     Returns `200 OK` when Docker client is connected, `503` otherwise                     |
|           [Startup](../endpoints/health/)            |        `health`         |   `GET`    |       `/startupz`        |                                    None                                     |                                                                                                                                                                              |                               Returns `200 OK` once the server has started                                |

{{< callout type="info" >}}
- Endpoints enforce HTTP method restrictions using method-based routing.
- Requests with unsupported methods will receive a `405 Method Not Allowed` response.
- Watchtower will not start if [`http-api-endpoints`](../../configuration/http-api/#http_api_endpoints) has incorrect values or is combined with `all` (which enables all endpoints).

{{< /callout >}}

## Examples

{{< tabs >}}

{{< tab name="Enable Check Endpoint" >}}
{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
        environment:
            - WATCHTOWER_HTTP_API_ENDPOINTS=check
            - WATCHTOWER_HTTP_API_TOKEN=your-secure-token
        ports:
            - 8080:8080
        restart: unless-stopped
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}
```bash
docker run -d \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -e WATCHTOWER_HTTP_API_ENDPOINTS=check \
    -e WATCHTOWER_HTTP_API_TOKEN=your-secure-token \
    -p 8080:8080 \
    --restart unless-stopped \
    nickfedor/watchtower
```

{{< /tab >}}

{{< /tabs >}}

{{< /tab >}}

{{< tab name="Enable Update Endpoint" >}}
{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
        environment:
            - WATCHTOWER_HTTP_API_ENDPOINTS=update
            - WATCHTOWER_HTTP_API_TOKEN=your-secure-token
        ports:
            - 8080:8080
        restart: unless-stopped
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}
```bash
docker run -d \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -e WATCHTOWER_HTTP_API_ENDPOINTS=update \
    -e WATCHTOWER_HTTP_API_TOKEN=your-secure-token \
    -p 8080:8080 \
    --restart unless-stopped \
    nickfedor/watchtower
```

{{< /tab >}}

{{< /tabs >}}

{{< /tab >}}

{{< tab name="Enable Swagger Endpoint" >}}
{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
        environment:
            - WATCHTOWER_HTTP_API_ENDPOINTS=swagger
            - WATCHTOWER_HTTP_API_TOKEN=your-secure-token
        ports:
            - 8080:8080
        restart: unless-stopped
```

{{< /tab >}}

{{< tab name="Docker CLI" >}}
```bash
docker run -d \
    --name watchtower \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -e WATCHTOWER_HTTP_API_ENDPOINTS=swagger \
    -e WATCHTOWER_HTTP_API_TOKEN=your-secure-token \
    -p 8080:8080 \
    --restart unless-stopped \
    nickfedor/watchtower
```

{{< /tab >}}

{{< /tabs >}}

{{< /tab >}}

{{< /tabs >}}

## Security considerations

Watchtower was originally designed to be a stateless application with direct access to the Docker socket.
By having direct access to the Docker socket, the application effectively runs with root-level privileges.
It is worth taking a moment to review and familiarize yourself with the [Docker Engine security documentation](https://docs.docker.com/engine/security/).

- !!! Warning "The HTTP API should never be directly exposed to the Internet!"
- !!! Warning "Using the HTTP API without TLS encryption is insecure and not recommended!"
- !!! Warning "Only enable the endpoints that you need!"

Authentication, TLS, trusted proxies, CORS, and rate limiting are covered under [Authentication](../configuration/authentication/), [TLS](../configuration/tls/), and the [HTTP API configuration reference](../../configuration/http-api/).
