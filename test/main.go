package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	capital := flag.Float64("capital", 0.0, "Initial invested capital")
	interest := flag.Float64("interest", 3.0, "Expected period ROI")
	interest_period := flag.String("interest_period", "weekly", "weekly, monthly or yearly")
	years := flag.Int("years", 5, "Years of investment")
	top := flag.Float64("top", 0.0, "Monthly investment")
	notifyAt := flag.Float64("notify_at", 0.0, "Stop at")
	flag.Parse()

	date := time.Now().UTC()

	for end := date.AddDate(*years, 0, 0); end.After(date); {
		earn := *capital * *interest / float64(100.0)
		*capital += earn
		fmt.Printf("At %s you have %.2f\n", date.String(), *capital)

		if *notifyAt > 0 && *capital >= *notifyAt {
			break
		}
		*capital += *top

		switch *interest_period {
		case "weekly":
			date = date.AddDate(0, 0, 7)
		case "monthly":
			date = date.AddDate(0, 1, 0)
		default:
			date = date.AddDate(1, 0, 0)
		}

	}
}
