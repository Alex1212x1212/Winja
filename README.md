# Winja (Go Edition)

Winja is a lightning-fast, native Windows desktop application built with Go and the [Fyne](https://fyne.io/) toolkit. It serves as a comprehensive system analysis and malware scanning tool, heavily integrating with the VirusTotal API.

## Features

- **VirusTotal Integration**: Deep integration with the VirusTotal API (v3) to scan files, running processes, and report malicious activity.
- **Windows Explorer Integration**: Seamlessly adds "Scan selected items" and "Scan Folder" to the Windows right-click context menu (cascading menu).
- **Process Scanner**: Real-time process scanning that lists active processes, their elevations, and executable paths, allowing one-click VirusTotal analysis.
- **Drag & Drop Support**: Simply drop any file or folder onto the application window to instantly queue it for malware analysis.
- **System Tray Mode**: Can be minimized to the system tray to run quietly in the background, receiving files via TCP sockets when launched from Explorer.
- **Fully Native & Thread-Safe**: Rewritten completely in Go, replacing the old Python/CustomTkinter implementation with a highly concurrent, memory-safe, and native architecture.

## Installation & Build

Make sure you have Go installed on your system.
This project uses CGO for the Fyne GUI, so a GCC compiler (like TDM-GCC or MSYS2) is required on Windows.

1. Clone or download the repository.
2. Initialize and download dependencies:
   `ash
   go mod tidy
   `
3. Build the application for Windows (hides the console window):
   `ash
   go build -ldflags="-s -w -H windowsgui" -o winja.exe .
   `

## Running the Application

Simply execute the compiled binary:
`cmd
winja.exe
`
Or right-click any file in Windows Explorer and select **Winja -> Scan selected items**!

## Configuration
To scan files, you must provide your free VirusTotal API key in the **General Settings** tab of the application.
