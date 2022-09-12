#session-counter_build_prereqs.ps1

& "$PSScriptRoot\install_go.ps1"
& "$PSScriptRoot\install_innosetup.ps1"
& "$PSScriptRoot\build_session-counter.ps1"
& "$PSScriptRoot\build_wifi-hw-search.ps1"
