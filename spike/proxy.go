package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func handleRequest(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		handleTunneling(rw, req)
	} else {
		handleHTTP(rw, req)
	}
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func New(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: http.HandlerFunc(handleRequest),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
}

// func ListenAndServe() {
// 	var pemPath string
// 	flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")
// 	var keyPath string
// 	flag.StringVar(&keyPath, "key", "server.key", "path to key file")
// 	var proto string
// 	flag.StringVar(&proto, "proto", "https", "Proxy protocol (http or https)")
// 	flag.Parse()
// 	if proto != "http" && proto != "https" {
// 		log.Fatal("Protocol must be either http or https")
// 	}
// 	if proto == "http" {
// 		log.Fatal(server.ListenAndServe())
// 	} else {
// 		log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
// 	}
// }
