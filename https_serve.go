/* http://www.kaihag.com/https-and-go/ */

package main

/* 
 * # to generate server.key and server.pem
 * openssl genrsa -out server.key 2048
 * openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
*/

/*
 * Port Forwarding on Router
 *
 * 80 -> 9090    http on 9090 will redirect
 * 443 -> 4430   so that https://fpkoehler.dyndns.org works
 * 4430 -> 4430  because we redirect to https://fpkoehler.dyndns.org:4430
 *
 * External                Local
 * IP_Addr Start/End Port  IP Address   Start/End Port  Description     Protocol
 * 0.0.0.0  443 / 443      192.168.0.18 4430 / 4430     ubuntu_https    TCP
 * 0.0.0.0   80 / 80       192.168.0.18 9090 / 9090     ubuntu_http     BOTH
 * 0.0.0.0 4430 / 4430     192.168.0.18 4430 / 4430     ubuntu_4430     TCP 
*/

import (
    "fmt"
    "log"
    "net/http"
)

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
    // Redirect the incoming HTTP request.
    fmt.Println("redirectToHttps", r.RequestURI)
//    http.Redirect(w, r, "https://localhost:4430"+r.RequestURI, http.StatusFound)
    http.Redirect(w, r, "https://fpkoehler.dyndns.org:4430"+r.RequestURI, http.StatusFound)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there!")
}

func main() {
    http.HandleFunc("/", handler)

    // Start the HTTPS server in a goroutine
    go func() {
        err := http.ListenAndServeTLS(":4430", "server.pem", "server.key", nil)
        if err != nil {
            log.Fatalf("ListenAndServeTLS error: %v", err)
        }
    } ()

    // Start the HTTP server and redirect all incoming connections to HTTPS
//    err := http.ListenAndServe(":8080", http.HandlerFunc(redirectToHttps))
    err := http.ListenAndServe(":9090", http.HandlerFunc(redirectToHttps))
    if err != nil {
        log.Fatalf("ListenAndServe error: %v", err)
    }
}
