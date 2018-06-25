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
