// run-on-change httpcon.go -- g build httpcon.go -- ./httpcon

package main

import (
	"net/http"
	"flag"
	"log"
	"net"
	"io"
	"time"
	"bufio"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler %s %q %q", r.Method, r.Host, r.URL.Path)
	w.Write([]byte("hi\n"))
}

type ErrChan chan error

func cls(out net.Conn) {

	tcp, ok := out.(net.TCPConn)
	if ok {
		tcp.CloseWrite()
		log.Printf ("CloseWrote")
	} else {
		log.Printf ("Not a tcp conn")
	}
	
}

func copy(fin ErrChan, in, out net.Conn) {
	//log.Printf("Beginning incopy")
	_, err := io.Copy(out, in)
	//log.Printf("Done incopy %v", err)
	cls(out)
	fin <- err
}

func copy_(fin ErrChan, in, out net.Conn, rdr *bufio.ReadWriter) {
	//log.Printf("Beginning outcopy")
	_, err := io.Copy(out, rdr)
		//log.Printf("Done bufcopy %v", err)
	if err == nil {
		_, err = io.Copy(out, in)
		//log.Printf("Done outcopy %v", err)
	}
	cls(out)
	fin <- err
}

func there(w http.ResponseWriter, r *http.Request) {
	log.Printf("hosthandler %s %q %q", r.Method, r.Host, r.URL.Path)

	hij, ok := w.(http.Hijacker)

	if !ok {
		http.NotFound(w, r)
		return
	}

	upconn, reader, err := hij.Hijack()
	if err != nil {
		http.Error(w, err.Error(), 516)
		return
	}
	defer upconn.Close()

	var t time.Time
	upconn.SetDeadline(t)

	conn, err := net.Dial("tcp", "localhost:22")
	if err != nil {
		log.Println("connect:", err)
		http.Error(w, err.Error(), 517)
		return
	}
	defer conn.Close()

	upconn.Write([]byte("HTTP/1.0 200 All is good.\r\n\r\n"))

	fin := make(chan error, 2)

	go copy(fin, conn, upconn)

	go copy_(fin, upconn, conn, reader)

	err = <-fin
	// Also wait for timeouts here, esp do timeout on the second close.
	// Apparently a stupid idea.
	// EOF-cleanliness would be good,
	// but seems unattainable.
	//	if err == nil {
	//		log.Printf("First here, waiting for second")
	//		err = <-fin
	//	}
	log.Printf("Conn done")

}

func main() {
	mux := http.NewServeMux()
	flag.Parse()
	log.Printf("Up")

	mux.HandleFunc("/", handler)
	mux.HandleFunc("there", there)
	mux.HandleFunc("localhost", there)
	mux.HandleFunc("localhost:443", there)
	mux.HandleFunc("127.0.0.1:443", there)
	log.Fatal(http.ListenAndServe(":4004", mux))
}
