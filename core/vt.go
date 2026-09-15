package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type VTReport struct {
	Data struct {
		Attributes struct {
			LastAnalysisStats struct {
				Harmless   int `json:"harmless"`
				Malicious  int `json:"malicious"`
				Suspicious int `json:"suspicious"`
				Undetected int `json:"undetected"`
			} `json:"last_analysis_stats"`
		} `json:"attributes"`
	} `json:"data"`
	Error *struct {
		Code string `json:"code"`
	} `json:"error"`
}

func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil { return "", err }
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil { return "", err }
	return hex.EncodeToString(h.Sum(nil)), nil
}

func GetVTReport(hash string, apiKey string) (*VTReport, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("no API key provided")
	}

	req, _ := http.NewRequest("GET", "https://www.virustotal.com/api/v3/files/"+hash, nil)
	req.Header.Add("x-apikey", apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("invalid API key")
	}

	var report VTReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, err
	}
	
	if report.Error != nil && report.Error.Code == "NotFoundError" {
		return nil, fmt.Errorf("not_found")
	}
	return &report, nil
}
