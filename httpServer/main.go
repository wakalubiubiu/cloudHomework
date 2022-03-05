package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	flag.Set("v", "4")
	err1 := flag.Set("logtostderr", "true")
	if err1 != nil {
		fmt.Println(err1)
	}
	http.HandleFunc("/healthz", log(healthz))
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		glog.Fatal("run failed, error: ", err.Error())
	}
}

func log(handler func(w http.ResponseWriter, r *http.Request) int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var status int
		defer func(w http.ResponseWriter, r *http.Request) {
			glog.Infof("ip:%s, status code:%s", getIP(r), status)
		}(w, r)
		defer glog.Flush()
		for k, v := range r.Header {
			w.Header().Set(k, v[0])
		}
		w.Header().Set("version", os.Getenv("VERSION"))
		status = handler(w, r)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) int {
	io.WriteString(w, "200")
	w.WriteHeader(200)
	return 200
}

func getIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
