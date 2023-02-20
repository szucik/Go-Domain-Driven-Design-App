package clock

import "time"

func FakeTime() time.Time {
	dateString := "2021-11-22"
	date, _ := time.Parse("2006-01-02", dateString)
	return date
}
