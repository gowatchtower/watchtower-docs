# Slack Notifications
{{< callout type="warning" >}}
**Deprecated Slack Notification Options**

The following legacy Slack notification options are deprecated and subject to removal with the release of Watchtower v2:

- [`notification-slack-hook-url`](../../../configuration/notifications/#slack-hook-url)
- [`notification-slack-identifier`](../../../configuration/notifications/#slack-identifier)
- [`notification-slack-channel`](../../../configuration/notifications/#slack-channel)

Use the [`notification-url`](../../../configuration/notifications/#notification-url) configuration option with a `slack://` Shoutrrr URL instead.

{{< /callout >}}

## Overview

Watchtower can use Shoutrrr's [Slack service](https://shoutrrr.nickfedor.com/latest/services/chat/slack/) to send [Slack](https://slack.com/) notifications.

## Examples

{{< tabs >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  watchtower:
    image: nickfedor/watchtower:latest
    environment:
      WATCHTOWER_NOTIFICATION_URL: "slack://hook:AAAA-BBBB-CCCC@webhook?botname=watchtower"
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
  --notification-url "slack://hook:AAAA-BBBB-CCCC@webhook?botname=watchtower"
```

{{< /tab >}}

{{< /tabs >}}

/// details | The following legacy Slack configuration options and examples are deprecated and will be removed with the release of Watchtower v2.
    type: warning

{{< tabs >}}

{{< tab name="Docker CLI (Env Vars)" >}}
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_NOTIFICATIONS=slack \
  -e WATCHTOWER_NOTIFICATION_SLACK_HOOK_URL="https://hooks.slack.com/services/xxx/yyyyyyyyyyyyyyy" \
  -e WATCHTOWER_NOTIFICATION_SLACK_IDENTIFIER=watchtower-server-1 \
  -e WATCHTOWER_NOTIFICATION_SLACK_CHANNEL=C0123456789 \
  nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  nickfedor/watchtower \
  --notifications slack \
  --notification-slack-hook-url "https://hooks.slack.com/services/xxx/yyyyyyyyyyyyyyy" \
  --notification-slack-identifier watchtower-server-1 \
  --notification-slack-channel "C0123456789"
```

{{< /tab >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  watchtower:
    image: nickfedor/watchtower:latest
    environment:
      WATCHTOWER_NOTIFICATIONS: slack
      WATCHTOWER_NOTIFICATION_SLACK_HOOK_URL: "https://hooks.slack.com/services/xxx/yyyyyyyyyyyyyyy"
      WATCHTOWER_NOTIFICATION_SLACK_IDENTIFIER: watchtower-server-1
      WATCHTOWER_NOTIFICATION_SLACK_CHANNEL: "C0123456789"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< /tabs >}}

///

