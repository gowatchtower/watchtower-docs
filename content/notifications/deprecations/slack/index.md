---
title: Slack Notifications
weight: 5
---
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
      WATCHTOWER_NOTIFICATION_URL: "slack://hook:WORKSPACE_ID-WEBHOOK_ID-WEBHOOK_TOKEN@webhook?botname=watchtower"
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
  --notification-url "slack://hook:WORKSPACE_ID-WEBHOOK_ID-WEBHOOK_TOKEN@webhook?botname=watchtower"
```

{{< /tab >}}

{{< /tabs >}}

Replace the following with the parts of your Slack webhook URL, `https://hooks.slack.com/services/WORKSPACE_ID/WEBHOOK_ID/WEBHOOK_TOKEN`:

- *`WORKSPACE_ID`*: the workspace ID, which starts with `T`
- *`WEBHOOK_ID`*: the webhook ID, which starts with `B`
- *`WEBHOOK_TOKEN`*: the webhook token

/// details | The following legacy Slack configuration options and examples are deprecated and will be removed with the release of Watchtower v2.
    type: warning

{{< tabs >}}

{{< tab name="Docker CLI (Env Vars)" >}}
```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_NOTIFICATIONS=slack \
  -e WATCHTOWER_NOTIFICATION_SLACK_HOOK_URL="https://hooks.slack.com/services/WORKSPACE_ID/WEBHOOK_ID/WEBHOOK_TOKEN" \
  -e WATCHTOWER_NOTIFICATION_SLACK_IDENTIFIER=watchtower-server-1 \
  -e WATCHTOWER_NOTIFICATION_SLACK_CHANNEL=CHANNEL_ID \
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
  --notification-slack-hook-url "https://hooks.slack.com/services/WORKSPACE_ID/WEBHOOK_ID/WEBHOOK_TOKEN" \
  --notification-slack-identifier watchtower-server-1 \
  --notification-slack-channel "CHANNEL_ID"
```

{{< /tab >}}

{{< tab name="Docker Compose" >}}
```yaml
services:
  watchtower:
    image: nickfedor/watchtower:latest
    environment:
      WATCHTOWER_NOTIFICATIONS: slack
      WATCHTOWER_NOTIFICATION_SLACK_HOOK_URL: "https://hooks.slack.com/services/WORKSPACE_ID/WEBHOOK_ID/WEBHOOK_TOKEN"
      WATCHTOWER_NOTIFICATION_SLACK_IDENTIFIER: watchtower-server-1
      WATCHTOWER_NOTIFICATION_SLACK_CHANNEL: "CHANNEL_ID"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< /tabs >}}

Replace the following:

- *`WORKSPACE_ID`*: the workspace ID, which starts with `T`
- *`WEBHOOK_ID`*: the webhook ID, which starts with `B`
- *`WEBHOOK_TOKEN`*: the webhook token
- *`CHANNEL_ID`*: the ID of the channel to post to, such as `C0123456789`

///
