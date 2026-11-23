# Microsoft Teams Notifications
{{< callout type="warning" >}}
**Deprecated Teams Notification Option**

The following legacy Teams notification option is deprecated and subject to removal with the release of Watchtower v2:

- [`notification-msteams-hook`](../../../configuration/notifications/#microsoft-teams-hook)

Use the [`notification-url`](../../../configuration/notifications/#notification-url) configuration option with a `teams://` Shoutrrr URL instead.

{{< /callout >}}

## Overview

Watchtower can use Shoutrrr's [Teams service](https://shoutrrr.nickfedor.com/latest/services/chat/teams/) to send [Microsoft Teams](https://www.microsoft.com/en-us/microsoft-teams) notifications.

## Examples

{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  watchtower:
    image: nickfedor/watchtower:latest
    environment:
      WATCHTOWER_NOTIFICATION_URL: "teams://?host=https%3A%2F%2Fprod-00.westus.logic.azure.com%3A443%2Fworkflows%2Fabc123%2Ftriggers%2Fmanual%2Fpaths%2Finvoke%3Fapi-version%3D2016-06-00%26sp%3D%2Ftriggers%2Fmanual%2Frun%26sv%3D1.0%26sig%3DXXXXXXXX&title=Alert"
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
  --notification-url "teams://?host=https%3A%2F%2Fprod-00.westus.logic.azure.com%3A443%2Fworkflows%2Fabc123%2Ftriggers%2Fmanual%2Fpaths%2Finvoke%3Fapi-version%3D2016-06-00%26sp%3D%2Ftriggers%2Fmanual%2Frun%26sv%3D1.0%26sig%3DXXXXXXXX&title=Alert"
```

{{< /tab >}}

{{< /tabs >}}

/// details | The following legacy Microsoft Teams examples are deprecated and will be removed with the release of Watchtower v2.
    type: warning

{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  watchtower:
    image: nickfedor/watchtower:latest
    environment:
      WATCHTOWER_NOTIFICATIONS: msteams
      WATCHTOWER_NOTIFICATION_MSTEAMS_HOOK_URL: "https://prod-00.westus.logic.azure.com:443/workflows/abc123/triggers/manual/paths/invoke?api-version=2016-06-00&sp=/triggers/manual/run&sv=1.0&sig=XXXXXXXX"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_NOTIFICATIONS=msteams \
  -e WATCHTOWER_NOTIFICATION_MSTEAMS_HOOK_URL="https://prod-00.westus.logic.azure.com:443/workflows/abc123/triggers/manual/paths/invoke?api-version=2016-06-00&sp=/triggers/manual/run&sv=1.0&sig=XXXXXXXX" \
  nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  nickfedor/watchtower \
  --notifications msteams \
  --notification-msteams-hook "https://prod-00.westus.logic.azure.com:443/workflows/abc123/triggers/manual/paths/invoke?api-version=2016-06-00&sp=/triggers/manual/run&sv=1.0&sig=XXXXXXXX"
```

{{< /tab >}}

{{< /tabs >}}

///

