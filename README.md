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

Instruction :  Open Winja

<img width="1333" height="780" alt="winja_TUdDZy3rjt" src="https://github.com/user-attachments/assets/d27c0713-37cd-4461-9738-20039257fef3" />

The Open Settings :
<img width="1333" height="780" alt="winja_gDbDYdIKEL" src="https://github.com/user-attachments/assets/eef2e2f0-5506-4de2-8c8e-e42f6d2a9393" />

The API : Custom API KEY (On slow Down)
<img width="1325" height="759" alt="explorer_176tPRNTLW" src="https://github.com/user-attachments/assets/179cdad3-4ce9-4389-9ca5-6870d36a23d8" />

It is done

