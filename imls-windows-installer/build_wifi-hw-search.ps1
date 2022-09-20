#build_wifi-hw-search.ps1

#Requires -RunAsAdministrator

$installer_path = Get-Location

# -- Remove session-counter if it already exists
If (Test-Path -Path wifi-hardware-search-windows.exe) {
	Remove-Item wifi-hardware-search-windows.exe
}

# -- Move from estimating-wifi/imls-windowsinstaller to estimating-wifi
Set-Location ..

# -- Build the exe
Set-Location imls-raspberry-pi\cmd\wifi-hardware-search-windows
Write-Host "Building wifi-hardware-search executable."
# call the Go executable directly since we might have just installed Go and it
# may not be in our Path yet
& 'C:\Program Files\Go\bin\go.exe' build
$wd = Get-Location
$exe_path = "$wd\wifi-hardware-search-windows.exe"

# -- Move the exe into imls-windowsinstaller
Write-Host "Executable built. Terminating script"
Move-Item -Path $exe_path -Destination $installer_path

# -- Update location
Set-Location $installer_path
Exit
