# Windows Installer

## Automated Build Process

This process assumes a user is using a Windows machine and has admin rights.

### Download the installer

Visit [this page](https://github.com/IMLS/estimating-wifi/actions/runs/3100996980) and click on the SessionCounterInstall artifact.

Navigate to your Downloads folder, and extract the .zip file.

### Run the Installer
- Right click on SessionCounterInstall.exe and click "Run as Administrator"
- Follow the instructions
- Note: the WiresharkPortable app cannot be placed in any Program Files directories

### Check to Ensure Session-Counter is Running
- Ctrl + Alt + Delete to open the Task Manager
- Scroll to the Background Processes section
- If `Windows Service Wrapper` is running the `estimating-wifi` program, then it's running in the background as expected

## Manual Build Process

This process assumes a user is using a Windows machine, has admin rights, has downloaded the entire repository as a .zip, extracted all files into a folder, and is currently in the imls-windows-installer directory.

### Download and Build the Prerequisite Software

- Open Windows powershell and cd to estimating-wifi\imls-windows-installer
- Enter `./session-counter_build_prereqs.ps1` to run the script
 1. Downloads, installs, and adds GoLang 1.19 to your Windows Path
- If GoLang is already installed on your system, the script proceeds to the next step
 2. Downloads, installs, and adds the latest version of Inno Setup to your system
- If Inno Setup is already installed on your system, the script proceeds to the nexts tep
 3. Builds session-counter.exe from the imls-wifi-sensor directory and moves it to the imls-windows-installer directory
- If session-counter.exe is already found in imls-windows-installer, the script deletes it and builds a new one
 4. Builds wifi-hardware-search-windows.exe from the imls-wifi-sensor directory and moves it to the imls-windows-installer directory
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
- If `Windows Service Wrapper` is running the `estimating-wifi` program, then it's running in the background as expected

# Full stack testing

To run end-to-end tests on Windows requires WSL2. The basic idea is to run the backend in a Linux VM and have the Windows service "talk" to this backend. Steps:

- Install WSL2 and your favorite Linux distribution.
  - Ubuntu has a good guide on [how to install Ubuntu on WSL2](https://ubuntu.com/tutorials/install-ubuntu-on-wsl2-on-windows-10#1-overview).
- In Ubuntu, install the tools necessary to run Docker. `sudo apt-get -y install docker-compose docker.io`
  - You might need to run `sudo apt-get update` beforehand.
- Running Docker in Ubuntu in WSL can be clunky. We'll have to start the service manually because systemd is not available (although this has [changed recently](https://devblogs.microsoft.com/commandline/systemd-support-is-now-available-in-wsl/)):
  - `sudo dockerd`
- Now we can build the backend.
  - Open up a separate Ubuntu terminal
  - Clone this repository
  - Change directory to the `imls-backend` directory
  - Run `DOCKER_HOST=unix:///var/run/docker.sock docker-compose up`
    - Please note, if this command fails as a normal user, you can add your user to the group `sudo usermod -aG docker $USER` or you can just `sudo` this command.
- You probably also want to run migrations. Unfortunately, `dbmate` is not easily installable on Ubuntu, so we'll download the binary ourselves.
  - Open up a separate (third!) Ubuntu terminal
  - `sudo curl -fsSL -o dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64`
  - `sudo chown +x dbmate`
  - `./dbmate up`
- To verify that all is running correctly:
  - Open up a separate Powershell terminal and run `curl.exe -v 127.0.0.1:300/presences`
- Install the Windows session-counter service as normal, if you haven't already.
- Edit your session-counter.ini:
  - Under `[api]`, you should have: `host=127.0.0.1:3000`
  - The windows session-counter should run with the default ini configuration
  - _But_ if you need to change any ini settings:
    - After changing the ini file, restart session-counter by running `WinSw-x64.exe restart` in the installed IMLS Session Counter directory.
