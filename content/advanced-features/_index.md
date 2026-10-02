---
title: Advanced Features
description: Advanced deployment patterns and operational features.
weight: 6
type: docs
cascade:
  type: docs
---

<p class="wt-section-lead">Advanced deployment patterns and operational features.</p>

{{< cards >}}
  {{< card link="ephemeral-self-updates/" title="Ephemeral Self-Updates" icon="refresh" subtitle="Update Watchtower itself safely" >}}
  {{< card link="git-monitoring/" title="Git Monitoring" icon="code" subtitle="Rebuild images from Git refs" >}}
  {{< card link="image-cooldown/" title="Image Cooldown" icon="pause" subtitle="Avoid update thrashing" >}}
  {{< card link="disk-usage-monitor/" title="Disk Usage Monitor" icon="database" subtitle="Cap Docker image storage usage" >}}
  {{< card link="copy-files/" title="Copy Files" icon="document-duplicate" subtitle="Copy labeled files across recreates" >}}
  {{< card link="lifecycle-hooks/" title="Lifecycle Hooks" icon="code" subtitle="Run commands around updates" >}}
  {{< card link="linked-containers/" title="Linked Containers" icon="link" subtitle="Coordinate dependent containers" >}}
  {{< card link="private-registries/" title="Private Registries" icon="lock-closed" subtitle="Authenticate to private images" >}}
  {{< card link="registry-mirrors/" title="Registry Mirrors" icon="duplicate" subtitle="Use registry mirrors" >}}
  {{< card link="remote-hosts/" title="Remote Hosts" icon="server" subtitle="Manage remote Docker hosts" >}}
  {{< card link="running-multiple-instances/" title="Running Multiple Instances" icon="collection" subtitle="Isolate multiple Watchtowers" >}}
  {{< card link="secure-connections/" title="Secure Connections" icon="shield-check" subtitle="TLS to the Docker API" >}}
  {{< card link="stop-signals/" title="Stop Signals" icon="stop" subtitle="Custom container stop signals" >}}
{{< /cards >}}
