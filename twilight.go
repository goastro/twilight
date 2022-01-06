package twilight

import (
	"time"
)

func lenToDuration(len float64) time.Duration {
	return time.Duration(float64(time.Hour) * len)
}

type SunriseStatus int

const (
	SunriseStatusOK           = SunriseStatus(0)
	SunriseStatusAboveHorizon = SunriseStatus(1)
	SunriseStatusBelowHorizon = SunriseStatus(-1)
)

func DayLength(d time.Time, latitude, longitude float64) time.Duration {
	len := dayLen(d.Year(), int(d.Month()), d.Day(), longitude, latitude, -35.0/60.0, true)

	return lenToDuration(len)
}

func CivilTwilightLength(d time.Time, latitude, longitude float64) time.Duration {
	len := dayLen(d.Year(), int(d.Month()), d.Day(), longitude, latitude, -6.0, false)

	return lenToDuration(len)
}

func NauticalTwilightLength(d time.Time, latitude, longitude float64) time.Duration {
	len := dayLen(d.Year(), int(d.Month()), d.Day(), longitude, latitude, -12.0, false)

	return lenToDuration(len)
}

func AstronomicalTwilightLength(d time.Time, latitude, longitude float64) time.Duration {
	len := dayLen(d.Year(), int(d.Month()), d.Day(), longitude, latitude, -18.0, false)

	return lenToDuration(len)
}

func riseSetToTime(now time.Time, rise, set float64) (time.Time, time.Time) {
	y, m, d := now.Date()

	today := time.Date(y, m, d, 0, 0, 0, 0, now.Location())

	riseDuration := lenToDuration(rise)
	setDuration := lenToDuration(set)

	return today.Add(riseDuration), today.Add(setDuration)
}

func SunRiseSet(date time.Time, latitude, longitude float64) (time.Time, time.Time, SunriseStatus) {
	r, s, status := sunRiseSet(date.Year(), int(date.Month()), date.Day(), longitude, latitude, -35.0/60.0, true)
	rise, set := riseSetToTime(date, r, s)

	return rise, set, status
}

func CivilTwilight(date time.Time, latitude, longitude float64) (time.Time, time.Time, SunriseStatus) {
	r, s, status := sunRiseSet(date.Year(), int(date.Month()), date.Day(), longitude, latitude, -6.0, true)
	rise, set := riseSetToTime(date, r, s)

	return rise, set, status
}

func NauticalTwilight(date time.Time, latitude, longitude float64) (time.Time, time.Time, SunriseStatus) {
	r, s, status := sunRiseSet(date.Year(), int(date.Month()), date.Day(), longitude, latitude, -12.0, true)
	rise, set := riseSetToTime(date, r, s)

	return rise, set, status
}

func AstronomicalTwilight(date time.Time, latitude, longitude float64) (time.Time, time.Time, SunriseStatus) {
	r, s, status := sunRiseSet(date.Year(), int(date.Month()), date.Day(), longitude, latitude, -18.0, true)
	rise, set := riseSetToTime(date, r, s)

	return rise, set, status
}
