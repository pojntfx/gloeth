# gloeth

A global ethernet.

## Installation

A Go package [is available](https://godoc.org/github.com/pojntfx/gloeth).

## Usage

```bash
% gloeth --help   
Usage of gloeth:
  -device string
        Name of the network device to create (default "gloeth0")
  -listen string
        Host:port to listen on (default "127.0.0.1:1234")
  -peer string
        Host:port the peer listens on (default "127.0.0.1:1235")
  -redis-host string
        Host:port of Redis (default "127.0.0.1:6379")
  -redis-password string
        Password for Redis
```

## License

gloeth (c) 2019 Felicitas Pojtinger

SPDX-License-Identifier: AGPL-3.0