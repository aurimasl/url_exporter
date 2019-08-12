package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

// GetURL Retrieve content from URL
func GetURL(url string) ([]string, error) {
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)
	r := regexp.MustCompile(`(\w+) (\d+)$`)
	lines := strings.Split(responseString, "\n")

	var result []string
	for _, line := range lines {
		if line == "" {
			continue
		}
		result = append(result, line)
		if true != r.MatchString(line) {
			fmt.Println("Could not parse response:")
			fmt.Println(responseString)
			return nil, fmt.Errorf("Could not parse response")
		}
	}

	return result, nil
}

func retrieve(url string, textfilePath string) {
	for {
		timer := time.NewTimer(15 * time.Second)
		go func() {
			<-timer.C
			res, err := GetURL(url)
			if err != nil {
				fmt.Println("Could not retrieve metrics from ", url)
			}
			res = append(res, "")
			err = ioutil.WriteFile(textfilePath, []byte(strings.Join(res, "\n")), 0644)
			if err != nil {
				fmt.Println("Could not write to textfile ", textfilePath)
			}
		}()
		time.Sleep(15 * time.Second)
	}
}

func main() {
	var (
		url          = kingpin.Flag("url", "URL to retrieve metrics from").Required().String()
		textfilePath = kingpin.Flag("textfile", "Textfile path to save metrics to").Default("/opt/prometheus/textfile/url_exporter.prom").String()
	)
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("url_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting url_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	done := make(chan bool)
	go retrieve(*url, *textfilePath)
	<-done
}
