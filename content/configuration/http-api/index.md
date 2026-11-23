---
title: HTTP API
weight: 12
---
## Deprecation Notice

{{< callout type="warning" >}}
**Watchtower v2 Legacy HTTP API Configuration Deprecation**

Endpoint-specific configuration options will be removed with the release of Watchtower v2.

Use the [`HTTP API Endpoints`](#http-api-endpoints) configuration option instead.

{{< /callout >}}

## HTTP API Host

Sets the host interface for binding the HTTP API.

```text
            Argument: --http-api-host
Environment Variable: WATCHTOWER_HTTP_API_HOST
                Type: String
             Default: empty (binds to all IPv4 interfaces)
```

{{< callout type="info" >}}
See the [HTTP API Host documentation](../../http-api/configuration/host-and-port/#http-api-host) for details
{{< /callout >}}

## HTTP API Port

Sets the listening port for the HTTP API.

```text
            Argument: --http-api-port
Environment Variable: WATCHTOWER_HTTP_API_PORT
                Type: String
             Default: 8080
```

{{< callout type="info" >}}
See the [HTTP API Port documentation](../../http-api/configuration/host-and-port/#http-api-port) for details
{{< /callout >}}

## HTTP API Token

Sets an authentication token for HTTP API requests.

```text
            Argument: --http-api-token
Environment Variable: WATCHTOWER_HTTP_API_TOKEN
                Type: String
             Default: None
```

{{< callout type="info" >}}
Supports file path for Docker Secrets (e.g., `/run/secrets/http_api_token`).

{{< /callout >}}

## HTTP API Events Token

Sets an authentication token with read-only permissions specifically for the events SSE endpoint (`/v1/events`).

```text
            Argument: --http-api-events-token
Environment Variable: WATCHTOWER_HTTP_API_EVENTS_TOKEN
                Type: String
             Default: None
```

This is a separate token from the [`http-api-token`](#http-api-token) in order to reduce credential exposure, because events tokens may be exposed in URL query parameters.

The token can be provided via the `Authorization` header (for programmatic clients) or as the `access_token` query parameter (for browser `EventSource` which cannot set custom headers).

{{< callout type="warning" >}}
**Events token leakage risk**

When using the `access_token` query parameter (required for browser `EventSource`), the token appears in browser history, proxy logs, and access logs.
This is a lower exposure than the main API token, but you should still treat the events token as a secret and rotate it if it has been exposed.
Prefer header-based auth when possible.

{{< /callout >}}

{{< callout type="warning" >}}
**This is **required** when the `events` endpoint is enabled via [`http-api-endpoints`](#http-api-endpoints).**

{{< /callout >}}

{{< callout type="info" >}}
Supports file path for Docker Secrets (e.g., `/run/secrets/http_api_events_token`).

{{< /callout >}}

## HTTP API Endpoints

Selects which HTTP API endpoints to enable.

```text
            Argument: --http-api-endpoints
Environment Variable: WATCHTOWER_HTTP_API_ENDPOINTS
                Type: String (comma or whitespace separated list)
             Default: empty (HTTP API disabled)
```

{{< callout type="info" >}}
For more information, review the [HTTP API](../../http-api/overview/) documentation
{{< /callout >}}

Valid names (case-insensitive):

| Name                                                         | Routes                                                                                                | Auth                                              |
|--------------------------------------------------------------|-------------------------------------------------------------------------------------------------------|---------------------------------------------------|
| [`health`](../../http-api/endpoints/health/)         | `/livez`, `/readyz`, `/startupz`                                                                      | None                                              |
| [`update`](../../http-api/endpoints/update/)         | `POST /v1/update`                                                                                     | [`http-api-token`](#http-api-token)               |
| [`metrics`](../../http-api/endpoints/metrics/)       | `GET /v1/metrics`, [`GET /v1/status`](../../http-api/endpoints/status/)                       | [`http-api-token`](#http-api-token)               |
| [`containers`](../../http-api/endpoints/containers/) | `GET /v1/containers`, [`/v1/containers/details`](../../http-api/endpoints/container-details/) | [`http-api-token`](#http-api-token)               |
| [`check`](../../http-api/endpoints/check/)           | `POST /v1/check`                                                                                      | [`http-api-token`](#http-api-token)               |
| [`history`](../../http-api/endpoints/history/)       | `GET /v1/history`                                                                                     | [`http-api-token`](#http-api-token)               |
| [`images`](../../http-api/endpoints/images/)         | `GET /v1/images`                                                                                      | [`http-api-token`](#http-api-token)               |
| [`config`](../../http-api/endpoints/config/)         | `GET /v1/config`                                                                                      | [`http-api-token`](#http-api-token)               |
| [`events`](../../http-api/endpoints/events/)         | `GET /v1/events`                                                                                      | [`http-api-events-token`](#http-api-events-token) |
| [`swagger`](../../http-api/endpoints/swagger/)       | `GET /swagger/*`                                                                                      | None                                              |

{{< callout type="warning" >}}
- Protected `/v1/*` endpoints require configuring a [`http-api-token`](#http-api-token).
- Enabling `events` requires configuring a [`http-api-events-token`](#http-api-events-token).
- The `all` value enables every endpoint and **MUST** be the only value.
- Watchtower will exit on startup if an incorrect name is used or combined with `all`.

{{< /callout >}}

{{< callout type="info" >}}
**Flag and environment value usage**

- The CLI flag can be specified multiple times.
- Defining the environment variable multiple times will not work (only the last value is used).
- For environment variables use comma or space-separated values.

{{< /callout >}}

## HTTP API Rate Limit

Sets the maximum number of API requests allowed per minute per IP address. Watchtower exits at startup if the value is less than 1.

```text
            Argument: --http-api-rate-limit
Environment Variable: WATCHTOWER_HTTP_API_RATE_LIMIT
                Type: Integer
             Default: 60
```

{{< callout type="info" >}}
When the limit is exceeded, the client receives HTTP 429 (Too Many Requests).

{{< /callout >}}

## HTTP API Check Timeout

Sets the maximum duration for the `/v1/check` API endpoint. Watchtower exits at startup if the value is negative.

```text
            Argument: --http-api-check-timeout
Environment Variable: WATCHTOWER_HTTP_API_CHECK_TIMEOUT
                Type: Duration
             Default: 5m
```

{{< callout type="info" >}}
The `/v1/check` endpoint queries registries for the latest image digest.
Increase this timeout if you see `context deadline exceeded` errors during checks.

{{< /callout >}}

## HTTP API Update Timeout

Sets the maximum duration for the `/v1/update` API endpoint. Watchtower exits at startup if the value is negative.

```text
            Argument: --http-api-update-timeout
Environment Variable: WATCHTOWER_HTTP_API_UPDATE_TIMEOUT
                Type: Duration
             Default: 10m
```

{{< callout type="info" >}}
The `/v1/update` endpoint performs a full container update scan.
Increase this timeout if you see `context deadline exceeded` errors during updates.

{{< /callout >}}

## HTTP API TLS Certificate

Path to the TLS certificate file for the HTTP API.

```text
            Argument: --http-api-tls-cert
Environment Variable: WATCHTOWER_HTTP_API_TLS_CERT
                Type: String
             Default: empty (HTTP)
```

{{< callout type="warning" >}}
When both the [`http-api-tls-cert`](#http-api-tls-certificate) and [`http-api-tls-key`](#http-api-tls-key) are provided, the server uses HTTPS.

{{< /callout >}}

## HTTP API TLS Key

Path to the TLS private key file for the HTTP API.

```text
            Argument: --http-api-tls-key
Environment Variable: WATCHTOWER_HTTP_API_TLS_KEY
                Type: String
             Default: empty (HTTP)
```

{{< callout type="warning" >}}
When both the [`http-api-tls-cert`](#http-api-tls-certificate) and [`http-api-tls-key`](#http-api-tls-key) are provided, the server uses HTTPS.

{{< /callout >}}

{{< callout type="info" >}}
See the [HTTP API Host documentation](../../http-api/configuration/tls/) for details
{{< /callout >}}

## HTTP API Trusted Proxies

Comma-separated list of trusted proxy IP addresses or CIDR ranges for reverse proxy support.
When set, enables proxy header processing for client IP and scheme detection.

```text
            Argument: --http-api-trusted-proxies
Environment Variable: WATCHTOWER_HTTP_API_TRUSTED_PROXIES
                Type: String Array
             Default: (unset)
```

{{< callout type="warning" >}}
Required for correct client IP detection and rate limiting behind a reverse proxy (Traefik, Caddy, Nginx, etc.).

{{< /callout >}}

## HTTP API Proxy Header

Header to read the real client IP from when behind a reverse proxy.

```text
            Argument: --http-api-proxy-header
Environment Variable: WATCHTOWER_HTTP_API_PROXY_HEADER
                Type: String
             Default: X-Forwarded-For
```

{{< callout type="info" >}}
Only used when [`http-api-trusted-proxies`](#http-api-trusted-proxies) is set.

{{< /callout >}}

Common values:

- `X-Forwarded-For` (Traefik, Caddy)
- `X-Real-IP` (Nginx)
- `CF-Connecting-IP` (Cloudflare)

## HTTP API CORS Origins

Comma-separated list of allowed CORS origins for cross-origin requests.
Use this to restrict which origins can access the API.

```text
            Argument: --http-api-cors-origins
Environment Variable: WATCHTOWER_HTTP_API_CORS_ORIGINS
                Type: String Array
             Default: Unset (same-origin only)
```

{{< callout type="info" >}}
When unset, CORS is disabled and only same-origin requests are allowed.
Set this to permit specific cross-origin origins.

{{< /callout >}}

## Deprecated Configuration Options

/// details | The following legacy configuration options are deprecated and will be removed with the release of Watchtower v2.
    type: warning

### HTTP API Containers

Runs Watchtower with the read-only containers API enabled.

```text
            Argument: --http-api-containers
Environment Variable: WATCHTOWER_HTTP_API_CONTAINERS
                Type: Boolean
             Default: false
```

{{< callout type="warning" >}}
Add `containers` to the  [`http-api-endpoints`](#http-api-endpoints) configuration instead.

For example: `WATCHTOWER_HTTP_API_ENDPOINTS=containers`

{{< /callout >}}

### HTTP API Metrics

Runs Watchtower with the Prometheus metrics API enabled.

```text
            Argument: --http-api-metrics
Environment Variable: WATCHTOWER_HTTP_API_METRICS
                Type: Boolean
             Default: false
```

{{< callout type="warning" >}}
Add `metrics` to the  [`http-api-endpoints`](#http-api-endpoints) configuration instead.

For example: `WATCHTOWER_HTTP_API_ENDPOINTS=metrics`

{{< /callout >}}

### HTTP API Update

Runs Watchtower in HTTP API mode, so that image updates must be triggered by a request.

```text
            Argument: --http-api-update
Environment Variable: WATCHTOWER_HTTP_API_UPDATE
                Type: Boolean
             Default: false
```

{{< callout type="warning" >}}
Add `update` to the  [`http-api-endpoints`](#http-api-endpoints) configuration instead.

For example: `WATCHTOWER_HTTP_API_ENDPOINTS=update`

{{< /callout >}}

///
