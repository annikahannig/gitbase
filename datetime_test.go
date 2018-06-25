package gitbase

import (
	"testing"
	"time"
)

func TestDatetimeParseLocation(t *testing.T) {
	time.Now()

	// Parse location
	loc, err := DatetimeParseOffset("+0200")
	if err != nil {
		t.Error(err)
	}
	if loc == nil {
		t.Error("did not parse offset")
	}

	_, err = DatetimeParseOffset("2342")
	if err != ErrInvalidOffsetFormat {
		t.Error("did not expect offset to parse")
	}

	_, err = DatetimeParseOffset("+asdf")
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestDatetimeSetLocation(t *testing.T) {
	utc0 := time.Unix(0, 0).UTC()

	loc0200, err := DatetimeParseOffset("+0200")
	if err != nil {
		t.Error(err)
	}

	dt0200 := DatetimeSetLocation(utc0, loc0200)

	// UTC time should now be utc200 - 2h
	utc0200 := dt0200.UTC()

	if utc0200.Hour() != 22 &&
		utc0200.Day() != 31 &&
		utc0200.Month() != 12 {
		t.Error("Date conversion failed. Expected 1969-12-31 22:00:00")
	}

}
