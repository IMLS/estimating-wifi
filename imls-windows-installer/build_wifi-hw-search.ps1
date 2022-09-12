#build_wifi-hw-search.ps1

#Requires -RunAsAdministrator

$installer_path = Get-Location

# -- Move from estimating-wifi/imls-windowsinstaller to estimating-wifi
Set-Location ..

# -- Build the exe
Set-Location imls-raspberry-pi\cmd\wifi-hardware-search-windows
Write-Host "Building wifi-hardware-search executable."
go build
$wd = Get-Location
$exe_path = "$wd\wifi-hardware-search-windows.exe"

# -- Move the exe into imls-windowsinstaller
Write-Host "Executable built. Terminating script"
Move-Item -Path $exe_path -Destination $installer_path
Exit