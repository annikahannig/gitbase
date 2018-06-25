package gitbase

/*
Datetime supplemental:
Parse unix timestamps with offsets

Update / set location
*/

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrInvalidOffsetFormat = errors.New("invalid offset format")
)

func DatetimeParseOffset(offset string) (*time.Location, error) {
	// Offset described as +-HHMM
	if len(offset) != 5 {
		return nil, ErrInvalidOffsetFormat
	}

	sign := string(offset[0])
	hh := string(offset[1:3])
	mm := string(offset[3:5])

	h, err := strconv.ParseInt(hh, 10, 64)
	if err != nil {
		return nil, err
	}

	m, err := strconv.ParseInt(mm, 10, 64)
	if err != nil {
		return nil, err
	}

	// Calculate offset in seconds
	locOffset := int(h*60*60 + m*60)

	if sign == "-" {
		locOffset *= -1
	}

	// As name use UTC+Offset in hours
	name := fmt.Sprintf("UTC%s%d", sign, h) // TODO: Room for improvement

	// Make location
	loc := time.FixedZone(name, locOffset)

	return loc, nil
}

func DatetimeSetLocation(datetime time.Time, loc *time.Location) time.Time {
	return time.Date(
		datetime.Year(),
		datetime.Month(),
		datetime.Day(),
		datetime.Hour(),
		datetime.Minute(),
		datetime.Second(),
		datetime.Nanosecond(),
		loc)
}
