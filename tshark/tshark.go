package tshark
	
import (
    "fmt"
    "io/ioutil"
    "os/exec"
	"strings"
)

func Tshark (wlan string, duration int) map[string]int {

		tshark_path := "/usr/bin/tshark"
		
		tsharkCmd := exec.Command(tshark_path, "-a", fmt.Sprintf("duration:%d", duration), "-I", "-i", wlan, "-Tfields", "-e", "wlan.sa")

		tsharkOut, _ := tsharkCmd.StdoutPipe()
		tsharkCmd.Start()
		tsharkBytes, _ := ioutil.ReadAll(tsharkOut)
		tsharkCmd.Wait()
		macs := strings.Split(string(tsharkBytes), "\n")

		counts := make(map[string]int)

		for _, a_mac := range macs {
			if v, ok := counts[a_mac]; ok {
				counts[a_mac] = v + 1
			} else {
				counts[a_mac] = 1
			}
		}
		return counts
}