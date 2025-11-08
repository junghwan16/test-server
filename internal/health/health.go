package health

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// Service는 애플리케이션의 헬스체크를 담당합니다.
type Service struct {
	db     *gorm.DB
	logger *slog.Logger
	ready  bool // 애플리케이션 시작 완료 여부
}

// NewService는 새로운 헬스체크 서비스를 생성합니다.
func NewService(db *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
		ready:  false,
	}
}

// MarkReady는 애플리케이션이 시작을 완료했음을 표시합니다.
func (s *Service) MarkReady() {
	s.ready = true
}

// CheckReadiness는 애플리케이션이 트래픽을 받을 준비가 되었는지 확인합니다.
func (s *Service) CheckReadiness(ctx context.Context) bool {
	if !s.ready {
		return false
	}

	return s.checkDatabase(ctx)
}

// checkDatabase는 데이터베이스 연결 상태를 확인합니다.
func (s *Service) checkDatabase(ctx context.Context) bool {
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	sqlDB, err := s.db.DB()
	if err != nil {
		s.logger.Warn("failed to get database instance", "error", err)
		return false
	}

	if err := sqlDB.PingContext(pingCtx); err != nil {
		s.logger.Warn("database health check failed", "error", err)
		return false
	}

	return true
}
