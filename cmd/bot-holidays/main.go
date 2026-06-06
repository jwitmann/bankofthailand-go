package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	bot "github.com/jwitmann/bankofthailand-go"
)

func main() {
	var (
		year   = flag.Int("year", time.Now().Year(), "Year to fetch holidays for")
		format = flag.String("format", "json", "Output format: json, thaifa, csv")
	)
	flag.Parse()

	client, err := bot.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	holidays, err := client.GetHolidays(ctx, *year)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch *format {
	case "json":
		outputJSON(holidays)
	case "thaifa":
		outputThaiFA(holidays)
	case "csv":
		outputCSV(holidays)
	default:
		fmt.Fprintf(os.Stderr, "Unknown format: %s\n", *format)
		os.Exit(1)
	}
}

func outputJSON(holidays []bot.Holiday) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(holidays); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func outputThaiFA(holidays []bot.Holiday) {
	type holiday struct {
		HolidayWeekDay         string `json:"HolidayWeekDay"`
		HolidayWeekDayThai     string `json:"HolidayWeekDayThai"`
		Date                   string `json:"Date"`
		DateThai               string `json:"DateThai"`
		HolidayDescription     string `json:"HolidayDescription"`
		HolidayDescriptionThai string `json:"HolidayDescriptionThai"`
	}

	data := make([]holiday, len(holidays))
	for i, h := range holidays {
		data[i] = holiday{
			HolidayWeekDay:         h.HolidayWeekDay,
			HolidayWeekDayThai:     h.HolidayWeekDayThai,
			Date:                   h.Date,
			DateThai:               h.DateThai,
			HolidayDescription:     h.HolidayDescription,
			HolidayDescriptionThai: h.HolidayDescriptionThai,
		}
	}

	result := struct {
		API       string    `json:"api"`
		Timestamp string    `json:"timestamp"`
		Data      []holiday `json:"data"`
	}{
		API:       "API_V2.FIHolidays",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Data:      data,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func outputCSV(holidays []bot.Holiday) {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write([]string{"Date", "DateThai", "WeekDay", "WeekDayThai", "Description", "DescriptionThai"}); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
		os.Exit(1)
	}
	for _, h := range holidays {
		if err := w.Write([]string{
			h.Date,
			h.DateThai,
			h.HolidayWeekDay,
			h.HolidayWeekDayThai,
			h.HolidayDescription,
			h.HolidayDescriptionThai,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
			os.Exit(1)
		}
	}
	w.Flush()
}
