# Windows Installer

## Manual Build Process

This process assumes a user is using a Windows machine, has downloaded the entire repository as a .zip, extracted all files into a folder, and is currently in the imls-windows-installer directory.

### Download and Build the Prerequisite Software

- Open Windows powershell and cd to estimating-wifi\imls-windows-installer
- Enter `./session-counter_build_prereqs.ps1` to run the script
 1. Downloads, installs, and adds GoLang 1.19 to your Windows Path
  - If GoLang is already installed on your system, the script proceeds to the next step
 2. Downloads, installs, and adds the latest version of Inno Setup to your system
  - If Inno Setup is already installed on your system, the script proceeds to the nexts tep
 3. Builds session-counter.exe from the imls-raspberry-pi directory and moves it to the imls-windows-installer directory
  - If session-counter.exe is already found in imls-windows-installer, the script deletes it and builds a new one
 4. Builds wifi-hardware-search-windows.exe from the imls-raspberry-pi directory and moves it to the imls-windows-installer directory
  - If wifi-hardware-search-windows.exe is already found in imls-windows-installer, the script deletes it and builds a new one

### Build the Installer

- Open Inno Setup and choose imls-windows-installer\setup.iss as the configuration script
- Compile the executable by clicking 'Build' at the top toolbar

### Run the Installer
- Open the newly created estimating-wifi\imls-windows-installer\Output folder
- Right click on SessionCounterInstall.exe and click "Run as Administrator"
- Follow the instructions
- Note: the WiresharkPortable app cannot be placed in any Program Files directories

### Check to Ensure Session-Counter is Running
- Ctrl + Alt + Delete to open the Task Manager
- Scroll to the Background Processes section
- If `The non-sucking service manager` is running the `IMLS Session Counter` program, then it's running in the background as expected

## Automated Build Process

Once this is implemented, the README will be updated.