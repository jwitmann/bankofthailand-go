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

	switch *format {
	case "json":
		holidays, err := client.GetHolidays(ctx, *year)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		outputJSON(holidays)
	case "thaifa":
		resp, err := client.GetHolidaysRaw(ctx, *year)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		outputThaiFA(resp)
	case "csv":
		holidays, err := client.GetHolidays(ctx, *year)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
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

func outputThaiFA(resp *bot.HolidaysResponse) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resp); err != nil {
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
