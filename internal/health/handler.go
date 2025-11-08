package health

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type Handler struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewHandler(db *gorm.DB, logger *slog.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: logger,
	}
}

func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	sqlDB, err := h.db.DB()
	if err != nil {
		h.logger.Error("health check: failed to get db instance", "error", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		h.logger.Warn("health check: database not ready", "error", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
