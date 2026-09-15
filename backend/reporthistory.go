package backend

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
	"sync"
)

type ReportEntry struct {
	ReportDate      string `json:"reportDate"`
	FileName        string `json:"fileName"`
	Status          string `json:"status"`
	DetectionRate   string `json:"detectionRate"`
	LastScannedDate string `json:"lastScannedDate"`
	SHA1            string `json:"sha1"`
}

var logMutex sync.Mutex

func LogReportHistory(filePath string, res ScanResult, sha1 string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	dir := "C:\\Winja"
	os.MkdirAll(dir, 0755)
	logFile := filepath.Join(dir, "ReportHistory.txt")

	var history []ReportEntry
	data, err := os.ReadFile(logFile)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, &history)
	}

	entry := ReportEntry{
		ReportDate:      time.Now().Format("2006-01-02 15:04:05"),
		FileName:        filepath.Base(filePath),
		Status:          res.Status,
		DetectionRate:   res.DetectionRate,
		LastScannedDate: time.Now().Format("2006-01-02 15:04:05"),
		SHA1:            sha1,
	}

	history = append([]ReportEntry{entry}, history...)

	out, _ := json.MarshalIndent(history, "", "  ")
	os.WriteFile(logFile, out, 0644)
}

func GetReportHistory() []ReportEntry {
	logMutex.Lock()
	defer logMutex.Unlock()
	
	logFile := "C:\\Winja\\ReportHistory.txt"
	var history []ReportEntry
	data, err := os.ReadFile(logFile)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, &history)
	}
	return history
}
