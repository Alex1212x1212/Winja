package core

import (
	"os"
	"golang.org/x/sys/windows/registry"
)

const runKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`
const appName = "Winja"

func SetRunAtStartup(enable bool) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	if enable {
		exePath, err := os.Executable()
		if err != nil {
			return err
		}
		// Wrap in quotes in case of spaces in path
		return k.SetStringValue(appName, `"`+exePath+`"`)
	} else {
		err := k.DeleteValue(appName)
		if err != nil && err != registry.ErrNotExist {
			return err
		}
		return nil
	}
}

func IsRunAtStartup() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	
	val, _, err := k.GetStringValue(appName)
	return err == nil && val != ""
}

