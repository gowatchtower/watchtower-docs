# Installation
## Sources

### Container Image

- [Docker Hub](https://hub.docker.com/r/nickfedor/watchtower)
- [GitHub Container Registry](https://github.com/nicholas-fedor/watchtower/pkgs/container/watchtower)

### Binary

- [GitHub Releases](https://github.com/nicholas-fedor/watchtower/releases)

## CLI Commands

### Pulling the Container Image

{{< tabs >}}

{{< tab name="Docker Hub" >}}
```bash
docker pull nickfedor/watchtower:latest
```

{{< /tab >}}

{{< tab name="GitHub" >}}
```bash
docker pull ghcr.io/nicholas-fedor/watchtower:latest
```

{{< /tab >}}

{{< /tabs >}}

### Downloading the Binary

The following will download and extract the binary to the current directory:

{{< tabs >}}

{{< tab name="Windows (amd64)" >}}
```powershell
iwr (iwr https://api.github.com/repos/nicholas-fedor/watchtower/releases/latest | ConvertFrom-Json).assets.where({$_.name -like "*windows_amd64*.zip"}).browser_download_url -OutFile watchtower.zip; Add-Type -AssemblyName System.IO.Compression.FileSystem; $zip = [System.IO.Compression.ZipFile]::OpenRead("$PWD\watchtower.zip"); $zip.Entries | Where-Object {$_.Name -eq 'watchtower.exe'} | ForEach-Object {[System.IO.Compression.ZipFileExtensions]::ExtractToFile($_, "$PWD\watchtower.exe", $true)}; $zip.Dispose(); Remove-Item watchtower.zip; if (Test-Path ".\watchtower.exe") { Write-Host "Successfully installed watchtower.exe to current directory" } else { Write-Host "Failed to install watchtower.exe" }
```

{{< /tab >}}

{{< tab name="Linux (amd64)" >}}
```bash
curl -L $(curl -s https://api.github.com/repos/nicholas-fedor/watchtower/releases/latest | grep -o 'https://[^"]*linux_amd64[^"]*\.tar\.gz') | tar -xz -C . watchtower && if [ -f ./watchtower ]; then echo "Successfully installed watchtower to current directory"; else echo "Failed to install watchtower"; fi
```

{{< /tab >}}

{{< tab name="macOS (Apple Silicon)" >}}
```bash
curl -L $(curl -s https://api.github.com/repos/nicholas-fedor/watchtower/releases/latest | grep -o 'https://[^"]*macOS_arm64v8[^"]*\.tar\.gz') | tar -xz -C . watchtower && if [ -f ./watchtower ]; then echo "Successfully installed watchtower to current directory"; else echo "Failed to install watchtower"; fi
```

{{< /tab >}}

{{< tab name="macOS (Intel)" >}}
```bash
curl -L $(curl -s https://api.github.com/repos/nicholas-fedor/watchtower/releases/latest | grep -o 'https://[^"]*macOS_amd64[^"]*\.tar\.gz') | tar -xz -C . watchtower && if [ -f ./watchtower ]; then echo "Successfully installed watchtower to current directory"; else echo "Failed to install watchtower"; fi
```

{{< /tab >}}

{{< /tabs >}}

{{< callout type="info" >}}
Review the [release page](https://github.com/nicholas-fedor/watchtower/releases) for additional architectures (e.g., arm, arm64, i386, riscv64).

{{< /callout >}}

## Supported Architectures

Watchtower supports the following architectures:

- amd64
- i386
- armhf
- arm64v8
- riscv64

