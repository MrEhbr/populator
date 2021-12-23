package helpers

import (
	"fmt"
	"time"
)

var timeFormats = [...]string{
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	"2006-01-02",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05",
}

func TryStrToDate(s string) (time.Time, error) {
	for _, l := range timeFormats {
		t, err := time.Parse(l, s)
		if err != nil {
			continue
		}
		return t, nil
	}
	return time.Time{}, fmt.Errorf(`could not convert string %q to time`, s)
}
