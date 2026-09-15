package backend

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ScanResult struct {
	Status        string
	DetectionRate string
	Color         string
}

func GetFileHashes(filePath string) (sha256Hash, sha1Hash string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	h256 := sha256.New()
	h1 := sha1.New()

	buf := make([]byte, 4096)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			h256.Write(buf[:n])
			h1.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", err
		}
	}

	return hex.EncodeToString(h256.Sum(nil)), hex.EncodeToString(h1.Sum(nil)), nil
}

func parseStats(stats map[string]interface{}) ScanResult {
	malicious := int(stats["malicious"].(float64))
	undetected := int(stats["undetected"].(float64))
	harmless := int(stats["harmless"].(float64))

	total := malicious + undetected + harmless
	rate := fmt.Sprintf("%d / %d", malicious, total)

	if malicious > 0 {
		return ScanResult{Status: "Infected", Color: "#e74c3c", DetectionRate: rate}
	} else {
		return ScanResult{Status: "Clean", Color: "#2ecc71", DetectionRate: rate}
	}
}

func PerformScan(filePath, apiKey string, doNotUpload bool, cb func(ScanResult)) {
	defer func() {
		if r := recover(); r != nil {
			cb(ScanResult{Status: fmt.Sprintf("Crash: %v", r), Color: "#e74c3c", DetectionRate: "? / ?"})
		}
	}()
	if apiKey == "" {
		cb(ScanResult{Status: "No API Key (See Settings)", Color: "#e74c3c", DetectionRate: "? / ?"})
		return
	}

	sha256Hash, _, err := GetFileHashes(filePath)
	if err != nil {
		cb(ScanResult{Status: "File Read Error", Color: "#e74c3c", DetectionRate: "? / ?"})
		return
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://www.virustotal.com/api/v3/files/%s", sha256Hash), nil)
	req.Header.Add("x-apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		cb(ScanResult{Status: "Network Error", Color: "#e74c3c", DetectionRate: "? / ?"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		data := result["data"].(map[string]interface{})
		attr := data["attributes"].(map[string]interface{})
		stats := attr["last_analysis_stats"].(map[string]interface{})

		cb(parseStats(stats))
	} else if resp.StatusCode == 404 {
		if doNotUpload {
			cb(ScanResult{Status: "Not Found (Upload Disabled)", Color: "#f39c12", DetectionRate: "? / ?"})
			return
		}
		cb(ScanResult{Status: "Uploading...", Color: "#f39c12", DetectionRate: "? / ?"})

		// Upload the file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(filePath))
		if err == nil {
			f, _ := os.Open(filePath)
			io.Copy(part, f)
			f.Close()
		}
		writer.Close()

		upReq, _ := http.NewRequest("POST", "https://www.virustotal.com/api/v3/files", body)
		upReq.Header.Add("x-apikey", apiKey)
		upReq.Header.Add("Content-Type", writer.FormDataContentType())

		upResp, err := client.Do(upReq)
		if err != nil || upResp.StatusCode != 200 {
			cb(ScanResult{Status: "Upload Failed", Color: "#e74c3c", DetectionRate: "? / ?"})
			return
		}
		
		var upResult map[string]interface{}
		json.NewDecoder(upResp.Body).Decode(&upResult)
		upResp.Body.Close()

		analysisData := upResult["data"].(map[string]interface{})
		analysisID := analysisData["id"].(string)

		cb(ScanResult{Status: "Queued / Analyzing...", Color: "#f39c12", DetectionRate: "? / ?"})
		
		// Poll for analysis completion
		for {
			time.Sleep(5 * time.Second)
			pollReq, _ := http.NewRequest("GET", fmt.Sprintf("https://www.virustotal.com/api/v3/analyses/%s", analysisID), nil)
			pollReq.Header.Add("x-apikey", apiKey)
			
			pollResp, err := client.Do(pollReq)
			if err != nil || pollResp.StatusCode != 200 {
				if pollResp != nil { pollResp.Body.Close() }
				continue
			}
			
			var pollResult map[string]interface{}
			json.NewDecoder(pollResp.Body).Decode(&pollResult)
			pollResp.Body.Close()
			
			pollData := pollResult["data"].(map[string]interface{})
			pollAttr := pollData["attributes"].(map[string]interface{})
			status := pollAttr["status"].(string)
			
			if status == "completed" {
				stats := pollAttr["stats"].(map[string]interface{})
				cb(parseStats(stats))
				return
			}
		}
	} else {
		cb(ScanResult{Status: "API Error", Color: "#e74c3c", DetectionRate: "? / ?"})
	}
}


