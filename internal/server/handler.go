package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func HandleLive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleReady(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		sqlDB, err := db.DB()
		if err != nil {
			logger.Error("health check: failed to get db instance", "error", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}

		if err := sqlDB.PingContext(ctx); err != nil {
			logger.Warn("health check: database not ready", "error", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func HandleHealth(db *gorm.DB, logger *slog.Logger) http.HandlerFunc {
	type health struct {
		Status     string            `json:"status"`
		Timestamp  time.Time         `json:"timestamp"`
		Components map[string]string `json:"components,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		h := health{
			Status:     "healthy",
			Timestamp:  time.Now(),
			Components: make(map[string]string),
		}

		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(ctx) != nil {
			h.Status = "unhealthy"
			h.Components["database"] = "unhealthy"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(h)
			return
		}

		h.Components["database"] = "healthy"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(h)
	}
}
