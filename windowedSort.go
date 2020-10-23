package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"time"
)

func main() {
	fileName := flag.String("fileName", "", "File to parse, without it stdin is used")
	dateFormat := flag.String("datefmt", "02/Jan/2006:15:04:05 +0000", "The date format, go style, to use to parse the date. Defaults to 02/Jan/2006:15:04:05 +0000")
	dateRegex := flag.String("dateRegex", "\\[(.+)\\]", "GO regex pattern to find the date in a line. Should return the date as the only group match.")
	windowSize := flag.Int("seconds", 10, "Window size in seconds. Defaults to 10.")

	flag.Parse()

	// read stdin forever
	buffer := make(map[time.Time][]string)
	var timeStamps []time.Time

	var scanner *bufio.Scanner
	if *fileName == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(*fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	regex := regexp.MustCompile(*dateRegex)

	for scanner.Scan() {
		line := scanner.Text()
		dateString := regex.FindStringSubmatch(line)[1]
		timeStamp, err := time.Parse(*dateFormat, dateString)
		if err != nil {
			log.Fatal(err)
		}
		_, found := buffer[timeStamp]
		if !found {
			buffer[timeStamp] = make([]string, 0)
			timeStamps = append(timeStamps, timeStamp)
			sort.Slice(timeStamps, func(i, j int) bool {
				return timeStamps[i].Before(timeStamps[j])
			})
		}
		buffer[timeStamp] = append(buffer[timeStamp], line)

		if timeStamps[0].Add(time.Duration(*windowSize) * time.Second).Before(timeStamp) {
			// pop the timestamp
			timeStamp, timeStamps = timeStamps[0], timeStamps[1:]
			// remove element from map
			lines := buffer[timeStamp]
			for _, oldLine := range lines {
				fmt.Println(oldLine)
			}
			delete(buffer, timeStamp)
		}
	}
}
