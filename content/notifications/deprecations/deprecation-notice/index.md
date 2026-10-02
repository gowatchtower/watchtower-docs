---
title: Deprecation Notice
weight: 1
---
## Overview

Watchtower has a number of legacy notification configuration options that will be removed with the release of Watchtower v2.

This deprecation and eventual removal of legacy notification configuration options is a long-term migration that started with the refactoring of Shoutrrr into a separate notification library. The removal of the legacy notification configuration options is simply a removal of an unnecessary abstraction layer in favor of using Shoutrrr URL's directly.

In order to continue receiving notifications, users will need to add Shoutrrr URL's directly using the [`NOTIFICATION URL`](../../../configuration/notifications/#notification-url) configuration option.

Watchtower users do not need to install any new software to implement Shoutrrr URL's.

Use Watchtower's CLI [`migration tool`](../migration-tool/) to help convert legacy email configurations to Shoutrrr URLs or use the [Shoutrrr Playground](https://shoutrrr.nickfedor.com/latest/playground/) to help convert configurations for other services to Shoutrrr URLs.

## General Options

- [`watchtower-notifications`](../../../configuration/notifications/#notifications)

## [Email Notifications](../email/)

- [`notification-email-from`](../../../configuration/notifications/#email-from)
- [`notification-email-to`](../../../configuration/notifications/#email-to)
- [`notification-email-server`](../../../configuration/notifications/#email-server)
- [`notification-email-server-tls-skip-verify`](../../../configuration/notifications/#email-server-tls-skip-verify)
- [`notification-email-server-user`](../../../configuration/notifications/#email-server-user)
- [`notification-email-server-password`](../../../configuration/notifications/#email-server-password)
- [`notification-email-subjecttag`](../../../configuration/notifications/#email-subject-tag)
- [`notification-email-server-port`](../../../configuration/notifications/#email-server-port)
- [`notification-email-delay`](../../../configuration/notifications/#email-delay)

## [Gotify Notifications](../gotify/)

- [`notification-gotify-url`](../../../configuration/notifications/#gotify-url)
- [`notification-gotify-token`](../../../configuration/notifications/#gotify-token)
- [`notification-gotify-tls-skip-verify`](../../../configuration/notifications/#gotify-tls-skip-verify)

## [Microsoft Teams Notifications](../teams/)

- [`notification-msteams-hook`](../../../configuration/notifications/#microsoft-teams-hook)

## [Slack Notifications](../slack/)

- [`notification-slack-hook-url`](../../../configuration/notifications/#slack-hook-url)
- [`notification-slack-identifier`](../../../configuration/notifications/#slack-identifier)
- [`notification-slack-channel`](../../../configuration/notifications/#slack-channel)
- [`notification-slack-icon-emoji`](../../../configuration/notifications/#slack-icon-emoji)
- [`notification-slack-icon-url`](../../../configuration/notifications/#slack-icon-url)
