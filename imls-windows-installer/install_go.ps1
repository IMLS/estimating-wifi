# install_go.ps1

#Requires -RunAsAdministrator

#TODO: Wanted to add a function to check for install status (for cleanliness), but the script seemed to be calling it incorrectly
<# function Check-Go-Install-Status () {
	$program = "Go Programming Language"
	$installation_status = (Get-ItemProperty HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\* | Where { $_.DisplayName -Match $program }) -ne $null
} #>

Write-Host "Running the install_go script. This downloads and installs Go 1.19."

# -- Detect if Golang is already installed on the machine
$program = "Go Programming Language"
$installation_status = (Get-ItemProperty HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\* | Where { $_.DisplayName -Match $program }) -ne $null

If(-Not $installation_status ) {	
	#TODO: Remove all testing purposes code eventually
	<# # -- Adding a new directory for testing purposes
	$current_time = Get-Date -Format o | ForEach-Object { $_ -replace ":", "-"}
	New-Item -Path C:\Users\Administrator\Downloads\TestGolang\$current_time -Type Directory #>

	# -- Use 1.19 Golang version
	$url = "https://go.dev/dl/go1.19.windows-amd64.msi"
	<# # -- Testing purposes
	$result = "C:\Users\Administrator\Downloads\TestGolang\$current_time\go1.19.windows-amd64.msi"
	 #>
	$result = "$env:userprofile\Downloads\go1.19.windows-amd64.msi"

	# -- Download
	try {
		$WebClient = New-Object System.Net.WebClient
		$WebClient.DownloadFile($url, $result)
		Write-Host "Downloaded successfully."
	} catch [Net.WebException] {
		Write-Host $_.Exception.ToString()
		Write-Host "Terminating install_go script because of exception."
		Exit
	}
	# -- Install on the system
<# 	# -- Testing purposes
	Start-Process C:\Users\Administrator\Downloads\TestGolang\$current_time\go1.19.windows-amd64.msi -Wait
 #>	
	Start-Process C:\Users\Administrator\Downloads\go1.19.windows-amd64.msi -Wait
	
	# -- Confirm install was successful
	$installation_status = (Get-ItemProperty HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\* | Where { $_.DisplayName -Match $program }) -ne $null
	
	If(-Not $installation_status ) {
		Write-Host "Go was not installed properly. Terminating install_go script."
		Exit
	} else {
		Write-Host "Go was installed. Terminating install_go script."
		Exit
	}
} else {
	Write-Host "Go is already installed on your system. Terminating install_go script."
	Exit
}