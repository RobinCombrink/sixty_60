package format

import (
	"fmt"
	"log"
	"time"
)

const htmlInputFormat = "02 January 2006, 15:04"
const htmlOutputFormat = "02 January 2006"
const datePickerSetterFormat = "2006-01-02"

func ToRand(valueInCents uint64) string {
	return fmt.Sprintf("R%.2f", float64(valueInCents)/100)
}
func ToDate(input string) (date time.Time) {
	date, err := time.Parse(htmlInputFormat, input)
	if err != nil {
		log.Fatalf("Failed to parse time: %v\n", input)
	}
	return date
}

func ToReadableDate(input time.Time) (readableDate string) {
	readableDate = input.Format(htmlOutputFormat)
	return readableDate
}

func ToDatePickerDate(input time.Time) (datePickerDate string) {
	datePickerDate = input.Format(datePickerSetterFormat)
	return datePickerDate
}
