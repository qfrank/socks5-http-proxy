package main

import (
	"flag"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
)

func main() {
	// Parse the command line arguments
	socksAddr := flag.String("socks", "127.0.0.1:1086", "SOCKS5 server address")
	port := flag.Int("port", 8080, "HTTP proxy server port")
	flag.Parse()

	// Create a new SOCKS5 dialer
	dialer, err := proxy.SOCKS5("tcp", *socksAddr, nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	// Create a Dial function that uses the SOCKS5 dialer
	dial := func(network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}

	// Create a new HTTP proxy
	httpProxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			// Do not set the SOCKS5 proxy for the request
			//r.URL.Scheme = "socks5"
			//r.URL.Host = *socksAddr
		},
		Transport: &http.Transport{
			Dial: dial,
		},
	}

	// Create a new HTTP handler that supports the CONNECT method
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			// Get the target host and port
			host, port, err := net.SplitHostPort(r.URL.Host)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Dial the target host and port
			conn, err := dial("tcp", net.JoinHostPort(host, port))
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			defer conn.Close()

			// Respond to the CONNECT request
			w.WriteHeader(http.StatusOK)
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

			// Hijack the connection
			hj, ok := w.(http.Hijacker)
			if !ok {
				http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
				return
			}
			clientConn, _, err := hj.Hijack()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer clientConn.Close()

			// Copy data between the client and the target host
			go func() {
				defer conn.Close()
				io.Copy(conn, clientConn)
			}()
			io.Copy(clientConn, conn)
		} else {
			// Use the HTTP proxy to handle the request
			httpProxy.ServeHTTP(w, r)
		}
	})

	// Start the HTTP proxy server
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), handler))
}
