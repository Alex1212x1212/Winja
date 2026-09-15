package core

import (
	"os"
	"golang.org/x/sys/windows/registry"
)

func SetExplorerExtension(enable bool) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	bases := []string{
		`Software\Classes\*\shell\Winja`,
		`Software\Classes\Directory\shell\Winja`,
		`Software\Classes\AllFilesystemObjects\shell\Winja`,
	}

	for _, base := range bases {
		if enable {
			// cleanup old keys
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\Scan\command`)
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\Scan`)
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\FileScan\command`)
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\FileScan`)
			
			// 1. Create main cascading menu key
			k1, _, err := registry.CreateKey(registry.CURRENT_USER, base, registry.ALL_ACCESS)
			if err != nil {
				continue
			}
			k1.SetStringValue("MUIVerb", "Winja")
			k1.SetStringValue("Icon", `"`+exePath+`"`)
			// Use ExtendedSubCommandsKey for more reliable Windows 11 cascading
			k1.SetStringValue("ExtendedSubCommandsKey", base)
			k1.Close()

			// 2. Create sub-item 1: Scan selected items
			k2, _, err := registry.CreateKey(registry.CURRENT_USER, base+`\shell\1ScanItems`, registry.ALL_ACCESS)
			if err == nil {
				k2.SetStringValue("MUIVerb", "Scan selected items")
				k2.SetStringValue("Icon", `"`+exePath+`"`)
				k2.Close()
				
				k2c, _, _ := registry.CreateKey(registry.CURRENT_USER, base+`\shell\1ScanItems\command`, registry.ALL_ACCESS)
				k2c.SetStringValue("", `"`+exePath+`" "%1"`)
				k2c.Close()
			}

			// 3. Create sub-item 2: Scan Folder
			k3, _, err := registry.CreateKey(registry.CURRENT_USER, base+`\shell\2ScanFolder`, registry.ALL_ACCESS)
			if err == nil {
				k3.SetStringValue("MUIVerb", "Scan Folder")
				k3.SetStringValue("Icon", `"`+exePath+`"`)
				k3.Close()
				
				k3c, _, _ := registry.CreateKey(registry.CURRENT_USER, base+`\shell\2ScanFolder\command`, registry.ALL_ACCESS)
				k3c.SetStringValue("", `"`+exePath+`" "%1"`)
				k3c.Close()
			}
		} else {
			// Delete the entire tree
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\2ScanFolder\command`)
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\2ScanFolder`)
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\1ScanItems\command`)
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\1ScanItems`)
			
			// cleanup old ones just in case
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\Scan\command`)
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\Scan`)
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\FileScan\command`) 
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell\FileScan`) 
			
			registry.DeleteKey(registry.CURRENT_USER, base+`\shell`)
			registry.DeleteKey(registry.CURRENT_USER, base)
		}
	}
	return nil
}

func IsExplorerExtensionEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\AllFilesystemObjects\shell\Winja\shell\1ScanItems\command`, registry.QUERY_VALUE)
	if err == nil {
		k.Close()
		return true
	}
	
	k2, err2 := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\*\shell\Winja\shell\1ScanItems\command`, registry.QUERY_VALUE)
	if err2 == nil {
		k2.Close()
		return true
	}
	
	return false
}
