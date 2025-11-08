package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

var (
	isReady     int32 = 1
	isLive      int32 = 1
	isStartupOk int32 = 1
)

func main() {
	http.HandleFunc("/startup", handleStartup)
	http.HandleFunc("/ready", handleReady)
	http.HandleFunc("/live", handleLive)

	// 부하 ON/OFF 및 관련 API
	http.HandleFunc("/server-load-on", handleServerLoadOn)
	http.HandleFunc("/server-load-off", handleServerLoadOff)

	addr := ":8080"
	fmt.Println("Starting server on", addr)
	http.ListenAndServe(addr, nil)
}

// 앱 초기화 체크
func handleStartup(w http.ResponseWriter, r *http.Request) {
	log.Println("startup probe called")

	// startupProbe 실패 강제 (응용 1)
	if atomic.LoadInt32(&isStartupOk) == 0 {
		http.Error(w, "startup failing (forced by /startup-fail-on)", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "startup ok")
}

// DB 연결, 캐시 로드 등 준비 상태 체크
func handleReady(w http.ResponseWriter, r *http.Request) {
	log.Println("ready probe called")

	if atomic.LoadInt32(&isReady) == 0 {
		http.Error(w, "ready=false", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ready=true")
}

// 단순 핑 응답 (일반적인 /healthz 역할)
func handleLive(w http.ResponseWriter, r *http.Request) {
	log.Println("live probe called")

	if atomic.LoadInt32(&isLive) == 0 {
		http.Error(w, "live=false", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "live=true")
}

// 부하 증가
func handleServerLoadOn(w http.ResponseWriter, r *http.Request) {
	atomic.StoreInt32(&isReady, 0)
	atomic.StoreInt32(&isLive, 0)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Server load ON. isReady=false, isLive=false. With proper probe settings, traffic will be stopped soon and the app will restart later.")
}

// 부하 감소
func handleServerLoadOff(w http.ResponseWriter, r *http.Request) {
	atomic.StoreInt32(&isReady, 1)
	atomic.StoreInt32(&isLive, 1)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Server load OFF. isReady=true, isLive=true")
}
