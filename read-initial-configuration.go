package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/fatih/color"
	"gsa.gov/18f/read-initial-configuration/wordlist"
)

const VERSION = "0.0.1"
const lookup = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const yamlPath = "/etc/session-counter/auth.yaml"

const states = "(AL|AK|AZ|AR|CA|CO|CT|DE|FL|GA|HI|ID|IL|IN|IA|KS|KY|LA|ME|MD|MA|MI|MN|MS|MO|MT|NE|NV|NH|NJ|NM|NY|NC|ND|OH|OK|OR|PA|RI|SC|SD|TN|TX|UT|VT|VA|WA|WV|WI|WY|AS|DC|FM|GU|MH|MP|PW|PR|VI)"

type config struct {
	Token string `yaml:"api_token"`
	FCFS  string `yaml:"fcfs_seq_id"`
	Tag   string `yaml:"tag"`
}

func decode(ndx int) string {
	mask := int(math.Exp2(6) - 1)
	result := ""
	for loopndx := 2; loopndx > -1; loopndx-- {
		c := (ndx & (mask << (6 * loopndx))) >> (6 * loopndx)
		if c < 63 {
			result += string(lookup[c])
		} else {
			result += " "
		}
	}
	return strings.TrimSpace(result)
}

func readFCFS() string {
	msg := ""
	msg += "Look up and enter the FCFS Id Seq for this device's location\n\n"
	msg += color.BlueString("https://www.imls.gov/search-compare/\n\n")
	msg += "The FCFS Seq Id should look like: "
	msg += color.New(color.FgYellow).Sprint("KY0069-003")
	fmt.Println(box(color.New(color.FgBlue), msg))

	reader := bufio.NewReader(os.Stdin)
	re := regexp.MustCompile(states + `[0-9]{4}-[0-9]{3}`)
	fcfsid := ""

	matched := false
	for !matched {
		fmt.Print("FCFS Seq Id: ")
		fcfsid, _ = reader.ReadString('\n')
		fcfsid = strings.TrimSpace(fcfsid)

		if re.MatchString(fcfsid) {
			matched = true
		} else {
			color.New(color.FgRed).Println("\nThat does not seem to be a full sequence ID.")
			color.New(color.FgYellow).Println("It should be a state abbreviation, four digits, a hyphen, and three more digits.")
			fmt.Printf("Please try again.\n\n")
		}
	}

	fmt.Println()
	yay := box(color.New(color.FgHiGreen), color.New(color.FgGreen).Sprint("Great! Thank you!"))
	fmt.Print(yay)

	// We will pause, because people might get confused if it goes too fast.
	time.Sleep(2 * time.Second)

	return fcfsid
}

func box(c *color.Color, s string) string {
	// Expand tabs for measurement.
	s = strings.Replace(s, "\t", "  ", -1)
	msg := strings.Split(s, "\n")
	max := 0
	for _, s := range msg {
		// String color codes before measuring
		s = stripansi.Strip(s)
		if len(s) > max {
			// fmt.Println(s, len(s))
			max = len(s)
		}
	}
	result := ""
	result += c.Sprint("╔")
	for ndx := 0; ndx < max; ndx++ {
		result += c.Sprint("═")
	}
	result += c.Sprint("══╗\n")
	for _, line := range msg {
		result += c.Sprint("║ ")
		result += line
		// Strip color codes before measuring.
		line = stripansi.Strip(line)
		if len(line) < max {
			for i := 0; i < max-len(line); i++ {
				result += " "
			}
		}
		result += c.Sprint(" ║\n")
	}
	result += c.Sprint("╚")
	for ndx := 0; ndx < max; ndx++ {
		result += c.Sprint("═")
	}
	result += c.Sprint("══╝\n")

	return result
}

func readToken() string {
	wordlist.Init()

	// We need to read things until they write DONE.
	reading := true
	reader := bufio.NewReader(os.Stdin)

	key := ""
	textc := color.New(color.FgCyan)
	boxc := color.New(color.FgBlue)

	msg := "Enter your token word-pairs:\n\n"
	msg += "\t1) one pair at a time, and\n"
	msg += "\t2) in order.\n\n"
	msg += "When you are done, type "
	msg += textc.Sprint("DONE")
	msg += " and press return.\n\n"
	msg += color.New(color.FgYellow).Sprint("There should be 14 word pairs.")
	fmt.Println(box(boxc, msg))
	wpndx := 1
	for reading {
		fmt.Printf("Word pair %d: ", wpndx)
		pair, _ := reader.ReadString('\n')
		pair = strings.TrimSpace(pair)
		if pair == "DONE" || pair == "done" || pair == "quit" || pair == "exit" {
			// fmt.Println("key is", key)
			reading = false
		} else {
			ndx, err := wordlist.GetPairIndex(pair)
			if err != nil {
				color.New(color.FgBlue).Println("\n[ BAD! ] I can't find that word pair. Please try again, or DONE if you have no more word pairs.\n")
			} else {
				wpndx += 1
				decoded := decode(ndx)
				color.New(color.FgGreen).Printf("\n[ GOOD! ] `%v` became `%v`\n\n", pair, decoded)
				key += decoded
			}
		}
	}
	return key
}

func readTag() string {
	msg := ""
	msg += "Enter your device tag.\n\n"
	msg += "This will end up in the data, and will help you identify the device.\n\n"
	msg += "Examples:\n\n"
	msg += "\t1) " + color.New(color.FgYellow).Sprint("behind refdesk") + "\n"
	msg += "\t2) " + color.New(color.FgYellow).Sprint("in collections") + "\n"
	msg += "\t3) " + color.New(color.FgYellow).Sprint("on first floor") + "\n\n"
	msg += "The purpose is to allow you to uniquely identify this Pi.\n\n"
	msg += color.New(color.FgYellow).Sprint("We will truncate this at 255 characters.")
	fmt.Print(box(color.New(color.FgBlue), msg))

	fmt.Print("Device tag: ")
	reader := bufio.NewReader(os.Stdin)
	tag, _ := reader.ReadString('\n')
	tag = strings.TrimSpace(tag)

	fmt.Println()
	yay := box(color.New(color.FgHiGreen), color.New(color.FgGreen).Sprint("Awesome!"))
	fmt.Print(yay)
	time.Sleep(2 * time.Second)

	return tag
}

func writeYAML(cfg *config, path string) {
	s := fmt.Sprintf(`api_token: "%v"`, cfg.Token)
	s += "\n"
	s += fmt.Sprintf(`fcfs_seq_id: "%v"`, cfg.FCFS)
	s += "\n"
	s += fmt.Sprintf(`tag: "%v"`, cfg.Tag)
	s += "\n"
	// This will truncate the file if it exists.
	f, err := os.Create(path)
	if err != nil {
		log.Fatal("could not open config for writing")
	}
	f.WriteString(s)
	f.Close()
}

func main() {

	configPathPtr := flag.String("path", yamlPath, "Where to write the config.")
	versionPtr := flag.Bool("version", false, "Get version and exit.")
	readFCFSPtr := flag.Bool("fcfs-seq", false, "Read in their FCFS ID.")
	readTokenPtr := flag.Bool("token", false, "Read in their API token.")
	tagPtr := flag.Bool("tag", false, "A local inventory tag or identifier.")
	allPtr := flag.Bool("all", false, "Enables all values for entry.")
	// writeFilePtr := flag.String("path", yamlPath, "Write a YAML file with requested fields.")
	flag.Parse()

	cfg := &config{}

	if *versionPtr {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *allPtr {
		*readFCFSPtr = true
		*tagPtr = true
		*readTokenPtr = true
	}

	if *readFCFSPtr {
		fmt.Println()
		cfg.FCFS = readFCFS()
	}

	if *tagPtr {
		fmt.Println()
		cfg.Tag = readTag()
	}

	if *readTokenPtr {
		fmt.Println()
		cfg.Token = readToken()
	}

	// Writes to the default location, or another location
	// if overwridden by the flag.
	writeYAML(cfg, *configPathPtr)

	fmt.Println()
	fmt.Println(box(color.New(color.FgHiBlue), color.New(color.FgYellow).Sprint("All done!")))
}
