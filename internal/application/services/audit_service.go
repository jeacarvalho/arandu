package services

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	appcontext "arandu/internal/platform/context"
	"github.com/google/uuid"
)

const (
	AuditActionAccessPatient = "ACCESS_PATIENT"
	AuditActionCreatePatient = "CREATE_PATIENT"
	AuditActionUpdatePatient = "UPDATE_PATIENT"
	AuditActionDeletePatient = "DELETE_PATIENT"
	AuditActionAccessSession = "ACCESS_SESSION"
	AuditActionCreateSession = "CREATE_SESSION"
	AuditActionLogin         = "LOGIN"
	AuditActionLogout        = "LOGOUT"
	AuditActionExportData    = "EXPORT_DATA"
)

type AuditLog struct {
	ID         string
	Timestamp  time.Time
	UserID     string
	TenantID   string
	Action     string
	ResourceID string
	IPAddress  string
	UserAgent  string
}

type AuditService struct {
	db         *sql.DB
	logChan    chan *AuditLog
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	bufferSize int
}

type AuditOption func(*AuditService)

func WithAuditBufferSize(size int) AuditOption {
	return func(s *AuditService) {
		s.bufferSize = size
	}
}

func NewAuditService(db *sql.DB, opts ...AuditOption) *AuditService {
	ctx, cancel := context.WithCancel(context.Background())

	s := &AuditService{
		db:         db,
		logChan:    make(chan *AuditLog, 100),
		ctx:        ctx,
		cancel:     cancel,
		bufferSize: 100,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.wg.Add(1)
	go s.worker()

	return s
}

func (s *AuditService) worker() {
	defer s.wg.Done()

	for {
		select {
		case auditLog := <-s.logChan:
			if auditLog != nil {
				s.persistLog(auditLog)
			}
		case <-s.ctx.Done():
			for len(s.logChan) > 0 {
				auditLog := <-s.logChan
				if auditLog != nil {
					s.persistLog(auditLog)
				}
			}
			return
		}
	}
}

func (s *AuditService) persistLog(auditLog *AuditLog) {
	query := `
		INSERT INTO audit_logs (id, timestamp, user_id, tenant_id, action, resource_id, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(context.Background(), query,
		auditLog.ID,
		auditLog.Timestamp,
		auditLog.UserID,
		auditLog.TenantID,
		auditLog.Action,
		auditLog.ResourceID,
		auditLog.IPAddress,
		auditLog.UserAgent,
	)

	if err != nil {
		log.Printf("❌ Failed to persist audit log: %v", err)
	} else {
		log.Printf("✅ Audit log persisted: action=%s, resource=%s, tenant=%s",
			auditLog.Action, auditLog.ResourceID, auditLog.TenantID)
	}
}

func (s *AuditService) Log(ctx context.Context, action, resourceID string) {
	userID, err := appcontext.GetUserID(ctx)
	if err != nil {
		log.Printf("⚠️ AuditService: skipping log - missing user_id in context")
		return
	}

	tenantID, err := appcontext.GetTenantID(ctx)
	if err != nil {
		log.Printf("⚠️ AuditService: skipping log - missing tenant_id in context")
		return
	}

	auditLog := &AuditLog{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		UserID:     userID,
		TenantID:   tenantID,
		Action:     action,
		ResourceID: resourceID,
	}

	select {
	case s.logChan <- auditLog:
	default:
		log.Printf("⚠️ AuditService: channel full, dropping audit log")
	}
}

func (s *AuditService) LogWithDetails(ctx context.Context, action, resourceID, ipAddress, userAgent string) {
	userID, err := appcontext.GetUserID(ctx)
	if err != nil {
		log.Printf("⚠️ AuditService: skipping log - missing user_id in context")
		return
	}

	tenantID, err := appcontext.GetTenantID(ctx)
	if err != nil {
		log.Printf("⚠️ AuditService: skipping log - missing tenant_id in context")
		return
	}

	auditLog := &AuditLog{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		UserID:     userID,
		TenantID:   tenantID,
		Action:     action,
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	select {
	case s.logChan <- auditLog:
	default:
		log.Printf("⚠️ AuditService: channel full, dropping audit log")
	}
}

func (s *AuditService) Close() error {
	s.cancel()
	s.wg.Wait()
	close(s.logChan)
	return nil
}

func (s *AuditService) GetLogsByTenant(ctx context.Context, tenantID string, limit int) ([]AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `
		SELECT id, timestamp, user_id, tenant_id, action, resource_id, ip_address, user_agent
		FROM audit_logs
		WHERE tenant_id = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		var ipAddress, userAgent sql.NullString
		var resourceID sql.NullString

		err := rows.Scan(&log.ID, &log.Timestamp, &log.UserID, &log.TenantID, &log.Action, &resourceID, &ipAddress, &userAgent)
		if err != nil {
			return nil, err
		}

		if resourceID.Valid {
			log.ResourceID = resourceID.String
		}
		if ipAddress.Valid {
			log.IPAddress = ipAddress.String
		}
		if userAgent.Valid {
			log.UserAgent = userAgent.String
		}

		logs = append(logs, log)
	}

	return logs, nil
}
