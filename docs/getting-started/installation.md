# Installation

Get MockingJay installed on your system.

---

## Prerequisites

- Go 1.26 or later
- Git

---

## Install from Source

The recommended way to install MockingJay:

```bash
# Clone the repository
git clone https://github.com/ashczar77/mockingjay.git
cd mockingjay/cli

# Build the CLI
go build -o mockingjay

# Move to your PATH (optional)
sudo mv mockingjay /usr/local/bin/
```

Verify installation:

```bash
mockingjay version
```

---

## Install from Release

!!! info "Coming Soon"
    Pre-built binaries will be available soon. For now, install from source.

Once releases are available:

=== "Linux"

    ```bash
    curl -L https://github.com/ashczar77/mockingjay/releases/latest/download/mockingjay-linux-amd64 -o mockingjay
    chmod +x mockingjay
    sudo mv mockingjay /usr/local/bin/
    ```

=== "macOS"

    ```bash
    curl -L https://github.com/ashczar77/mockingjay/releases/latest/download/mockingjay-darwin-arm64 -o mockingjay
    chmod +x mockingjay
    sudo mv mockingjay /usr/local/bin/
    ```

=== "Windows"

    Download from [Releases](https://github.com/ashczar77/mockingjay/releases) and add to PATH.

---

## Install with Go

```bash
go install github.com/ashczar77/mockingjay@latest
```

---

## Verify Installation

```bash
mockingjay version
```

Expected output:
```
MockingJay v0.1.0
```

---

## Next Steps

- [Quick Start Guide](quickstart.md) - Get started in 5 minutes
- [Configuration](../configuration/overview.md) - Learn how to configure tests
