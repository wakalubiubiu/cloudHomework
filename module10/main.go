package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"math/rand"
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
	Register()
	mux := http.NewServeMux()
	mux.HandleFunc("/index", index)
	mux.HandleFunc("/healthz", healthz)
	mux.Handle("/metrics", promhttp.Handler())
	srv := http.Server{
		Addr:    "0.0.0.0:" + os.Getenv("GO_PORT"),
		Handler: mux,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Info("listen: %s\n", err)
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

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200")
	w.WriteHeader(200)
}

func index(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("index start")
	timer := NewTimer()
	defer timer.ObserveTotal()
	for k, v := range r.Header {
		w.Header().Set(k, v[0])
	}
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
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
