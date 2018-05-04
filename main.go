package main

//http hijack check
//use ip ttl to check  http hijack

import (
	"bufio"
	"flag"
	"fmt"
	//"io"
	"bytes"
	"golang.org/x/net/ipv4"
	"net"
	//"os"
	"strconv"
	"time"
)

var (
	serverip string
	port     int
	url      string
	hostname string
	//useragent
	ua    string
	refer string
	//max ip ttl
	maxttl int
	debug  bool
	//full raw http get
	rawhttpget string
	timeout    int
)

func init() {
	flag.StringVar(&serverip, "s", "127.0.0.1", "serverip")
	flag.IntVar(&port, "p", 80, "http port")
	flag.IntVar(&maxttl, "t", 64, "maxttl")
	flag.IntVar(&timeout, "o", 3, "http get timeout")
	flag.StringVar(&url, "u", "/index.html", "the hijacked url")
	flag.StringVar(&hostname, "h", "www.4399.com", "hostname")
	flag.StringVar(&ua, "a", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.119 Safari/537.36", "useragent")
	flag.StringVar(&refer, "r", "", "http refer")
	flag.BoolVar(&debug, "d", false, "debug this program")

}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func hijack_http(ttl int, ch chan bool) {

	conn, err := net.DialTimeout("tcp", serverip+":"+strconv.Itoa(port), time.Second*5)
	if err != nil {
		fmt.Println(err)
	}

	if err := ipv4.NewConn(conn).SetTTL(ttl); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n\nnow ttl is %d,maxttl is %d \n", ttl, maxttl)

	//发起请求
	fmt.Fprintf(conn, rawhttpget)
	scanner := bufio.NewScanner(conn)
	scanner.Split(ScanCRLF)
	fmt.Println("------------------------------------")
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if scanner.Text() == "" {
			ch <- true
			break
		}
	}

}

func main() {
	flag.Parse()
	//拼接生成raw http get
	//fmt.Fprintf(conn, "GET /union2/channel/8166/1465802583522770/gsdzz-0.1.30.1.1-b-n-8166.apk HTTP/1.1\r\nHost: sy-cdnres.unionsy.com\r\nX-Real-IP: 127.0.0.1\r\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.119 Safari/537.36\r\n\r\n")
	rawhttpget = "GET " + url + " HTTP/1.1\r\nHost: " + hostname + "\r\nUser-Agent: " + ua + "\r\nReferer: " + refer + "\r\n\r\n"

	fmt.Printf("server info is %s:%d \n", serverip, port)
	fmt.Printf("raw http get info is %s \n", rawhttpget)
	for i := 1; i <= maxttl; i++ {
		ch := make(chan bool, 1)

		go hijack_http(i, ch)

		select {
		case <-ch:
			fmt.Printf("test finished\n")
		case <-time.After(time.Second * time.Duration(timeout)):
			fmt.Printf("test timeout after %d second\n", timeout)
		}
	}

}
