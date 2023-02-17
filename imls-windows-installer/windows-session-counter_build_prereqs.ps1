#windows-session-counter_build_prereqs.ps1

& "$PSScriptRoot\install_go.ps1"
& "$PSScriptRoot\install_innosetup.ps1"
& "$PSScriptRoot\build_windows-session-counter.ps1"
& "$PSScriptRoot\build_wifi-hw-search.ps1"
& "$PSScriptRoot\build_wifi-device-reset.ps1"
