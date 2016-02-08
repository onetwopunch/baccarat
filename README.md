# Baccarat

A command line utility used to issue shell background tasks on any platform. Works by reading command line args and callback URI's from redis, executes the command and sends the output to the URI as a POST request with a status of 1 or 0 and the full stdout or error as a string.

You can configure the thread pool size as well as any headers you'll need for the callback URI

## Install
go get -u github.com/onetwopunch/baccarat

## Usage

baccarat -q default -c config.ini --pool 5
