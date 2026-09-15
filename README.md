# Winja (Go Edition)

Winja is a lightning-fast, native Windows desktop application built with Go  toolkit. It serves as a comprehensive system analysis and malware scanning tool, heavily integrating with the VirusTotal API.

## Features

- **VirusTotal Integration**: Deep integration with the VirusTotal API (v3) to scan files, running processes, and report malicious activity.
- **Windows Explorer Integration**: Seamlessly adds "Scan selected items" and "Scan Folder" to the Windows right-click context menu (cascading menu).
- **Process Scanner**: Real-time process scanning that lists active processes, their elevations, and executable paths, allowing one-click VirusTotal analysis.
- **Drag & Drop Support**: Simply drop any file or folder onto the application window to instantly queue it for malware analysis.

## Running the Application

Simply execute the compiled binary:
`cmd
winja.exe
`
Or right-click any file in Windows Explorer and select **Winja -> Scan selected items**!

## Configuration
To scan files, you must provide your free VirusTotal API key in the **General Settings** tab of the application.
