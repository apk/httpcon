Try to expose the local ssh via HTTP CONNECT (later)
on port 443, and write a C client using openssl.

`httpcon.go`: work in progress, for doing the CONNECT
(w/o ssl yet).

`sslconnect.c`: sample openssl client code, snarfed from the 'net.
May be working.
