package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

func main() {
	trains, err := Parse()
	if err != nil {
		panic(err)
	}
	csvWriter := csv.NewWriter(os.Stdout)
	err = csvWriter.Write([]string{"category", "number", "name", "from_station",
		"from_time", "to_station", "to_time", "current_station", "scheduled_arrival",
		"actual_arrival", "delay"})
	if err != nil {
		panic(err)
	}
	err = csvWriter.WriteAll(makeRecords(trains))
	if err != nil {
		panic(err)
	}
}

func makeRecords(trains []TrainInfo) [][]string {
	records := make([][]string, len(trains), len(trains))
	for i, train := range trains {
		records[i] = makeRecord(train)
	}
	return records
}

func makeRecord(t TrainInfo) []string {
	fields := make([]string, 11, 11)
	fields[0] = t.Category
	fields[1] = strconv.Itoa(t.Number)
	fields[2] = t.Name
	fields[3] = t.From.Station
	fields[4] = formatTime(t.From.Time)
	fields[5] = t.To.Station
	fields[6] = formatTime(t.To.Time)
	fields[7] = t.Current.Station
	fields[8] = formatTime(t.Current.Time)
	fields[9] = formatTime(t.Current.Actually)
	fields[10] = strconv.Itoa(t.Current.Delay)
	return fields
}

func formatTime(t time.Time) string {
	if t.Year() > 1 {
		return t.Format("02.01.2006 15:04")
	} else {
		return ""
	}
}
