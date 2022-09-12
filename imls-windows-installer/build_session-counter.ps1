#build_session-counter.ps1

#Requires -RunAsAdministrator

$installer_path = Get-Location

# -- Move to project directory
Set-Location ..

# -- Build the exe
Set-Location imls-raspberry-pi\cmd\session-counter
Write-Host "Building session-counter executable."
go build session-counter.go
$wd = Get-Location
$exe_path = "$wd\session-counter.exe"

# -- Move the exe into imls-windows-installer
Move-Item -Path $exe_path -Destination $installer_path
Write-Host "Executable built and moved to imls-windows-installer. Terminating script"
Exit