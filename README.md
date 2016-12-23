# go\_keep\_alived\_http\_proxy

**a keep-alived http/https proxy written in go**

It is not easy for php(cgi mode) to reuse http connection(s) between different requests(because php almost cleans up resources after each request, you can reuse http connection if you write a php extension in C, but it will take you a lot of time), so a keep-alived http(s) proxy is required to improve performance(by avoiding tcp 3-way handshake, and TLS handshake(only for https) after the first request), you can proxy https request via an http proxy, let's get started:

_build:_

    go build http.go
    
<br>
_start the server:_

    ./http_pool :9999 16 X-Real-Protocol
this will start service listening on :9999 and it's ready to proxy requests
<br>
_proxy http requests:_
    
    curl -i -x 127.0.0.1:9999 www.baidu.com
will output(from curl):

    HTTP/1.1 200 OK
    Cache-Control: private, no-cache, no-store, proxy-revalidate, no-transform
    Connection: Keep-Alive
    Content-Type: text/html
    Date: Fri, 23 Dec 2016 19:27:17 GMT
    Last-Modified: Mon, 25 Jul 2016 11:11:41 GMT
    Pragma: no-cache
    Server: bfe/1.0.8.18
    Set-Cookie: BDORZ=27315; max-age=86400; domain=.baidu.com; path=/
    Transfer-Encoding: chunked
    
    <!DOCTYPE html>
    <!--STATUS OK--><html> <head><meta http-equiv=content-type content=text/html;charset=utf-8><meta http-equiv=X-UA-Compatible content=IE=Edge><meta content=always name=referrer><link rel=stylesheet type=text/css href=http://s1.bdstatic.com/r/www/cache/bdorz/baidu.min.css> ... body ommitted
and output(from http_pool)(or maybe output error to stderr):

    2016-12-24 03:09:09.4616205 +0800 CST   127.0.0.1:19747 GET     http://www.baidu.com/   200     2381    0.0811881       http://www.baidu.com/
    #request_time                           local_ip_port   method  real_request_url status_code conent_length time_spent   http_request_url
    
if any error occurred during the request, it will respond 500 internal server error

    curl -i -x 127.0.0.1:9999 non_exists_website.com
    HTTP/1.1 500 Internal Server Error
    Date: Fri, 23 Dec 2016 19:12:47 GMT
    Content-Length: 153
    Content-Type: text/plain; charset=utf-8
    
    error get response: Get http://non_exists_website.com/: dial tcp: lookup non_exists_website.com: getaddrinfow: No such host is known., response is <nil>
<br>
_proxy https requests:_
you can proxy an https url via an http proxy by setting **X-Real-Protocol** header(the 3rd command line argument) with value **https**:

    curl -H 'X-Real-Protocol: https' -i -x 127.0.0.1:9999 www.baidu.com
    HTTP/1.1 200 OK
    Cache-Control: private, no-cache, no-store, proxy-revalidate, no-transform
    Connection: keep-alive
    Content-Type: text/html
    Date: Fri, 23 Dec 2016 19:27:03 GMT
    Last-Modified: Mon, 25 Jul 2016 11:13:20 GMT
    Pragma: no-cache
    Server: bfe/1.0.8.18
    Set-Cookie: BDORZ=27315; max-age=86400; domain=.baidu.com; path=/
    Set-Cookie: __bsi=12187426595763402941_00_44_N_N_1_0303_C02F_N_N_Y_0; expires=Fri, 23-Dec-16 19:27:08 GMT; domain=www.baidu.com; path=/
    Transfer-Encoding: chunked
    
    <!DOCTYPE html>
    <!--STATUS OK--><html> <head><meta http-equiv=content-type content=text/html;charset=utf-8><meta http-equiv=X-UA-Compatible content=IE=Edge><meta content=always name=referrer><link rel=stylesheet type=text/css href=https://ss1.bdstatic.com/5eN1bjq8AAUYm2zgoY3K/r/www/cache/bdorz/baidu.min.css> ... body omitted
    
**Full Usage is: ./http_pool [ListenHost]:ListenPort [MaxIdleConnsPerHost=10 [ProtocolHeader = X-Protocol]]**

**MaxIdleConnsPerHost** specifies maximum idle connections for a certain host, increasing this number will make it faster to proxy request to the same host simultaneously, because go http client usually(but not always) reuses existing http connection to avoid (expensive) tcp 3-way handshake if target host support HTTP **'Connection: Keep-Alive'**, but if the number is too big, it is a waste of resources because there may be too many unused keep-alive http connections
you can change **X-Real-Protocol** to any value you want by modify the 3rd command line argument(ProtocolHeader)

before you start it, please make sure that target host supports HTTP **'Connection: Keep-Alive'**, now please have a try~

