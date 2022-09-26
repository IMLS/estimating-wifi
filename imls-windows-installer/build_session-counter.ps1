#build_session-counter.ps1

#Requires -RunAsAdministrator

$installer_path = Get-Location

# -- Remove session-counter if it already exists
If (Test-Path -Path session-counter.exe) {
	Remove-Item session-counter.exe
}

# -- Move to project directory
Set-Location ..

# -- Build the exe
Set-Location imls-wifi-sensor\cmd\windows-session-counter
Write-Host "Building session-counter executable."
# call the Go executable directly since we might have just installed Go and it
# may not be in our Path yet
& 'C:\Program Files\Go\bin\go.exe' build windows-session-counter.go
$wd = Get-Location
$exe_path = "$wd\session-counter.exe"

# -- Move the exe into imls-windows-installer
Move-Item -Path $exe_path -Destination $installer_path
Write-Host "Executable built and moved to imls-windows-installer. Terminating script"

# -- Update location
Set-Location $installer_path
Exit
