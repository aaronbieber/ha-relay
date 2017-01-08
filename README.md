# Home Automation Relay #

This is a server/client pair written in Go that is designed to act as a relay
for home automation requests. This will allow your local network to receive
commands through a dedicated TCP pipe without having to run a local web server
open to the Internet.

## Why? ##

The reason for having a thing like this is so that you can use Google Assistant
or Amazon Alexa, which operate in the cloud themselves or through a service like
IFTTT Maker, to receive API requests on some public web server somewhere, then
relay those requests back into your home network through a single dedicated TCP
connection; you do not have to forward ports or run a listening HTTP server on
your local network.

This should be a more secure way to receive messages from Internet services into
your home network.

## How? ##

Simply run the server component on a public server somewhere. Specify the
interface to bind to and the port to listen on. It's probably a good idea to run
this within Supervisor or a similar process management system.

From your local network, run the client component. Specify the server IP and
port similarly.

The client will attempt to remain connected at all times, retrying every 10
seconds if disconnected. When a command is received, take some action on
it. You'll have to modify the code to do this.

The server expects to receive commands through a named pipe. The easiest way to
receive commands is to run the sister Python application, `ha-web-server.py`,
which is a Flask HTTP server that writes to this pipe.

## Caveats ##

This is provided as-is, has barely been tested, requires a ton of manual work to
get going, and might destroy everything that you hold dear. I make no guarantees
or even casual claims that this works, or is even decent code.

So, be careful, friends.
