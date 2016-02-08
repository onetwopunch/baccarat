# Baccarat
## Concurrent and distributable background task scheduler usable on any platform


*Disclaimer: This is still being worked on and should not be used in production yet*

Works by reading command line args and callback URI's from Redis, executes the command and sends the output to a URI as a POST request with a status of 1 or 0 and the full stdout or error as a string.

You can configure the thread pool size as well as any headers you'll need for the callback URI

## Install

    go get -u github.com/onetwopunch/baccarat

## Usage

    baccarat -q default -c config.ini

Where the `config.ini` looks like this:

```
[default]
pool = 5

[redis]
url = localhost:6379
password = xyz

[headers]
Authorization = TOKEN xyz
X-Baccarat = Some Value
```

You can have whatever headers you want. The POST request will be sent out with a JSON payload, but you can add any Authorization or Signature you like.

Then just call the command line to store in Redis:

    baccarat write <queue name> <command> <uri> -c <path to config>
    baccarat write default 'echo hello world' http://localhost:8080 -c config.ini

Behind the scenes, it's just calling Redis `LPUSH` so you could do this from any Redis client. Just prefix your queue name with `baccarat:` namespace.

    redis-cli> LPUSH baccarat:default '{"Command":"echo hello world","URI":"http://localhost:8080"}'

## Development


    git clone git@github.com:onetwopunch/baccarat.git
    gopm Install

Then to test, you'll need three terminal windows open:

Listener:

    ./vendor/bin/baccarat -q default -c config.ini

Node Server:

    node test/test_server.js

Redis Control:

    .vendor/bin/baccarat write default 'echo hello world' http://localhost:8080 -c config.ini
