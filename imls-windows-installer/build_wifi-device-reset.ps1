#build_wifi-device-reset.ps1

#Requires -RunAsAdministrator

$installer_path = Get-Location

# -- Remove wifi-hardware-search if it already exists
If (Test-Path -Path wifi-device-reset-windows.exe) {
	Remove-Item wifi-device-reset-windows.exe
}

# -- Move from estimating-wifi/imls-windowsinstaller to estimating-wifi
Set-Location ..

# -- Build the exe
Set-Location imls-wifi-sensor\cmd\wifi-device-reset-windows
Write-Host "Building wifi-device-reset-windows executable."
# call the Go executable directly since we might have just installed Go and it
# may not be in our Path yet
& 'C:\Program Files\Go\bin\go.exe' build
$wd = Get-Location
$exe_path = "$wd\wifi-device-reset-windows.exe"

# -- Move the exe into imls-windowsinstaller
Write-Host "Executable built. Terminating script"
Move-Item -Path $exe_path -Destination $installer_path

# -- Update location
Set-Location $installer_path
Exit
