package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	flag.Set("v", os.Getenv("GO_LEVEL"))
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("healthz", healthz)
	srv := http.Server{
		Addr:    "0.0.0.0:" + os.Getenv("GO_PORT"),
		Handler: mux,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("listen: %s\n", err)
		}
	}()
	glog.V(4).Info("Server Started")
	<-done
	glog.V(4).Info("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		glog.Fatalf("Server Shutdown Failed:%+v", err)
	}
	glog.V(4).Info("Server Exited Properly")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200")
	w.WriteHeader(200)
}

func index(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("index start")
	for k, v := range r.Header {
		w.Header().Set(k, v[0])
	}
	w.Header().Set("version", os.Getenv("VERSION"))
	io.WriteString(w, "welcome to index!")
	w.WriteHeader(200)
	glog.V(4).Infof("ip:%s, status code:%d", getIP(r), 200)
	glog.V(4).Info("index end")
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
