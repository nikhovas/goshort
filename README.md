# GoShort

![Linter](https://img.shields.io/github/workflow/status/nikhovas/goshort/lint?style=flat-square&label=Linter)
![Testing](https://img.shields.io/github/workflow/status/nikhovas/goshort/test?style=flat-square&label=Testing)
![Dockerized](https://img.shields.io/github/workflow/status/nikhovas/goshort/dockerize?style=flat-square&label=Dockerize)

Url shortener written in Go. Supports generic/specified key setting and REST API for urls.
Needs Redis. Docker-friendly.

* [Build requirements](#build-requirements)
* [Docker](#docker)
* [Build](#build)
* [Configuration](#configuration)
* [Rest API](#rest-api)

  

## Build requirements

Example for OS Ubuntu with redis/SQL data storage.

|Name|Minimal version| Installation  |
|---|---|---|
|OS Ubuntu|18.04|https://ubuntu.com/tutorials/install-ubuntu-desktop|
|Golang|1.14.2| https://golang.org/doc/install  |
|go.mod libraries|-|`go get dep_name` for all dep_name in `go.mod` file (optionally, go will automatically install deps during go build)|
|Redis|6.2|https://redis.io/topics/quickstart|
|Postgres|12.0|`sudo apt -y install postgresql`

## Docker

`nikhovas/goshort:alpine` in DockerHub

## Build

Just `go build -o urlshort .` in the folder of repository.

## Configuration

GoShort uses config file which is the first parameter to start a program.

Example of config file:

```yml
inputs:
  server:
    name: Server
    ip: ''
    port: 80
    mode: tcp
    token: asdf
database:
  redis:
    name: Redis
    ip: 127.0.0.1:6379
    port: 6379
    mode: tcp
    pool_size: 10
loggers:
  console:
    name: consoleLogger
    extra_logger: true
    common_logger: true
middlewares:
  - url_normalizer
limits:
  max_connections: 2000
```

## Rest API

See `docs/swagger.yaml` for swagger Rest API.

To regenerate `swag init` in the root of the project.

### Authorization

For simplicity the project uses a simple Bearer token auth, which you can specify with
`GOSHORT_TOKEN` env var or in config file. The token in single, so make sure that it is
really secret.

If no token is specified, you don't have to write anything to `Authorization` header,
the program won't look to it.
If the token is specified, but something is wrong with auth in your request, you'll get
error 401 (Unauthorized).
