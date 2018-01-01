package main

import (
	"errors"
	"strconv"
	"time"
)

// JSONTimezone custom timeone struct for unmarshall
type JSONTimezone struct {
	*time.Location
}

// Schedule when a job shuld be running
type Schedule struct {
	Timezone  JSONTimezone `json:"timezone"`
	DayOfWeek []uint8      `json:"dayOfWeek"`
	Month     []uint8      `json:"month"`
	Day       []uint8      `json:"day"`
	Hour      []uint8      `json:"hour"`
	Minute    []uint8      `json:"minute"`
}

// UnmarshalJSON unmarshall a sjson timezone string
func (l *JSONTimezone) UnmarshalJSON(b []byte) error {
	// you can now parse b as thoroughly as you want
	zoneString, _ := strconv.Unquote(string(b))
	loc, err := time.LoadLocation(zoneString)
	if err != nil {
		return err
	}
	l.Location = loc
	return nil
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func containsInt(this int, in []uint8) bool {
	for _, i := range in {
		if uint8(this) == i {
			return true
		}
	}
	return false
}

func (s *Schedule) isThisDayOfWeek(date time.Time) bool {
	if s.DayOfWeek == nil || len(s.DayOfWeek) == 0 {
		return true
	}
	dayOfWeek := int(date.Weekday())
	return containsInt(dayOfWeek, s.DayOfWeek)
}

func (s *Schedule) isThisMonth(date time.Time) bool {
	if s.Month == nil || len(s.Month) == 0 {
		return true
	}
	month := int(date.Month())
	return containsInt(month, s.Month)
}

func (s *Schedule) isThisDay(date time.Time) bool {
	if s.Day == nil || len(s.Day) == 0 {
		return true
	}
	day := date.Day()
	return containsInt(day, s.Day)
}

func (s *Schedule) isThisHour(date time.Time) bool {
	if s.Hour == nil || len(s.Hour) == 0 {
		return true
	}
	hour := date.Hour()
	return containsInt(hour, s.Hour)
}

func (s *Schedule) isThisMinute(date time.Time) bool {
	if s.Minute == nil || len(s.Minute) == 0 {
		return true
	}
	minute := date.Minute()
	return containsInt(minute, s.Minute)
}

func (s *Schedule) nextHour(date time.Time) (int, error) {
	var err error

	if s.Minute == nil {
		date = date.Add(time.Hour)
		return date.Hour(), err
	}

	for i := date.Hour(); i <= 23; i++ {
		if containsInt(i, s.Hour) {
			return i, err
		}
	}
	for i := 0; i <= date.Hour(); i++ {
		if containsInt(i, s.Hour) {
			return i, err
		}
	}

	err = errors.New("Failed to find next hour")
	return 0, err
}

func (s *Schedule) nextMinute(date time.Time) (int, error) {
	var err error
	date = date.Add(time.Minute)

	if s.Minute == nil {
		return date.Minute(), err
	}

	for i := date.Minute(); i <= 59; i++ {
		if containsInt(i, s.Minute) {
			return i, err
		}
	}
	for i := 0; i <= date.Minute(); i++ {
		if containsInt(i, s.Minute) {
			return i, err
		}
	}

	err = errors.New("Failed to find next minute")
	return 0, err
}

func (s *Schedule) isNow(date time.Time) bool {
	/*
		println("===")
		spew.Dump(s.isThisDayOfWeek(date),
			s.isThisMonth(date),
			s.isThisDay(date),
			s.isThisHour(date),
			s.isThisMinute(date))
		println("---")
	*/
	if s.Timezone.String() != "" {
		date = date.In(s.Timezone.Location)
	}
	if s.isThisDayOfWeek(date) &&
		s.isThisMonth(date) &&
		s.isThisDay(date) &&
		s.isThisHour(date) &&
		s.isThisMinute(date) {
		return true
	}
	return false
}

func (s *Schedule) findNextRun(date time.Time) time.Time {
	dateExpire := date.AddDate(1, 0, 0)
	date = date.Add(time.Minute)

	for checkDate := date; checkDate.Before(dateExpire); checkDate = checkDate.Add(time.Minute) {
		if s.isNow(checkDate) {
			return checkDate
		}
	}

	return time.Time{}
}
