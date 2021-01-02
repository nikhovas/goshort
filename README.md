# GoShort

Url shortener written in Go. Supports generic/specified key setting and REST api for urls.
Needs Redis. Docker-friendly.

- [GoShort](#goshort)
    * [Build status](#build-status)
    * [Configuration](#configuration)
    * [Rest API](#rest-api)
        + [Examples](#examples)
            - [Key specified POST request](#key-specified-post-request)
            - [Generic key POST request](#generic-key-post-request)
            - [GET request](#get-request)
            - [PATCH/PUT request](#patch-put-request)
            - [DELETE request](#delete-request)
        + [Authorization](#authorization)

## Build status

![Linter](https://github.com/nikhovas/goshort/workflows/lint/badge.svg)
![Testing](https://github.com/nikhovas/goshort/workflows/test/badge.svg)
![Dockerized](https://github.com/nikhovas/goshort/workflows/dockerize/badge.svg)

## Configuration

You can use enviromental variables or config file to setup app.

| Name                   | Description                                        | Default value  |
|------------------------|----------------------------------------------------|----------------|
| GOSHORT_PORT           | Port to listen to connections. App uses localhost. | 80             |
| GOSHORT_TOKEN          | Auth token for managing with urls                  |                |
| GOSHORT_REDIS_IP       | Address of redis server                            | 127.0.0.1:6379 |
| GOSHORT_REDIS_POOLSIZE | Number of opened connections                       | 10             |

Example of config file:

```yml
todo
```

## Rest API

You can manage urls with rest api on /urls/

### Examples

#### Key specified POST request

```http request
POST /urls/ HTTP/1.1
Content-Type: application/json
Authorization: Bearer asdf
Host: 127.0.0.1:80
Content-Length: 49

{
  "key": "yandex",
  "url": "https://yandex.ru"
}
```

The default answer's code status is 301 (Moved permanently).
If all is OK, the answer should be:
```http request
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 02 Jan 2021 11:37:20 GMT
Content-Length: 75

{
  "key": "yandex",
  "url": "https://yandex.ru",
  "code": 301,
  "autogenerated": false
}
```

You can also set your own response code (for example 302 temporary redirect):

Request
```http request
POST /urls/ HTTP/1.1
Content-Type: application/json
Authorization: Bearer asdf
Host: 127.0.0.1:80
Content-Length: 63

{
  "key": "yandex",
  "url": "https://yandex.ru",
  "code": 302
}
```

Response
```http request
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 02 Jan 2021 11:42:29 GMT
Content-Length: 75

{
  "key": "yandex",
  "url": "https://yandex.ru",
  "code": 302,
  "autogenerated": false
}
```

#### Generic key POST request

If you don't need a specified key for your link, you can use adding link with generic key.
It's an increasing string.

In this case there will be no new key if database already has a key-url pair with the url you want to add.
For example, if there is a generic key `b` for `example.com`, and you send a new generic POST
request for `example.com`, there result will be `b`, not `c`.

Example:

```http request
POST /urls/ HTTP/1.1
Content-Type: application/json
Authorization: Bearer asdf
Host: 127.0.0.1:80
Content-Length: 25

{
  "url": "example.com"
}
```
```http request
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 02 Jan 2021 15:19:52 GMT
Content-Length: 63

{
  "key": "b",
  "url": "example.com",
  "code": 301,
  "autogenerated": true
}
```

The next generic key request:

```http request
POST /urls/ HTTP/1.1
Content-Type: application/json
Authorization: Bearer asdf
Host: 127.0.0.1:80
Content-Length: 25

{
  "url": "example.com"
}
```
```http request
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 02 Jan 2021 15:21:21 GMT
Content-Length: 63

{
  "key": "b",
  "url": "example.com",
  "code": 301,
  "autogenerated": true
}
```

Make attention that you can understand if the key is generic by looking at `autogenerated` field.

If there is a key conflict with new generic key and existing not generic key, this value
will be skipped, and the next alphabet-ordered free generic key will be used.

For example, if the last generic value was `b` and there is already a `c` key, it 
will be skipped. The next generic key will be `d` if it is free.

#### GET request

```http request
GET /urls/yandex/ HTTP/1.1
Content-Type: application/json
Authorization: Bearer asdf
Host: 127.0.0.1:80
```
```http request
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 02 Jan 2021 15:30:43 GMT
Content-Length: 75

{
  "key": "yandex",
  "url": "https://yandex.ru",
  "code": 302,
  "autogenerated": false
}
```

#### PATCH/PUT request

```http request
PATCH /urls/yandex/ HTTP/1.1
Content-Type: application/json
Authorization: Bearer asdf
Host: 127.0.0.1:80
Content-Length: 27

{
  "url": "www.yandex.ru"
}
```
```http request
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 02 Jan 2021 15:32:16 GMT
Content-Length: 71

{
  "key": "yandex",
  "url": "www.yandex.ru",
  "code": 302,
  "autogenerated": false
}
```

#### DELETE request
```http request
DELETE /urls/yandex/ HTTP/1.1
Authorization: Bearer asdf
Host: 127.0.0.1:80
```
```http request
HTTP/1.1 200 OK
Date: Sat, 02 Jan 2021 15:34:53 GMT
Content-Length: 0
```

### Authorization

For simplicity the project uses a simple Bearer token auth, which you can specify with
`GOSHORT_TOKEN` env var or in config file. The token in single, so make sure that it is
really secret.

If no token is specified, you don't have to write anything to `Authorization` header,
the program won't look to it.
If the token is specified, but something is wrong with auth in your request, you'll get
error 401 (Unauthorized).