---
title: Deprecation Notice
weight: 1
---
## Overview

Watchtower has a number of legacy notification configuration options that will be removed with the release of Watchtower v2.

This deprecation and eventual removal of legacy notification configuration options is a long-term migration that started with the refactoring of Shoutrrr into a separate notification library. The removal of the legacy notification configuration options is simply a removal of an unnecessary abstraction layer in favor of using Shoutrrr URL's directly.

In order to continue receiving notifications, users will need to add Shoutrrr URL's directly using the [`NOTIFICATION URL`](../../../configuration/notifications/#notification_url) configuration option.

Watchtower users do not need to install any new software to implement Shoutrrr URL's.

Use Watchtower's CLI [`migration tool`](../migration-tool/) to help convert legacy email configurations to Shoutrrr URLs or use the [Shoutrrr Playground](https://shoutrrr.nickfedor.com/latest/playground/) to help convert configurations for other services to Shoutrrr URLs.

## General Options

- [`watchtower-notifications`](../../../configuration/notifications/#notifications_1)

## [Email Notifications](../email/)

- [`notification-email-from`](../../../configuration/notifications/#email_from)
- [`notification-email-to`](../../../configuration/notifications/#email_to)
- [`notification-email-server`](../../../configuration/notifications/#email_server)
- [`notification-email-server-tls-skip-verify`](../../../configuration/notifications/#email_server_tls_skip_verify)
- [`notification-email-server-user`](../../../configuration/notifications/#email_server_user)
- [`notification-email-server-password`](../../../configuration/notifications/#email_server_password)
- [`notification-email-subjecttag`](../../../configuration/notifications/#email_subject_tag)
- [`notification-email-server-port`](../../../configuration/notifications/#email_server_port)
- [`notification-email-delay`](../../../configuration/notifications/#email_delay)

## [Gotify Notifications](../gotify/)

- [`notification-gotify-url`](../../../configuration/notifications/#gotify_url)
- [`notification-gotify-token`](../../../configuration/notifications/#gotify_token)
- [`notification-gotify-tls-skip-verify`](../../../configuration/notifications/#gotify_tls_skip_verify)

## [Microsoft Teams Notifications](../teams/)

- [`notification-msteams-hook`](../../../configuration/notifications/#microsoft_teams_hook)

## [Slack Notifications](../slack/)

- [`notification-slack-hook-url`](../../../configuration/notifications/#slack_hook_url)
- [`notification-slack-identifier`](../../../configuration/notifications/#slack_identifier)
- [`notification-slack-channel`](../../../configuration/notifications/#slack_channel)
- [`notification-slack-icon-emoji`](../../../configuration/notifications/#slack_icon_emoji)
- [`notification-slack-icon-url`](../../../configuration/notifications/#slack_icon_url)
