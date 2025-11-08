package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
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

	// Hello API
	http.HandleFunc("/hello", handleHello)

	// 응용 1: startupProbe 실패 유발
	http.HandleFunc("/startup-fail-on", func(w http.ResponseWriter, r *http.Request) {
		atomic.StoreInt32(&isStartupOk, 0)
		fmt.Fprintln(w, "Set to fail startupProbe. The pod will restart indefinitely if a startupProbe is configured.")
	})

	// 응용 2: 부하 ON/OFF 및 관련 API
	http.HandleFunc("/server-load-on", handleServerLoadOn)
	http.HandleFunc("/server-load-off", handleServerLoadOff)

	addr := ":8080"
	fmt.Println("Starting server on", addr)
	http.ListenAndServe(addr, nil)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	greetings := []string{
		"Hello, World!",
		"Live long and prosper.",
		"May the Force be with you.",
		"So long, and thanks for all the fish.",
		"Never give up, never surrender!",
		"Go forth, and multiply.",
		"To infinity, and beyond!",
		"Excelsior!",
	}

	rand.Seed(time.Now().UnixNano())
	greeting := greetings[rand.Intn(len(greetings))]

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, greeting)
}

// 앱 초기화 체크
func handleStartup(w http.ResponseWriter, r *http.Request) {
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
