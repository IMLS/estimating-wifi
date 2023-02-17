package netadapter

// netadapter: pull in basic network adapter properties and property keys from
// the Plug and Play Device Property API on Windows.

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
)

var (
	// backticks are used as escape sequences in PowerShell, and go multiline
	// strings (`...`) do not support escaping backticks.
	findNetPSCommand = "Get-NetAdapter -Physical | `\n" +
		"Select-Object -Property Name,MacAddress,DeviceID,InterfaceDescription,PnpDeviceID | `\n" +
		"ForEach-Object { `\n" +
		"    $Manufacturer = Get-PnpDeviceProperty -InstanceId $_.PnpDeviceID -KeyName DEVPKEY_Device_Manufacturer | Select -ExpandProperty Data\n" +
		"    Add-Member -InputObject $_ -NotePropertyName Manufacturer -NotePropertyValue $Manufacturer -PassThru `\n" +
		"} | `\n" +
		"ConvertTo-Json"
)

func RestartNetAdapter(AdapterName string) {
	ps := New()
	var restartNetPSCommand = "Get-NetAdapter -Physical -Name \"" + AdapterName + "\"| Restart-NetAdapter -Confirm:$false"
	ps.Execute(restartNetPSCommand)
}

type PowerShell struct {
	powerShell string
}

type NetInfo struct {
	Name                 string
	MacAddress           string
	DeviceID             string
	InterfaceDescription string
	PnpDeviceID          string
	Manufacturer         string
}

func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) Execute(args ...string) []byte {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).Msg("Powershell: cannot start command")
	}
	return out.Bytes()
}

func GetDeviceHash(wlan *models.Device) []map[string]string {
	ps := New()
	lines := ps.Execute(findNetPSCommand)
	var netinfo []NetInfo
	json.Unmarshal([]byte(lines), &netinfo)
	result := make([]map[string]string, len(netinfo))
	for _, net := range netinfo {
		hash := make(map[string]string)
		hash["logical name"] = net.Name
		hash["serial"] = net.MacAddress
		hash["physical id"] = net.DeviceID
		hash["description"] = net.InterfaceDescription
		hash["vendor"] = net.Manufacturer
		// below fields are not applicable: leave blank
		hash["bus info"] = ""
		hash["configuration"] = ""
		result = append(result, hash)
	}
	return result
}
