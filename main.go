package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	isReady int32 = 1
	isLive  int32 = 1
	start         = time.Now()
)

func main() {
	http.HandleFunc("/startup", handleStartup)
	http.HandleFunc("/ready", handleReady)
	http.HandleFunc("/live", handleLive)

	http.HandleFunc("/server-load-on", handleLoadOn)
	http.HandleFunc("/server-load-off", handleLoadOff)

	addr := ":8080"
	fmt.Println("Starting server on", addr)
	http.ListenAndServe(addr, nil)
}

// 앱 초기화 체크
func handleStartup(w http.ResponseWriter, r *http.Request) {
	// startupProbe 실패 강제
	if os.Getenv("STARTUP_ALWAYS_FAIL") == "true" {
		http.Error(w, "startup failing (forced)", http.StatusInternalServerError)
		return
	}

	// 선택적 지연 (STARTUP_DELAY_SECONDS)
	if d := os.Getenv("STARTUP_DELAY_SECONDS"); d != "" {
		var delay int
		fmt.Sscanf(d, "%d", &delay)
		if time.Since(start) < time.Duration(delay)*time.Second {
			http.Error(w, "starting up", http.StatusServiceUnavailable)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "startup ok")
}

// DB 연결, 캐시 로드 등 준비 상태 체크
func handleReady(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&isReady) == 0 {
		http.Error(w, "ready=false", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ready=true")
}

// 단순 핑 응답 (일반적인 /healthz 역할)
func handleLive(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&isLive) == 0 {
		http.Error(w, "live=false", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "live=true")
}

// 가상의 부하 상태를 가정
// isReady, isLive 를 랜덤으로 변경
func handleLoadOn(w http.ResponseWriter, r *http.Request) {
	// 30초 뒤 ready=false, 3분 뒤 live=false
	go func() {
		time.Sleep(30 * time.Second)
		atomic.StoreInt32(&isReady, 0)
	}()
	go func() {
		time.Sleep(3 * time.Minute)
		atomic.StoreInt32(&isLive, 0)
	}()
	fmt.Fprintln(w, "load-on: ready->false@30s, live->false@180s")
}

func handleLoadOff(w http.ResponseWriter, r *http.Request) {
	atomic.StoreInt32(&isReady, 1)
	atomic.StoreInt32(&isLive, 1)
	fmt.Fprintln(w, "load-off: ready/live restored")
}
