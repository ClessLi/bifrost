package test

import (
	"testing"
	"time"
)

func TestTimeZone(t *testing.T) {
	TZ, err := time.LoadLocation("Local")
	t.Log("time zone", TZ, err)

	now := time.Now().Local().Add(time.Hour * -12)
	t.Log("now", now)

	dt := now.Format("20060102")
	t.Log("date time", dt)

	tt := func(date time.Time) {
		t.Log(date.Format("20060102"))
	}
	tt(now)
}
