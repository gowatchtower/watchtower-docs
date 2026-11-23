# Authentication
## Overview

Watchtower's HTTP API uses token-based authentication to protect sensitive endpoints.

Two separate tokens are available:

- **[HTTP API Token](../../../configuration/http-api/#http-api-token)**: Used for most authenticated endpoints.
- **[HTTP API Events Token](../../../configuration/http-api/#http-api-events-token)**: Used exclusively for the [`/v1/events`](../../endpoints/events/) endpoint.

Authentication is enforced at the route level after other middleware (rate limiting, CORS, etc.).

## HTTP API Token

Use the [HTTP API Token](../../../configuration/http-api/#http-api-token) configuration option to set the primary authentication token.

This token is required when any of the following endpoints are enabled:

- [`/v1/check`](../../endpoints/check/)
- [`/v1/config`](../../endpoints/config/)
- [`/v1/containers`](../../endpoints/containers/)
- [`/v1/containers/details`](../../endpoints/container-details/)
- [`/v1/history`](../../endpoints/history/)
- [`/v1/images`](../../endpoints/images/)
- [`/v1/metrics`](../../endpoints/metrics/)
- [`/v1/status`](../../endpoints/status/)
- [`/v1/update`](../../endpoints/update/)

Clients send the token in the `Authorization` header, usually as `Authorization: Bearer API_TOKEN`:

```bash
curl -H "Authorization: Bearer API_TOKEN" http://localhost:8080/v1/metrics
```

Replace *`API_TOKEN`* with your HTTP API token.

The header may also carry the token on its own, as `Authorization: API_TOKEN`. Swagger UI sends the token this way when you authorize with it.

Invalid or missing tokens result in `401 Unauthorized`.
Failed authentication attempts are logged with the client IP address.

A cookie-based fallback is also supported for clients that cannot set custom headers.
The token may be provided in a cookie named `access_token`:

```http
Cookie: access_token=API_TOKEN
```

{{< callout type="info" >}}
The cookie fallback exists primarily for browser-based or limited clients.
The header `Authorization: Bearer` is the recommended method.

{{< /callout >}}

{{< callout type="warning" >}}
**Cookie auth and CSRF**

When the [`http-api-cors-origins`](../../../configuration/http-api/#http-api-cors-origins) is configured and credentials are allowed (for example, `Access-Control-Allow-Credentials: true`), the cookie fallback is vulnerable to cross-site request forgery (CSRF).
A browser on an attacker-controlled origin can issue authenticated POST requests to the API without the user's knowledge.
Prefer header-based `Authorization: Bearer` auth.

{{< /callout >}}

## HTTP API Events Token

The [`/v1/events`](../../endpoints/events/) endpoint uses a separate token set via the [HTTP API Events Token](../../../configuration/http-api/#http-api-events-token) configuration option.

This token is **required** when the [`/v1/events`](../../endpoints/events/) endpoint is enabled.

The events token can be supplied in three ways:

- `Authorization: Bearer` header (recommended for most clients)
- `Authorization` header with the token on its own, as Swagger UI sends it
- `access_token` query parameter (required for browser `EventSource`, which cannot set custom headers)

Example using the header:

```bash
curl -N -H "Authorization: Bearer EVENTS_TOKEN" http://localhost:8080/v1/events
```

Replace *`EVENTS_TOKEN`* with your events token.

Example using the query parameter:

```bash
curl -N -G --data-urlencode "access_token=EVENTS_TOKEN" "http://localhost:8080/v1/events"
```

In JavaScript (for browsers):

```javascript
const eventSource = new EventSource('http://localhost:8080/v1/events?access_token=' + encodeURIComponent('EVENTS_TOKEN'));
```

{{< callout type="warning" >}}
The events token is intentionally separate from the main API token. Query parameters can appear in access logs, browser history, and proxy logs. Using a dedicated token limits exposure.

{{< /callout >}}

See the [Events endpoint documentation](../../endpoints/events/) for more details.

## Examples

Tokens can be provided to Watchtower using Docker Secrets, environment variables, or CLI flags.

{{< tabs >}}

{{< tab name="Docker Compose" >}}
{{< tabs >}}

{{< tab name="Docker Secrets" >}}
Use Compose secrets to mount the token file inside the container and point the environment variable at the mounted path (commonly `/run/secrets/SECRET_NAME`, where *`SECRET_NAME`* is the name of the Compose secret).

```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
        secrets:
            - http_api_token
        environment:
            - WATCHTOWER_HTTP_API_TOKEN=/run/secrets/http_api_token
            # Enable an endpoint that requires authentication
            - WATCHTOWER_HTTP_API_ENDPOINTS=metrics
        ports:
            - "8080:8080"
        restart: unless-stopped

secrets:
    http_api_token:
        file: ./secrets/http_api_token.txt
```

{{< /tab >}}

{{< tab name="Environment  Variables" >}}
```yaml
services:
    watchtower:
        image: nickfedor/watchtower:latest
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
        environment:
            - WATCHTOWER_HTTP_API_TOKEN=API_TOKEN
            # Enable an endpoint that requires authentication
            - WATCHTOWER_HTTP_API_ENDPOINTS=metrics
        ports:
            - "8080:8080"
        restart: unless-stopped
```

{{< /tab >}}

{{< /tabs >}}

{{< /tab >}}

{{< tab name="Docker CLI" >}}
{{< tabs >}}

{{< tab name="Environment Variables" >}}
```bash
docker run -d \
    --name watchtower \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -e WATCHTOWER_HTTP_API_TOKEN=API_TOKEN \
    -e WATCHTOWER_HTTP_API_ENDPOINTS=metrics \
    -p 8080:8080 \
    --restart unless-stopped \
    nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="CLI Flags" >}}
```bash
docker run -d \
    --name watchtower \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -p 8080:8080 \
    --restart unless-stopped \
    nickfedor/watchtower \
    --http-api-token=API_TOKEN \
    --http-api-endpoints=metrics
```

{{< /tab >}}

{{< /tabs >}}

{{< /tab >}}

{{< /tabs >}}

Replace *`API_TOKEN`* with your HTTP API token.

## Unauthenticated Endpoints

The following endpoints do not require authentication when enabled:

- [Health probes](../../endpoints/health/): `/livez`, `/readyz`, `/startupz`
- [Swagger UI](../../endpoints/swagger/): `/swagger/*` ("Try it out" functionality still requires authorization)

## Best Practices

- Generate strong, random tokens (for example: `openssl rand -base64 32`).
- Use a separate events token when enabling the [`/v1/events`](../../endpoints/events/) endpoint.
- Always run the HTTP API behind TLS in production (see the [TLS documentation](../tls/)).
- **Never** expose the HTTP API directly to the public internet.
- Store tokens using Docker Secrets or a secrets manager rather than plain environment variables or CLI flags when possible.
- Rotate tokens if they may have been exposed.
- Combine authentication with network controls (firewalls, reverse proxies with IP allowlists).

## Related Documentation

- [HTTP API Configuration](../../../configuration/http-api/)
- [HTTP API TLS](../tls/)
- [HTTP API Host and Port](../host-and-port/)
- [Events Endpoint](../../endpoints/events/)
- [Docker Secrets](../../../getting-started/docker-secrets/)
- [HTTP API Overview](../../overview/#security-considerations)

