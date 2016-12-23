package main

import (
  "fmt"
  "net/http"
  "net"
  "time"
  "log"
  "io/ioutil"
  "os"
  "strconv"
)

var ProtocolHeader string = "X-Protocol"

type ProxyHanlder struct {
    http.ServeMux
}

func PoolDial(network, addr string) (net.Conn, error) {
    dial := net.Dialer{
        Timeout:   70 * time.Second,
        KeepAlive: 70 * time.Second,
    }

    conn, err := dial.Dial(network, addr)
    if err != nil {
        return conn, err
    }

    // fmt.Fprintf(os.Stderr, "connecting %s, local addr is %s\n", addr, conn.LocalAddr().String())

    return conn, err
}

var client *http.Client

func (mux *ProxyHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    startTime := time.Now()
    url := r.URL
    protocol := r.Header.Get(ProtocolHeader)
    originalUrl := url.String()
    if protocol != "" {
        url.Scheme = protocol
        // fmt.Printf("protocol is %q\n", protocol)
    }

    // defer r.Body.Close()
    // body, err := ioutil.ReadAll(r.Body)
    // if err != nil {
    //    panic("error reading request body")
    // }
    req, req_err := http.NewRequest(r.Method, url.String(), r.Body)
    // set headers here
    if req_err != nil {
        // w.Header().Add("Connection", "closed")
        w.WriteHeader(500)
        req_err_msg := fmt.Sprintf("error construct new request: %v\n", req_err)
        fmt.Fprintf(w, req_err_msg)
        fmt.Fprintf(os.Stderr, fmt.Sprintf("%v %s", startTime, req_err_msg))
        fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", startTime, r.RemoteAddr, r.Method, url.String(), 500, len(req_err_msg), (float64)(time.Now().UnixNano() - startTime.UnixNano()) / 1e9, originalUrl)
        return
        // panic("error construct NewRequest")
    }
    req.Header = r.Header
    req.Header.Del("Proxy-Connection")
    if protocol != "" {
        req.Header.Del(ProtocolHeader)
    }

    response, res_err := client.Do(req)
    if res_err != nil {
        // w.Header().Add("Connection", "closed")
        w.WriteHeader(500)
        res_err_msg := fmt.Sprintf("error get response: %v, response is %v\n", res_err, response)
        // fmt.Fprintf(w, res_err_msg)
        fmt.Fprintf(os.Stderr, fmt.Sprintf("%v %s", startTime, res_err_msg))
        fmt.Fprintf(w, res_err_msg)
        fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", startTime, r.RemoteAddr, r.Method, url.String(), 500, len(res_err_msg), (float64)(time.Now().UnixNano() - startTime.UnixNano()) / 1e9, originalUrl)
        // panic("error get response")
        return
    }
    defer response.Body.Close()
    for key, values := range response.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    // w.Header().Del("Connection")
    // w.Header().Add("Connection", "closed")
    w.WriteHeader(response.StatusCode)
    body, _ := ioutil.ReadAll(response.Body)
    w.Write(body)
    fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", startTime, r.RemoteAddr, r.Method, url.String(), response.StatusCode, len(body), (float64)(time.Now().UnixNano() - startTime.UnixNano()) / 1e9, originalUrl)
    // fmt.Printf("headers: %q\n", req.Header)
    // write headers here
    // fmt.Fprintf(w, "Hello, url=%s, response=%q", r.URL, response)

    // fmt.Fprintf(w, "Hello, url=%s, body=%s, w=%q, r=%q\n", r.URL, body, w, r)
    // mux.ServeMux.ServeHTTP(w, r)
}

func main() {
    argc := len(os.Args)
    var listenAddr string
    var maxIdleConnsPerHost int = 10
    if argc == 1 {
        fmt.Fprintf(os.Stderr, "%s [ListenHost]:ListenPort [MaxIdleConnsPerHost=10 [ProtocolHeader = X-Protocol]]\n", os.Args[0]);
        return
    }
    if argc > 1 {
        listenAddr = os.Args[1]
    }
    if argc > 2 {
        maxIdleConnsPerHost, _ = strconv.Atoi(os.Args[2])
    }
    if argc > 3 {
        ProtocolHeader = os.Args[3]
    }
    client = &http.Client{
        Transport: &http.Transport{
            Dial: PoolDial,
            MaxIdleConnsPerHost: maxIdleConnsPerHost,
        },
    }
    log.Fatal(http.ListenAndServe(listenAddr, &ProxyHanlder{}))
}

