package main

import (
	"flag"
	"fmt"
	"time"

	monitor "github.com/acidleroy/http-log-monitor"
)

// Example 1: A single string flag called "species" with default value "gopher".
var site = flag.String("site", "my.site.com", "The URL you are interested in acquiring statistics on.")
var logFile = flag.String("logfile", "test.log", "The log file you want to monitor")
var rollingAverage = flag.Int64("roll", 2, "Amount of time in minutes to keep a rolling average (must be an integer)")
var alarmRate = flag.Float64("alarmthreshold", 1.0, "Threshold in accesses per minute to use for trigering a high traffic alarm.")

func main() {

	flag.Parse()
	avg := monitor.NewOverallTimeAverage()
	rollingAverage := monitor.NewRollingTimeAverage(*rollingAverage)
	stats := monitor.NewLogStats(*site, avg, rollingAverage, float32(*alarmRate))
	reader := monitor.NewLogReader(*logFile)

	startTs := time.Now()
	for true {
		entries, err := reader.GetNewLogEntries()
		if err != nil {
			fmt.Println("There was an error getting the log entries: ", err)
			return
		}
		for _, v := range entries {
			err2 := stats.ProcessEntry(&v)
			if err2 != nil {
				fmt.Printf("There was an error processing %s: %s\n", v, err2)
			}
		}
		time.Sleep(100 * time.Millisecond)
		if time.Now().Sub(startTs) >= time.Second*10 {
			fmt.Println("Printing top 10 most popular sections")
			stats.PrintPopulartSections(10) // Print the 10 most popular sections
			startTs = time.Now()
		}
	}
}
