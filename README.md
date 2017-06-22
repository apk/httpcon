Try to expose the local ssh via HTTP CONNECT (later)
on port 443, and write a C client using openssl.

The point is to have a ssh `ProxyCommand` that first
uses HTTP CONNECT on a proxy to connect to the actual
server, then performs the ssl setup, and then runs
another HTTP CONNECT inside the ssl connection to
connect to the actual ssh port (or one of them).

This would allow to use a single IP addess and a
HTTPS server to accept connections to multiple
ssh (and other) services, all while acting as
a regular web server as well. (But the server
side would need to directly expose the go web
server, as not many reverse proxies will forward
and handle the CONNECT in the way required here.)

Nowadays, the proper way to do this would be
via websockets (https://github.com/yydesa/wscat)
but sometimes the client side is so ancient that
you simply have no websocket implementation available.

`httpcon.go`: work in progress, for doing the CONNECT
(w/o ssl yet).

`sslconnect.c`: sample openssl client code, snarfed from the 'net.
May be working.
