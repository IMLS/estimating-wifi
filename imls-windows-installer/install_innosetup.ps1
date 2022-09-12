# install_innosetup.ps1

#Requires -RunAsAdministrator

#TODO: Wanted to add a function to check for install status (for cleanliness), but the script seemed to be calling it incorrectly
<# function Check-Inno-Install-Status () {
	$program = "Inno Setup"
	$installation_status = (Get-ItemProperty HKLM:\Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\* | Where { $_.DisplayName -Match $program }) -ne $null
} #>

Write-Host "Running the install_innosetup script. This downloads and installs the latest version of Inno Setup."

# -- Detect if Inno Setup is already installed on the machine
$program = "Inno Setup"
$installation_status = (Get-ItemProperty HKLM:\Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\* | Where { $_.DisplayName -Match $program }) -ne $null
# --NOTE: Wow6432 is needed because Inno Setup is a 32 bit program

If(-Not $installation_status ) {	
	#TODO: Remove all testing purposes code eventually
	<# # -- Adding a new directory for testing purposes
	$current_time = Get-Date -Format o | ForEach-Object { $_ -replace ":", "-"}
	New-Item -Path $env:userprofile\Downloads\TestInno\$current_time -Type Directory #>

	# -- Use the latest Inno Setup version
	$url = "https://jrsoftware.org/download.php/is.exe?site=1"
	<# # -- Testing purposes
	$result = "$env:userprofile\Downloads\TestInno\$current_time\innosetup.exe"
	 #>
	$result = "$env:userprofile\Downloads\innosetup.exe"

	# -- Download
	try {
		$WebClient = New-Object System.Net.WebClient
		$WebClient.DownloadFile($url, $result)
		Write-Host "Downloaded successfully."
	} catch [Net.WebException] {
		Write-Host $_.Exception.ToString()
		Write-Host "Terminating install_innosetup script because of exception."
		Exit
	}
	# -- Install on the system
<# 	# -- Testing purposes
	Start-Process $env:userprofile\Downloads\TestInno\$current_time\innosetup.exe -Wait
 #>	
	Start-Process $env:userprofile\Downloads\innosetup.exe -Wait -ArgumentList "/quiet /norestart /verysilent "
	
	# -- Confirm install was successful
	$installation_status = (Get-ItemProperty HKLM:\Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\* | Where { $_.DisplayName -Match $program }) -ne $null
	
	If(-Not $installation_status ) {
		Write-Host "Inno Setup was not installed properly. Terminating install_innosetup script."
		Exit
	} else {
		Write-Host "Inno Setup was installed. Terminating install_innosetup script."
		Exit
	}
} else {
	Write-Host "Inno Setup is already installed on your system. Terminating install_innosetup script."
	Exit
}