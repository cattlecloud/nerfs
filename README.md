# nerfs

[![Go Reference](https://pkg.go.dev/badge/cattlecloud.net/go/nerfs.svg)](https://pkg.go.dev/cattlecloud.net/go/nerfs)
[![License](https://img.shields.io/github/license/cattlecloud/nerfs?color=7C00D8&style=flat-square&label=License)](https://github.com/cattlecloud/nerfs/blob/main/LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/cattlecloud/nerfs/ci.yaml?style=flat-square&color=0FAA07&label=Tests)](https://github.com/cattlecloud/nerfs/actions/workflows/ci.yaml)

`nerfs` provides a tool intended to be run as a cronjob that collects domains known to be
malicious / nsfw / unwanted content. It also provides a client library for accessing the
produced artifacts, which are simply written as files on disk.

### Getting started

Note that `nerfs` may be updated with breaking changes at any time; fork this repository if you
intend to use either the cronjob or client library.

The `nerfs` library package can be added to a Go project with `go get`.

```shell
go get cattlecloud.net/go/nerfs@latest
```

```go
import "cattlecloud.net/go/nerfs"
```

### Sources

The vast majority of sources come from the very popular adblock lists by StevenBlack.

See the [hosts](https://github.com/StevenBlack/hosts) files.

### Example running cronjob

```
➜ .bin/nerfs build -o /tmp
06/20 09:25:43 INFO  [domains-builder] starting the build ...
06/20 09:25:43 INFO  [domains-builder] writing artifact to /tmp/domains.txt
06/20 09:25:43 INFO  [domains-builder] complete in 757.083293ms
06/20 09:25:43 INFO  [wordlist-builder] starting the build ...
06/20 09:25:43 INFO  [wordlist-builder] writing artifact to /tmp/wordlist.json
06/20 09:25:44 INFO  [wordlist-builder] complete in 36.967782ms
```

### Systemd Timer

An example systemd unit file; assuming `nerfs` system user.

```
[Unit]
Description=Build nerf file artifacts.

[Service]
Type=oneshot
ExecStart=/opt/bin/nerfs build -o /path/to/output
User=nerfs
Group=nerfs

ReadWritePaths=/path/to/output
PrivateTmp=yes
NoNewPrivileges=true
ProtectSystem=strict

CPUQuota=100%
CPUWeight=100
OOMScoreAdjust=900
MemoryLow=16M
MemoryMax=64M
```

An example systemd timer file

```
[Unit]
Description=Run nerfs every 24 hours.

[Timer]
OnBootSec=5min
OnUnitActiveSec=24h
Unit=nerfs.service

[Install]
WantedBy=timers.target
```
