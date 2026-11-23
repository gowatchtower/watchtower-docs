---
title: Gotify Notifications
weight: 4
---
{{< callout type="warning" >}}
**Deprecated Gotify Notification Options**

The following legacy Gotify notification options are deprecated and subject to removal with the release of Watchtower v2:

- [`notification-gotify-url`](../../../configuration/notifications/#gotify-url)
- [`notification-gotify-token`](../../../configuration/notifications/#gotify-token)
- [`notification-gotify-tls-skip-verify`](../../../configuration/notifications/#gotify-tls-skip-verify)

Use the [`notification-url`](../../../configuration/notifications/#notification-url) configuration option with a `gotify://` Shoutrrr URL instead.

{{< /callout >}}

## Overview

Watchtower can use Shoutrrr's [Gotify service](https://shoutrrr.nickfedor.com/latest/services/push/gotify/) to send [Gotify](https://gotify.net/) notifications.

## Examples

{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  watchtower:
    image: nickfedor/watchtower:latest
    environment:
      WATCHTOWER_NOTIFICATION_URL: "gotify://gotify.example.com/TOKEN"
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
  --notification-url "gotify://gotify.example.com/TOKEN"
```

{{< /tab >}}

{{< /tabs >}}

Replace *`TOKEN`* with the application token from your Gotify server. Use your server's host name in place of `gotify.example.com`.

/// details | The following legacy Gotify configuration examples are deprecated and will be removed with the release of Watchtower v2.
    type: warning

{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  watchtower:
    image: nickfedor/watchtower:latest
    environment:
      WATCHTOWER_NOTIFICATIONS: gotify
      WATCHTOWER_NOTIFICATION_GOTIFY_URL: "https://gotify.example.com/"
      WATCHTOWER_NOTIFICATION_GOTIFY_TOKEN: "TOKEN"
      WATCHTOWER_NOTIFICATION_GOTIFY_TLS_SKIP_VERIFY: true
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_NOTIFICATIONS=gotify \
  -e WATCHTOWER_NOTIFICATION_GOTIFY_URL="https://gotify.example.com/" \
  -e WATCHTOWER_NOTIFICATION_GOTIFY_TOKEN="TOKEN" \
  -e WATCHTOWER_NOTIFICATION_GOTIFY_TLS_SKIP_VERIFY=true \
  nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  nickfedor/watchtower \
  --notifications gotify \
  --notification-gotify-url "https://gotify.example.com/" \
  --notification-gotify-token "TOKEN" \
  --notification-gotify-tls-skip-verify
```

{{< /tab >}}

{{< /tabs >}}

///
