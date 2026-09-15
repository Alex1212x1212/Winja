package core

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Settings struct {
	WebcamActive         bool   `json:"webcamActive"`
	WebcamNotifStart     bool   `json:"webcamNotifStart"`
	WebcamNotifStop      bool   `json:"webcamNotifStop"`
	WebcamRequestReport  bool   `json:"webcamRequestReport"`
	WebcamAction         int    `json:"webcamAction"`
	WidgetShow     bool   `json:"widgetShow"`
	WidgetFixed    bool   `json:"widgetFixed"`
	WidgetTopMost  bool   `json:"widgetTopMost"`
	WidgetOpacity  bool   `json:"widgetOpacity"`
	NotifSticky    bool   `json:"notifSticky"`
	NotifCloseTime string `json:"notifCloseTime"`
	NotifPlacement string `json:"notifPlacement"`
	ProxyMode      int    `json:"proxyMode"`
	ProxyHost      string `json:"proxyHost"`
	ProxyPort      string `json:"proxyPort"`
	MinimizeToTray bool   `json:"minimizeToTray"`
	PrivateMode    bool   `json:"privateMode"`
	DoNotUpload    bool   `json:"doNotUpload"`
	OpenWeb        bool   `json:"openWeb"`
	NotifEnabled   bool   `json:"notifEnabled"`
	NotifLevel     int    `json:"notifLevel"`
	CritFail       bool   `json:"critFail"`
	APIKey         string `json:"apiKey"`
}

var (
	settingsFile string
	currentSettings Settings
	settingsMutex sync.Mutex
)

func init() {
	appData, err := os.UserConfigDir() // AppData\Roaming on Windows
	if err == nil {
		dir := filepath.Join(appData, "Winja")
		os.MkdirAll(dir, 0755)
		settingsFile = filepath.Join(dir, "settings.json")
	} else {
		settingsFile = "settings.json"
	}
	loadSettings()
}

func protectAPIKey(value string) (string, error) {
	if value == "" {
		return "", nil
	}

	keyBytes := []byte(value)
	in := windows.DataBlob{Data: &keyBytes[0], Size: uint32(len(keyBytes))}
	var out windows.DataBlob
	if err := windows.CryptProtectData(&in, nil, nil, 0, nil, 0, &out); err != nil {
		return "", err
	}
	defer func() {
		_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	}()

	return base64.StdEncoding.EncodeToString(unsafe.Slice(out.Data, int(out.Size))), nil
}

func unprotectAPIKey(value string) (string, error) {
	if value == "" {
		return "", nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	if len(ciphertext) == 0 {
		return "", nil
	}

	in := windows.DataBlob{Data: (*byte)(unsafe.Pointer(&ciphertext[0])), Size: uint32(len(ciphertext))}
	var out windows.DataBlob
	if err := windows.CryptUnprotectData(&in, nil, nil, 0, nil, 0, &out); err != nil {
		return "", err
	}
	defer func() {
		_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	}()

	return string(unsafe.Slice(out.Data, int(out.Size))), nil
}

func loadSettings() {
	// Not acquiring lock here because init runs sequentially and we lock in getters/setters
	
	currentSettings = Settings{
		WebcamNotifStart:    true,
		WebcamAction:        0,
		WidgetShow:     true,
		NotifCloseTime: "10",
		NotifPlacement: "Bottom Centered",
		ProxyMode:      1,
		MinimizeToTray: true,
		NotifEnabled:   true,
		CritFail:       true,
	}

	data, err := os.ReadFile(settingsFile)
	if err == nil {
		json.Unmarshal(data, &currentSettings)
		if currentSettings.APIKey != "" {
			if key, err := unprotectAPIKey(currentSettings.APIKey); err == nil {
				currentSettings.APIKey = key
			}
		}
	}
}

func saveSettings() error {
	settings := currentSettings
	if settings.APIKey != "" {
		encryptedKey, err := protectAPIKey(settings.APIKey)
		if err != nil {
			return err
		}
		settings.APIKey = encryptedKey
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(settingsFile, data, 0644)
}

func SetMinimizeToTray(enable bool) error {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.MinimizeToTray = enable
	return saveSettings()
}

func IsMinimizeToTray() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.MinimizeToTray
}

func SetPrivateMode(enable bool) error {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.PrivateMode = enable
	return saveSettings()
}

func IsPrivateMode() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.PrivateMode
}

func SetDoNotUpload(enable bool) error {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.DoNotUpload = enable
	return saveSettings()
}

func IsDoNotUpload() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.DoNotUpload
}

func SetOpenWeb(enable bool) error {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.OpenWeb = enable
	return saveSettings()
}

func IsOpenWeb() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.OpenWeb
}

func SetNotifEnabled(enable bool) error {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.NotifEnabled = enable
	return saveSettings()
}

func IsNotifEnabled() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.NotifEnabled
}

func SetNotifLevel(level int) error {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.NotifLevel = level
	return saveSettings()
}

func GetNotifLevel() int {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.NotifLevel
}

func SetCritFail(enable bool) error {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.CritFail = enable
	return saveSettings()
}

func IsCritFail() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.CritFail
}

func SetAPIKey(key string) error {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.APIKey = key
	return saveSettings()
}

func GetAPIKey() string {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.APIKey
}

func SetWidgetShow(enable bool) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.WidgetShow = enable
	saveSettings()
}
func IsWidgetShow() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.WidgetShow
}

func SetWidgetFixed(enable bool) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.WidgetFixed = enable
	saveSettings()
}
func IsWidgetFixed() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.WidgetFixed
}

func SetWidgetTopMost(enable bool) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.WidgetTopMost = enable
	saveSettings()
}
func IsWidgetTopMost() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.WidgetTopMost
}

func SetWidgetOpacity(enable bool) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.WidgetOpacity = enable
	saveSettings()
}
func IsWidgetOpacity() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.WidgetOpacity
}

func SetNotifSticky(enable bool) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.NotifSticky = enable
	saveSettings()
}
func IsNotifSticky() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.NotifSticky
}

func SetNotifCloseTime(time string) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.NotifCloseTime = time
	saveSettings()
}
func GetNotifCloseTime() string {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.NotifCloseTime
}

func SetNotifPlacement(place string) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.NotifPlacement = place
	saveSettings()
}
func GetNotifPlacement() string {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.NotifPlacement
}

func SetProxyMode(mode int) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.ProxyMode = mode
	saveSettings()
}
func GetProxyMode() int {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.ProxyMode
}

func SetProxyHost(host string) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.ProxyHost = host
	saveSettings()
}
func GetProxyHost() string {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.ProxyHost
}

func SetProxyPort(port string) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.ProxyPort = port
	saveSettings()
}
func GetProxyPort() string {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.ProxyPort
}

func SetWebcamActive(enable bool) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.WebcamActive = enable
	saveSettings()
}
func IsWebcamActive() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.WebcamActive
}

func SetWebcamNotifStart(enable bool) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.WebcamNotifStart = enable
	saveSettings()
}
func IsWebcamNotifStart() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.WebcamNotifStart
}

func SetWebcamNotifStop(enable bool) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.WebcamNotifStop = enable
	saveSettings()
}
func IsWebcamNotifStop() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.WebcamNotifStop
}

func SetWebcamRequestReport(enable bool) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.WebcamRequestReport = enable
	saveSettings()
}
func IsWebcamRequestReport() bool {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.WebcamRequestReport
}

func SetWebcamAction(action int) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	currentSettings.WebcamAction = action
	saveSettings()
}
func GetWebcamAction() int {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	return currentSettings.WebcamAction
}
