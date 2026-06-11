package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// CreateBackup writes a timestamped JSON dump of the whole database into the
// backup directory and prunes copies older than the retention window.
// Returns the created file path.
func (s *Service) CreateBackup() (string, error) {
	if err := os.MkdirAll(s.Cfg.BackupDir, 0o755); err != nil {
		return "", err
	}
	data := s.dumpData()
	payload := map[string]interface{}{
		"created_at": time.Now().Format(time.RFC3339),
		"tables":     data,
	}
	target := filepath.Join(s.Cfg.BackupDir, "rentacar_backup_"+time.Now().Format("20060102_150405")+".json")
	out, err := os.Create(target)
	if err != nil {
		return "", err
	}
	defer out.Close()
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		return "", err
	}
	s.pruneOldBackups()
	return target, nil
}

// pruneOldBackups removes backups older than the retention window.
func (s *Service) pruneOldBackups() {
	cutoff := time.Now().AddDate(0, 0, -s.Cfg.BackupRetentionDays)
	entries, err := os.ReadDir(s.Cfg.BackupDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), "rentacar_backup_") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(s.Cfg.BackupDir, e.Name()))
		}
	}
}

// ListBackups returns the available backups, newest first.
func (s *Service) ListBackups() []map[string]interface{} {
	entries, err := os.ReadDir(s.Cfg.BackupDir)
	if err != nil {
		return []map[string]interface{}{}
	}
	items := make([]map[string]interface{}, 0)
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), "rentacar_backup_") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		items = append(items, map[string]interface{}{
			"name":       e.Name(),
			"size":       info.Size(),
			"created_at": info.ModTime().Format(time.RFC3339),
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i]["name"].(string) > items[j]["name"].(string)
	})
	return items
}
