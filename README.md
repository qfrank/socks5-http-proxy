Overview

This is a Go program that creates an HTTP proxy server that uses a SOCKS5 server as a proxy for incoming connections. It is capable of handling both regular HTTP requests and HTTPS requests using the CONNECT method.

Running the Proxy Server

To start the HTTP proxy server, use the following command:
go run socks5-http-proxy.go -socks <socks-server-address> -port <proxy-server-port>

For example, to start the proxy server on localhost port 8080 and use a SOCKS5 server at 192.168.0.188:1086 as the proxy, use the following command:
go run socks5-http-proxy.go -socks 192.168.0.188:1086 -port 8080

Once the proxy server is running, clients can configure their HTTP client to use the proxy server by specifying the proxy server address and port. For example, in the curl command line utility, the -x flag can be used to specify the proxy server, like this:
curl -x http://localhost:8080 http://example.com

This command will send a request to http://example.com using the proxy server. If the proxy server is configured and running correctly, the response will be sent back to the client via the proxy server and the SOCKS5 server.

Conclusion

This Go program creates an HTTP proxy server that uses a SOCKS5 server as a proxy to forward requests to target hosts. It can be started using the go run command and configured using command line arguments. Once running, clients can use the proxy server by specifying the proxy server address and port in their HTTP client settings.
