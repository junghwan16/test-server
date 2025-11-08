package health

import (
	"encoding/json"
	"net/http"
)

// Handler는 헬스체크 HTTP 엔드포인트를 처리합니다.
type Handler struct {
	service *Service
}

// NewHandler는 새로운 헬스체크 핸들러를 생성합니다.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// HandleLiveness는 liveness probe 엔드포인트를 처리합니다.
func (h *Handler) HandleLiveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// HandleReadiness는 readiness probe 엔드포인트를 처리합니다.
func (h *Handler) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	if h.service.CheckReadiness(r.Context()) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(map[string]string{"status": "unavailable"})
}
