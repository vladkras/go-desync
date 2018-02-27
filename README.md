# Request Desynchronizer

## What is "desynchronizer"?
Modern apps make multiple HTTP request during their life-cycle: loading data, querying databases, sending notifications, calling APIs. All this actions together require a decent amount of time and may slow down or even completely freeze your app if some services are broken or unreachable.

You can avoid delays and shutdowns by delegating all request to desynchronizer. It will send your request asynchronously and return response when ready (if needed).

## Usage
1. Run `./go-desync [--debug] [--port=8080]`

   desynchronizer will start server (port 8080 by default). Add `--debug` (false by default) flag too see all output
2. Prepend your requests with desynchronizer endpoint like:

   http://localhost:8080/https://your-request.url/what-ever


 That's it. Desynchronizer uses your original request method and payload, and sends request when ready. It prints both request and response to standard error if `--debug` flag is passed.

 Desynchronizer **always** returns `200 OK` status even if your path or data is senseless.

## TODO
 * Callback url
 * Additional custom headers: (retry, ttl, etc.)

## License
This project is licensed under [MIT License](https://github.com/vladkras/go-desync/blob/master/LICENSE) and developed by &copy; 2018, GraphitLab R&D
