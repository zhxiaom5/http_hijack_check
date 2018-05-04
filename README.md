# http_hijack_check
## 简介
http_hijack_check是使用go写的一个小工具，可以修改http请求的ttl值。通过ttl值的逐渐增大，可以大致判断劫持发生的位置

## 使用说明
 ./http_hijack_check   -h
flag needs an argument: -h
Usage of ./http_hijack_check:
  -a string
        useragent (default "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.119 Safari/537.36")
  -d    debug this program
  -h string
        hostname (default "www.4399.com")
  -o int
        http get timeout (default 3)
  -p int
        http port (default 80)
  -r string
        http refer
  -s string
        serverip (default "127.0.0.1")
  -t int
        maxttl (default 64)
  -u string
        the hijacked url (default "/index.html")