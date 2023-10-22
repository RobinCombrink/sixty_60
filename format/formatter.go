package format

import (
	"fmt"
	"log"
	"time"
)
const inputFormat = "02 January 2006, 15:04"
const outputFormat = "02 January 2006"
func ToRand(valueInCents uint64) string {
	return fmt.Sprintf("R%.2f", float64(valueInCents)/100)
}
func ToDate(input string) (date time.Time) {
	date, err := time.Parse(inputFormat, input)
	if err != nil {
		log.Fatalf("Failed to parse time: %v\n", input)
	}
	return date
}

func ToReadableDate(input time.Time) (readableDate string) {
	readableDate = input.Format(outputFormat)
	return readableDate
}