---
title: Email Notifications
weight: 3
---
{{< callout type="warning" >}}
**Deprecated Email Notification Options**

The following legacy email notification options are deprecated and subject to removal with the release of Watchtower v2:

- [`notification-email-from`](../../../configuration/notifications/#email-from)
- [`notification-email-to`](../../../configuration/notifications/#email-to)
- [`notification-email-server`](../../../configuration/notifications/#email-server)
- [`notification-email-server-tls-skip-verify`](../../../configuration/notifications/#email-server-tls-skip-verify)
- [`notification-email-server-user`](../../../configuration/notifications/#email-server-user)
- [`notification-email-server-password`](../../../configuration/notifications/#email-server-password)
- [`notification-email-subjecttag`](../../../configuration/notifications/#email-subject-tag)
- [`notification-email-server-port`](../../../configuration/notifications/#email-server-port)
- [`notification-email-delay`](../../../configuration/notifications/#email-delay)

Use the [`notification-url`](../../../configuration/notifications/#notification-url) configuration option with a `smtp://` Shoutrrr URL instead.

Use the [legacy notifications migration tool](../../deprecations/migration-tool/#migration-walkthrough) to help with migrating to a Shoutrrr URL.

{{< /callout >}}

## Overview

Watchtower uses Shoutrrr's [SMTP service](https://shoutrrr.nickfedor.com/latest/services/email/smtp/) to send email notifications.

## Examples

{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
  watchtower:
    image: nickfedor/watchtower:latest
    environment:
      WATCHTOWER_NOTIFICATION_URL: smtp://USERNAME:PASSWORD@smtp.example.com:587/?fromaddress=sender@example.com&toaddresses=recipient@example.com
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="Docker CLI (Env Vars)" >}}

```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_NOTIFICATION_URL="smtp://USERNAME:PASSWORD@smtp.example.com:587/?fromaddress=sender@example.com&toaddresses=recipient@example.com" \
  nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Docker CLI (Flags)" >}}

```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  nickfedor/watchtower \
  --notification-url "smtp://USERNAME:PASSWORD@smtp.example.com:587/?fromaddress=sender@example.com&toaddresses=recipient@example.com"
```

{{< /tab >}}

{{< /tabs >}}

Replace the following:

- *`USERNAME`*: the user name for your SMTP server, with an `@` in it written as `%40`
- *`PASSWORD`*: the password for your SMTP server, percent-encoded

Use your SMTP server's host name in place of `smtp.example.com` and your own addresses in place of `sender@example.com` and `recipient@example.com`.

/// details | The following legacy SMTP configuration examples are deprecated and will be removed with the release of Watchtower v2.
    type: warning

{{< tabs >}}

{{< tab name="Docker Compose" >}}
{{< tabs >}}

{{< tab name="Docker Compose" >}}

```yaml
services:
  watchtower:
    image: nickfedor/watchtower:latest
    environment:
      WATCHTOWER_NOTIFICATIONS: email
      WATCHTOWER_NOTIFICATION_EMAIL_FROM: sender@example.com
      WATCHTOWER_NOTIFICATION_EMAIL_TO: recipient@example.com
      WATCHTOWER_NOTIFICATION_EMAIL_SERVER: smtp.example.com
      WATCHTOWER_NOTIFICATION_EMAIL_SERVER_PORT: 587
      WATCHTOWER_NOTIFICATION_EMAIL_SERVER_USER: USERNAME
      WATCHTOWER_NOTIFICATION_EMAIL_SERVER_PASSWORD: 'PASSWORD'
      WATCHTOWER_NOTIFICATION_EMAIL_DELAY: 10
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

{{< /tab >}}

{{< tab name="SMTP Relay" >}}

```yaml
services:
  watchtower:
    image: nickfedor/watchtower:latest
    environment:
      WATCHTOWER_NOTIFICATIONS: email
      WATCHTOWER_NOTIFICATION_EMAIL_FROM: sender@example.com
      WATCHTOWER_NOTIFICATION_EMAIL_TO: recipient@example.com
      WATCHTOWER_NOTIFICATION_EMAIL_SERVER: relay.example.com
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

{{< callout type="info" >}}
This example sends through an SMTP relay at `relay.example.com` that accepts mail on port 25 without authentication.

If the relay offers STARTTLS, Watchtower checks its certificate against the server name, so the certificate must be valid for `relay.example.com`. A relay in the same Compose project can take that name as a network alias.

{{< /callout >}}

{{< /tab >}}

{{< /tabs >}}

{{< /tab >}}

{{< tab name="Docker CLI" >}}
{{< tabs >}}

{{< tab name="Env Vars" >}}

```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_NOTIFICATIONS=email \
  -e WATCHTOWER_NOTIFICATION_EMAIL_FROM=sender@example.com \
  -e WATCHTOWER_NOTIFICATION_EMAIL_TO=recipient@example.com \
  -e WATCHTOWER_NOTIFICATION_EMAIL_SERVER=smtp.example.com \
  -e WATCHTOWER_NOTIFICATION_EMAIL_SERVER_PORT=587 \
  -e WATCHTOWER_NOTIFICATION_EMAIL_SERVER_USER=USERNAME \
  -e WATCHTOWER_NOTIFICATION_EMAIL_SERVER_PASSWORD='PASSWORD' \
  -e WATCHTOWER_NOTIFICATION_EMAIL_DELAY=10 \
  nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Flags" >}}

```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  nickfedor/watchtower \
  --notifications email \
  --notification-email-from sender@example.com \
  --notification-email-to recipient@example.com \
  --notification-email-server smtp.example.com \
  --notification-email-server-port 587 \
  --notification-email-server-user USERNAME \
  --notification-email-server-password 'PASSWORD' \
  --notification-email-delay 10
```

{{< /tab >}}

{{< tab name="Env Vars (SMTP Relay)" >}}

```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e WATCHTOWER_NOTIFICATIONS=email \
  -e WATCHTOWER_NOTIFICATION_EMAIL_FROM=sender@example.com \
  -e WATCHTOWER_NOTIFICATION_EMAIL_TO=recipient@example.com \
  -e WATCHTOWER_NOTIFICATION_EMAIL_SERVER=relay.example.com \
  nickfedor/watchtower
```

{{< /tab >}}

{{< tab name="Flags (SMTP Relay)" >}}

```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  nickfedor/watchtower \
  --notifications email \
  --notification-email-from sender@example.com \
  --notification-email-to recipient@example.com \
  --notification-email-server relay.example.com
```

{{< callout type="info" >}}
These examples send through an SMTP relay at `relay.example.com` that accepts mail on port 25 without authentication.

If the relay offers STARTTLS, Watchtower checks its certificate against the server name, so the certificate must be valid for `relay.example.com`.

{{< /callout >}}

{{< /tab >}}

{{< /tabs >}}

{{< /tab >}}

{{< /tabs >}}

Replace the following:

- *`USERNAME`*: the user name for your SMTP server
- *`PASSWORD`*: the password for your SMTP server, inside the single quotes. In a Compose file, write each `'` in it as `''` and each `$` as `$$`. In a shell command, write each `'` as `'\''`.

Use your SMTP server's or relay's host name in place of `smtp.example.com` or `relay.example.com` and your own addresses in place of `sender@example.com` and `recipient@example.com`.

///

## Common SMTP Configurations

{{< tabs >}}

{{< tab name="Gmail" >}}

| Property        | Value                     |
|-----------------|---------------------------|
| Port            | `587`                     |
| Authentication  | OAuth2 client credentials |
| RequireStartTLS | `Yes`                     |

```text
smtp://GMAIL_USER%40gmail.com:APP_PASSWORD@smtp.gmail.com:587/?fromaddress=GMAIL_USER%40gmail.com&toaddresses=TO_ADDRESS&encryption=ExplicitTLS&usestarttls=yes&timeout=30s
```

Replace the following:

- *`GMAIL_USER`*: the part of your Gmail address before the `@`
- *`APP_PASSWORD`*: your Gmail password, or an App Password
- *`TO_ADDRESS`*: the address that receives the notifications, with its `@` written as `%40`

{{< callout type="info" >}}
For Gmail, use an [App Password](https://support.google.com/accounts/answer/185833) if two-factor authentication is enabled.

Omit spaces from the App Password. Google displays the 16-character passcode in groups of four, but those spaces are not part of the password.

Do not use `&auth=OAuth2` with an App Password. That method expects a SASL XOAUTH2 access token in the password field.

{{< /callout >}}

{{< /tab >}}

{{< tab name="AWS SES" >}}

| Property        | Value                     |
|-----------------|---------------------------|
| Port            | `587`                     |
| Authentication  | OAuth2 client credentials |
| RequireStartTLS | `Yes`                     |

```text
smtp://SMTP_USERNAME:SMTP_PASSWORD@email-smtp.REGION.amazonaws.com:587/?fromaddress=FROM_ADDRESS&toaddresses=TO_ADDRESS&encryption=ExplicitTLS&usestarttls=yes&timeout=30s
```

Replace the following:

- *`SMTP_USERNAME`*: the user name of your Amazon SES SMTP credentials
- *`SMTP_PASSWORD`*: the password of your Amazon SES SMTP credentials, percent-encoded
- *`REGION`*: the AWS Region of your Amazon SES endpoint, such as `us-east-1`
- *`FROM_ADDRESS`*: a verified sender address, with its `@` written as `%40`
- *`TO_ADDRESS`*: the address that receives the notifications, with its `@` written as `%40`

{{< /tab >}}

{{< tab name="Microsoft 365" >}}

| Property        | Value                     |
|-----------------|---------------------------|
| Port            | `587`                     |
| Authentication  | OAuth2 client credentials |
| RequireStartTLS | `Yes`                     |

```text
smtp://MAILBOX@smtp.office365.com:587/?fromaddress=FROM_ADDRESS&toaddresses=TO_ADDRESS&clientid=CLIENT_ID&clientsecret=CLIENT_SECRET&scopes=https%3A%2F%2Foutlook.office365.com%2F.default&tokenurl=https%3A%2F%2Flogin.microsoftonline.com%2FTENANT_ID%2Foauth2%2Fv2.0%2Ftoken&requirestarttls=yes&timeout=30s
```

Replace the following:

- *`MAILBOX`*: the address of the sending mailbox, with its `@` written as `%40`
- *`FROM_ADDRESS`*: the sender address, with its `@` written as `%40`
- *`TO_ADDRESS`*: the address that receives the notifications, with its `@` written as `%40`
- *`CLIENT_ID`*: the application (client) ID of your app registration
- *`CLIENT_SECRET`*: the client secret of your app registration, percent-encoded
- *`TENANT_ID`*: your Microsoft Entra tenant ID

{{< callout type="info" >}}
Shoutrrr gets OAuth2 access tokens from Microsoft Entra ID with the client credentials grant. Follow Shoutrrr's [Microsoft 365 guide](https://shoutrrr.nickfedor.com/latest/guides/microsoft-365/) to register the application and give it access to the sending mailbox.

Exchange Online turns off Basic authentication for SMTP AUTH by default in existing tenants at the end of December 2026. A URL with a password, such as `smtp://MAILBOX:PASSWORD@smtp.office365.com:587/`, works only while the tenant allows it.

{{< /callout >}}

{{< /tab >}}

{{< tab name="Generic (SSL)" >}}

| Property        | Value                     |
|-----------------|---------------------------|
| Port        | `465`         |
| Encryption  | `ImplicitTLS` |
| UseStartTLS | `No`          |

```text
smtp://USERNAME:PASSWORD@smtp.example.com:465/?fromaddress=FROM_ADDRESS&toaddresses=TO_ADDRESS&encryption=ImplicitTLS&usestarttls=no&timeout=30s
```

Replace the following:

- *`USERNAME`*: the user name for your SMTP server, with an `@` in it written as `%40`
- *`PASSWORD`*: the password for your SMTP server, percent-encoded
- *`FROM_ADDRESS`*: the sender address, with its `@` written as `%40`
- *`TO_ADDRESS`*: the address that receives the notifications, with its `@` written as `%40`

{{< /tab >}}

{{< tab name="Generic (Plain)" >}}

| Property    | Value  |
|-------------|--------|
| Port        | `25`   |
| Encryption  | `None` |
| UseStartTLS | `No`   |

```text
smtp://USERNAME:PASSWORD@smtp.example.com:25/?fromaddress=FROM_ADDRESS&toaddresses=TO_ADDRESS&encryption=None&usestarttls=no&timeout=30s
```

Replace the following:

- *`USERNAME`*: the user name for your SMTP server, with an `@` in it written as `%40`
- *`PASSWORD`*: the password for your SMTP server, percent-encoded
- *`FROM_ADDRESS`*: the sender address, with its `@` written as `%40`
- *`TO_ADDRESS`*: the address that receives the notifications, with its `@` written as `%40`

{{< /tab >}}

{{< /tabs >}}

## Notes

- **Timeout**:

  - The default SMTP timeout is 10 seconds.
  - If you experience timeouts (e.g., `failed to send: timed out: using smtp`), add `&timeout=30s` to the URL to allow more time for server responses, especially with proxies or slow networks.

- **Authentication**:

  - Use `&auth=Plain` for username and password authentication, including [Gmail App Passwords](https://support.google.com/accounts/answer/185833). This is the default when credentials are provided.
  - Set `clientid`, `clientsecret`, `scopes` and `tokenurl` to have Shoutrrr get OAuth2 access tokens with the [client credentials grant](https://shoutrrr.nickfedor.com/latest/services/email/smtp/#oauth2_client_credentials), as the Microsoft 365 example does. Shoutrrr then uses OAuth2 without `&auth=OAuth2`.
  - Use `&auth=OAuth2` only for [SASL XOAUTH2](https://shoutrrr.nickfedor.com/latest/services/email/smtp/#authentication) with a static access token in the password field. App Passwords are not OAuth2.

- **Testing**:

  - Install Shoutrrr using one of the various [installation methods](https://shoutrrr.nickfedor.com/installation/).
  - Test your URL with the Shoutrrr CLI:

      ```bash
      shoutrrr send -u "NOTIFICATION_URL" -m "Test message"
      ```

      Replace *`NOTIFICATION_URL`* with the `smtp://` URL to test.

- **Proxy Issues**:

  - If using a Docker proxy (e.g., `tcp://dockerproxy:2375`), ensure it allows outbound connections to your SMTP server's host and port.
  - Test connectivity with `telnet SMTP_HOST SMTP_PORT` inside the container. Replace *`SMTP_HOST`* with your SMTP server's host name and *`SMTP_PORT`* with its port.
