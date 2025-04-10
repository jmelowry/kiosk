# kiosk

| Environment         | Status                                                                                               |
|---------------------|------------------------------------------------------------------------------------------------------|
| Release (main)      | [![Release](https://github.com/jmelowry/kiosk/actions/workflows/release.yml/badge.svg)](https://github.com/jmelowry/kiosk/actions/workflows/release.yml) |
| Pre-Release (dev)   | [![Pre-Release](https://github.com/jmelowry/kiosk/actions/workflows/pre-release.yml/badge.svg?branch=dev)](https://github.com/jmelowry/kiosk/actions/workflows/pre-release.yml) |
| CI (main)           | [![CI - main](https://github.com/jmelowry/kiosk/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/jmelowry/kiosk/actions/workflows/test.yml) |
| CI (dev)            | [![CI - dev](https://github.com/jmelowry/kiosk/actions/workflows/test.yml/badge.svg?branch=dev)](https://github.com/jmelowry/kiosk/actions/workflows/test.yml) |

**kiosk** is a terminal-based tool designed to simplify the management of code and services across different environments.

This project is in early development and aims to assist developers working with complex infrastructure and platforms.

## Installation

You can install `kiosk` in a few different ways depending on your platform.

---

### Option 1: Prebuilt Binaries

Go to the [Releases page](https://github.com/jmelowry/kiosk/releases) and download the appropriate `.tar.gz` file for your system.

Example for macOS (arm64):

```sh
curl -L https://github.com/jmelowry/kiosk/releases/latest/download/kiosk-main-<timestamp>-darwin-arm64.tar.gz | tar -xz
sudo mv kiosk /usr/local/bin/
```

Example for Linux (amd64):

```sh
curl -L https://github.com/jmelowry/kiosk/releases/latest/download/kiosk-main-<timestamp>-linux-amd64.tar.gz | tar -xz
sudo mv kiosk /usr/local/bin/
```

Replace `<timestamp>` with the version from the release you want.

To verify the download:

```sh
curl -L https://github.com/jmelowry/kiosk/releases/latest/download/kiosk-main-<timestamp>-linux-amd64.tar.gz.sha256 | sha256sum --check
```

---

### Option 2: Build from Source

```sh
git clone https://github.com/jmelowry/kiosk.git
cd kiosk
go build -o kiosk .
sudo mv kiosk /usr/local/bin/
```

This avoids macOS Gatekeeper warnings and gives you a local build.

---

### Option 3: Development Version

Use the latest version from the `dev` branch:

```sh
git clone -b dev https://github.com/jmelowry/kiosk.git
cd kiosk
go build -o kiosk .
sudo mv kiosk /usr/local/bin/
```

---

### Checksum Verification

To verify a downloaded archive:

```sh
sha256sum -c kiosk-main-<timestamp>-darwin-arm64.tar.gz.sha256
```

---

A Homebrew formula is planned and will be added soon.
