package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"gsa.gov/18f/read-initial-configuration/wordlist"
)

const VERSION = "0.0.1"
const lookup = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const yamlPath = "/etc/session-counter/auth.yaml"

const states = "(AL|AK|AZ|AR|CA|CO|CT|DE|FL|GA|HI|ID|IL|IN|IA|KS|KY|LA|ME|MD|MA|MI|MN|MS|MO|MT|NE|NV|NH|NJ|NM|NY|NC|ND|OH|OK|OR|PA|RI|SC|SD|TN|TX|UT|VT|VA|WA|WV|WI|WY|AS|DC|FM|GU|MH|MP|PW|PR|VI)"

type config struct {
	Token string `yaml:"token"`
	FCFS  string `yaml:"fcfs"`
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
	fmt.Print(text.FgYellow.Sprint("Look up and enter the FCFS Id Seq for this device's location: "))
	fmt.Println("https://www.imls.gov/search-compare/")
	fmt.Print(text.FgYellow.Sprint("It should look like: "))
	fmt.Println(text.FgWhite.Sprint("KY0069-003\n"))

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
			fmt.Println(text.FgRed.Sprint("That does not seem to be a full sequence ID."))
			fmt.Printf("Please try again.\n\n")
		}
	}

	fmt.Println(text.FgGreen.Sprint("\nGreat! Thank you!"))
	fmt.Println(text.FgWhite.Sprint("---------\n"))
	return fcfsid
}

func readToken() string {
	wordlist.Init()

	// We need to read things until they write DONE.
	reading := true
	reader := bufio.NewReader(os.Stdin)

	key := ""

	fmt.Println(text.FgYellow.Sprint("Enter a word pair and press enter."))
	fmt.Printf(text.FgYellow.Sprint("When you are done, type "))
	fmt.Printf(text.FgHiWhite.Sprint("DONE"))
	fmt.Println(text.FgYellow.Sprint(" and press return.\n"))
	fmt.Println(text.FgHiWhite.Sprint("------------------------------------"))

	for reading {
		fmt.Printf("word pair: ")
		pair, _ := reader.ReadString('\n')
		pair = strings.TrimSpace(pair)
		if pair == "DONE" || pair == "done" || pair == "quit" || pair == "exit" {
			// fmt.Println("key is", key)
			reading = false
		} else {
			ndx, err := wordlist.GetPairIndex(pair)
			if err != nil {
				fmt.Println(text.FgRed.Sprint("\n[ BAD! ] I can't find that word pair. Please try again, or DONE if you have no more word pairs.\n"))
			} else {
				decoded := decode(ndx)
				fmt.Println(text.FgGreen.Sprintf("\n[ GOOD! ] `%v` became `%v`\n", pair, decoded))
				key += decoded
			}
		}
	}
	return key
}

func writeYAML(cfg *config) {
	s, err := json.Marshal(cfg)
	if err != nil {
		errors.New("cannot marshal config for writing")
	}
	fmt.Println(string(s))
}

func main() {

	versionPtr := flag.Bool("version", false, "Get version and exit.")
	readFCFSPtr := flag.Bool("FCFS", false, "Read in their FCFS ID.")
	readTokenPtr := flag.Bool("token", false, "Read in their API token.")
	// writeFilePtr := flag.String("path", yamlPath, "Write a YAML file with requested fields.")
	flag.Parse()

	cfg := &config{}

	if *versionPtr {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *readFCFSPtr {
		cfg.FCFS = readFCFS()
	}

	if *readTokenPtr {
		cfg.Token = readToken()
	}

	writeYAML(cfg)
	fmt.Println(text.FgYellow.Sprint("All done!"))
}
