package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
)

func main() {
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	certPem := flag.String("cert", "certs/server_chain.pem", "fullchain certificate")
	certKey := flag.String("key", "certs/server.key", "certificate key")
	flag.Parse()
	goproxyCa, err := tls.LoadX509KeyPair(*certPem, *certKey)
	if err != nil {
		log.Fatalln("failed to load certs:", err)
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		log.Fatalln("failed to parse certs:", err)
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("baidu.*:443$"))).
		HandleConnect(goproxy.AlwaysReject)
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
		HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*:80$"))).
		HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
			defer func() {
				if e := recover(); e != nil {
					ctx.Logf("error connecting to remote: %v", e)
					client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
				}
				client.Close()
			}()
			clientBuf := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
			remote, err := net.Dial("tcp", req.URL.Host)
			if err != nil {
				log.Println("failed to connect on", req.URL.Host, ":", err)
			}
			client.Write([]byte("HTTP/1.1 200 Ok\r\n\r\n"))
			remoteBuf := bufio.NewReadWriter(bufio.NewReader(remote), bufio.NewWriter(remote))
			for {
				req, err := http.ReadRequest(clientBuf.Reader)
				if err != nil {
					log.Println("failed to read request:", err)
					return
				}
				if err := req.Write(remoteBuf); err != nil {
					log.Println("failed to write request:", err)
					return
				}
				if err := remoteBuf.Flush(); err != nil {
					log.Println("failed to flush remote buffer:", err)
					return
				}
				resp, err := http.ReadResponse(remoteBuf.Reader, req)
				if err != nil {
					log.Println("failed to read response:", err)
					return
				}
				if err := resp.Write(clientBuf.Writer); err != nil {
					log.Println("failed to write response:", err)
					return
				}
				if err := clientBuf.Flush(); err != nil {
					log.Println("failed to flush client buffer:", err)
					return
				}
			}
		})
	ip := getLocalIP()
	log.Printf("Starting proxy server on %s:8080", ip)
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		local := conn.LocalAddr().(*net.UDPAddr)
		return local.IP.String()
	}
	return ""
}
