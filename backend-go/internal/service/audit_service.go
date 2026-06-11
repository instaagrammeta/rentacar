package service

import "github.com/instaagrammeta/rentacar/backend-go/internal/models"

// ListAuditLogs returns recent audit log entries (newest first).
func (s *Service) ListAuditLogs(limit, offset int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	var logs []models.AuditLog
	if err := s.DB.Order("id desc").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(logs))
	for i := range logs {
		out = append(out, logs[i].ToMap())
	}
	return out, nil
}
