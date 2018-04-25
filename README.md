Send HTTP-requests without waiting.

![coverage: 33%](https://img.shields.io/badge/coverage-33%25-yellow.svg)

## What is "desynchronizer"?
Modern apps make multiple HTTP request during their life-cycle: loading data, querying databases, sending notifications, calling APIs. All this actions together require a decent amount of time and may slow down or even completely freeze your app if some services are broken or unreachable.

You can avoid delays and shutdowns by delegating all request to desynchronizer. It will send your request asynchronously and return response when ready (if needed).

## Usage
1. Run `./go-desync [--debug] [--port=8080] [--cert=""]`

   desynchronizer will start server (port 8080 by default). Add `--debug` (false by default) flag too see all output
2. Prepend your requests with desynchronizer endpoint like:

   http://localhost:8080/https://your-request.url/what-ever


 That's it. Desynchronizer uses your original request method and payload, and sends request when ready. It prints both request and response to standard error if `--debug` flag is passed.

 Desynchronizer **always** returns `200 OK` status even if path, data or response are senseless.

## Docker

Thanks to [namsral/flag](https://github.com/namsral/flag) package use flags for CLI (`--debug`, `--port 8080`, `--cert=/var/ssl`) or uppercase `--env` for Docker (`PORT=8080`, `DEBUG=true`, `CERT=/var/ssl`). Usually, this looks like

`docker run -d -p 8080:8080 vladkras/go-desync`

or with TLS/SSL support

`docker run -d -p 8080:8080 -v "/path/to/certs:/var/ssl" -e CERT=/var/ssl -e DEBUG=true vladkras/go-desync`

! Do not use `/etc/ssl` - it's already used for [certificates for request](https://github.com/vladkras/go-desync/blob/master/Dockerfile#L3).

## HTTPS

Define `cert` flag or `CERT` environment variable as path to your \*.crt and \*.key files. If found they will be checked and used for secured server.

## TODO
 * Callback url
 * Additional custom headers: (retry, ttl, etc.)

## License
This project is licensed under [MIT License](https://github.com/vladkras/go-desync/blob/master/LICENSE) and developed by &copy; 2018, GraphitLab R&D

This project uses [namsral/flag](https://github.com/namsral/flag) package licensed under [BSD 3-Clause License](https://github.com/namsral/flag/blob/master/LICENSE)
