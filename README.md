# kiosk

| Environment         | Status                                                                                               |
|---------------------|------------------------------------------------------------------------------------------------------|
| Release (main)      | [![Release](https://github.com/jmelowry/kiosk/actions/workflows/release.yml/badge.svg)](https://github.com/jmelowry/kiosk/actions/workflows/release.yml) |
| Pre-Release (dev)   | [![Pre-Release](https://github.com/jmelowry/kiosk/actions/workflows/pre-release.yml/badge.svg?branch=dev)](https://github.com/jmelowry/kiosk/actions/workflows/pre-release.yml) |
| CI (main)           | [![CI - main](https://github.com/jmelowry/kiosk/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/jmelowry/kiosk/actions/workflows/test.yml) |
| CI (dev)            | [![CI - dev](https://github.com/jmelowry/kiosk/actions/workflows/test.yml/badge.svg?branch=dev)](https://github.com/jmelowry/kiosk/actions/workflows/test.yml) |

**kiosk** is a terminal-based tool designed to reduce the cognitive overhead of managing code and services across different environments. Instead of juggling terminal tabs, virtual environments, and tmux sessions, you can step into your kiosk and find everything set up just the way you need it.

This project is in early development and aims to make life a little easier for developers working with complex, heterogeneous infrastructure and platforms.

## Installation

### From Source
```sh
git clone https://github.com/jmelowry/kiosk.git
go build
```

### From Releases
```sh
curl -L https://github.com/jmelowry/kiosk/releases/latest/download/kiosk.tar.gz | tar xz
sudo mv kiosk /usr/local/bin
sudo chmod +x /usr/local/bin/kiosk
```
