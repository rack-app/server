# Project Turtle

## HTTP CONNECT tunneling

Suppose client wants to use either HTTPS or WebSockets in order to talk to server. Client is aware of using proxy. Simple HTTP request / response flow cannot be used since client needs to e.g. establish secure connection with server (HTTPS) or wants to use other protocol over TCP connection (WebSockets). Technique which works is to use HTTP CONNECT method. It tells the proxy server to establish TCP connection with destination server and when done to proxy the TCP stream to and from the client. This way proxy server won’t terminate SSL but will simply pass data between client and destination server so these two parties can establish secure connection.

![](/readme/http_connect_tunneling.png)
